param (
	[parameter(Mandatory=$false)]
	$WorkingDirectory,
	[parameter(Mandatory=$false)]
	$ApkFile
)

# Handle working directory
# Make absolute and check it's existence
if (-not $WorkingDirectory) {
	$WorkingDirectory = "."
}
$ResWD = (Resolve-Path -Path $WorkingDirectory -ErrorAction SilentlyContinue)
if (-not $ResWD) {
	Write-Error "The specified directory '$WorkingDirectory' does not exist."
	return
} else {
	$WorkingDirectory = $ResWD.Path
}

Write-Host "Working Directory: $WorkingDirectory"
Write-Host "Starting..."
Write-Host "------------------"

# Apk files was not provided via cmdline
if (-not $ApkFile) {
	
	Write-Host "No APK file specified. Looking in working directory..."
	$ApkFiles = Get-ChildItem -Path $WorkingDirectory -File | Where-Object { $_.Extension -in ".apk",".xapk" }
	
	# No apkfiles in working directory. download it.
	if ($ApkFiles.Count -eq 0) {
		Write-Host "No APK files found in the directory '$WorkingDirectory'"
		Write-Host "Downloading from UpToDown.com"
		
		# Static download link to 4.0.2 because the script may break with new versions anyways
		# Also, it will take months from now until they can feasibly change anything about the certificates or api
		
		# Send a request to the URL, then extract the dynamic download link from the response
		$BaseUrl = "https://rewe.de.uptodown.com/android/download/1047828586-x" 
		try {
			$Response = Invoke-WebRequest -Uri $BaseUrl -UseBasicParsing
			if ($response.StatusCode -ne 200) {
				throw "Issue getting data-url"
			}
		} catch {
			Write-Error "Failed to fetch the APK. Download the apk manually from https://rewe.en.uptodown.com/android/"
			return
		}
		
		# Match the regex pattern in the response content
		if ($Response.Content -match 'data-url="([^"]+)"') {
			$DataUrl = $matches[1] 
		} else {
			Write-Error "No data-url found in the response content."
			return
		}
		
		# Downloading the file
		Write-Host "Starting download. Please wait..."
		$ApkUrl = "https://dw.uptodown.com/dwn/$DataUrl"
		$ApkFile = Join-Path -Path $WorkingDirectory -ChildPath "rewe.apk"
		$ApkFileDownload = (New-Object Net.WebClient).DownloadFile($ApkUrl, $ApkFile)
		Write-Host "Done"
	
	} elseif ($ApkFiles.Count -gt 1) {
		Write-Host "Multiple APK files found in the directory"
		$ApkFile = ($ApkFiles[0]).FullName
	
	} else {
		$ApkFile = ($ApkFiles[0]).FullName
	}
} 

# Check if the file exists
if (-not (Test-Path -Path $ApkFile -PathType Leaf)) {
	Write-Error "APK file '$ApkFile' does not exist."
	return
} else {
	Write-Host "Using APK file: $($ApkFile)"
}
Write-Host "------------------"

$ErrorActionPreference = "Stop"

# Open the APK file as a zip archive
Add-Type -AssemblyName System.IO.Compression.FileSystem
$Zip = [System.IO.Compression.ZipFile]::OpenRead($ApkFile)

# Check for nested APK (apkx packing)
$NestedApk = $Zip.GetEntry("de.rewe.app.mobile.apk")
if ($NestedApk) {
    # Read nested APK into memory stream
    $NestedStream = New-Object System.IO.MemoryStream
    $NestedApk.Open().CopyTo($NestedStream)
    $Zip.Dispose()
    
    # Reset stream position and create new zip archive
    $NestedStream.Position = 0
    $Zip = [System.IO.Compression.ZipArchive]::new($NestedStream)
}

# Find the pfx and copy it out of the zip
$ExpectedCertName = "mtls_prod.pfx"
$Entry = $Zip.GetEntry("res/raw/$ExpectedCertName")
if($Entry) {
	$Dest = Join-Path -Path $WorkingDirectory -ChildPath $ExpectedCertName
	$EntrStr = $Entry.Open()
	$FileStr = [System.IO.File]::Create($Dest)
	$EntrStr.CopyTo($FileStr)
	$FileStr.Close()
	$EntrStr.Close()
	Write-Host "Extracted $ExpectedCertName to $Dest"
	$Zip.Dispose()
} else {
	Write-Error "Certificate $ExpectedCertName not found in the APK."
	$Zip.Dispose()
	return
}

Write-Host "------------------"

# Extract pem and key
$Cert = New-Object System.Security.Cryptography.X509Certificates.X509Certificate2($Dest, "NC3hDTstMX9waPPV", [System.Security.Cryptography.X509Certificates.X509KeyStorageFlags]::Exportable)

$PemPath = Join-Path -Path $WorkingDirectory -ChildPath "certificate.pem"
$CertBytes = $cert.Export([System.Security.Cryptography.X509Certificates.X509ContentType]::Cert)
$Base64Cert = [Convert]::ToBase64String($CertBytes)
$PemContent = "-----BEGIN CERTIFICATE-----`r`n$Base64Cert`r`n-----END CERTIFICATE-----"
Set-Content -Path $PemPath -Value $PemContent

$KeyPath = Join-Path -Path $WorkingDirectory -ChildPath "private.key"
try {
    if ($IsLinux) {
        $RSA = [System.Security.Cryptography.RSA]::Create()
        $KeyBytes = $Cert.PrivateKey.ExportPkcs8PrivateKey()
    } else {
        $RSACng = [Security.Cryptography.X509Certificates.RSACertificateExtensions]::GetRSAPrivateKey($Cert)
        $KeyBytes = $RSACng.Key.Export([Security.Cryptography.CngKeyBlobFormat]::Pkcs8PrivateBlob)
    }

    $KeyBase64 = [Convert]::ToBase64String($KeyBytes, [Base64FormattingOptions]::InsertLineBreaks)
    Set-Content -Path $KeyPath -Value @"
-----BEGIN PRIVATE KEY-----
$KeyBase64
-----END PRIVATE KEY-----
"@
} finally {
    if ($null -ne $RSA) {
        $RSA.Dispose()
    }
}

Write-Host "Keys exported successfully."
Write-Host "Done :)"
