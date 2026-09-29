param (
	[parameter(Mandatory=$false)]
	$WorkingDirectory,
	
	[parameter(Mandatory=$false)]
	$ApkFile,
	
	[parameter(Mandatory = $false, ParameterSetName = "DefaultPath")]
	[switch]$Pfx,
	
	[parameter(Mandatory = $true, ParameterSetName = "CustomPath")]
	[string]$PfxPath
)

# ——————————————————————————————————————————————————————————————————————————————
# Functions
# ——————————————————————————————————————————————————————————————————————————————

function Export-PfxToPemKey {
    [CmdletBinding()]
    param (
        [Parameter(Mandatory = $true)]
        [string]$PfxPath,
        
        [Parameter(Mandatory = $true)]
        [string]$WorkingDirectory
    )
    
    try {
        # Create X509Certificate2 object from the PFX file
        $Cert = New-Object System.Security.Cryptography.X509Certificates.X509Certificate2(
            $PfxPath, 
            "NC3hDTstMX9waPPV", 
            [System.Security.Cryptography.X509Certificates.X509KeyStorageFlags]::Exportable
        )
        
        # Export certificate to PEM format
        $PemPath = Join-Path -Path $WorkingDirectory -ChildPath "certificate.pem"
        $CertBytes = $Cert.Export([System.Security.Cryptography.X509Certificates.X509ContentType]::Cert)
        $Base64Cert = [Convert]::ToBase64String($CertBytes)
        $PemContent = "-----BEGIN CERTIFICATE-----`r`n$Base64Cert`r`n-----END CERTIFICATE-----"
        Set-Content -Path $PemPath -Value $PemContent
        
        # Export private key
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
        return @{
            CertificatePath = $PemPath
            PrivateKeyPath = $KeyPath
        }
    }
    catch {
        Write-Error "Error exporting PFX to PEM and key: $_"
        throw
    }
}

# ——————————————————————————————————————————————————————————————————————————————
# Main
# ——————————————————————————————————————————————————————————————————————————————

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

if ($Pfx -or $PfxPath) {
    $PfxFileToUse = if ($Pfx) {
        Join-Path -Path $WorkingDirectory -ChildPath "mtls_prod.pfx"
    } else {
        $PfxPath
    }
	
	if (-not [System.IO.Path]::IsPathRooted($PfxFileToUse)) {
		$PfxFileToUse = [System.IO.Path]::GetFullPath((Join-Path -Path (Get-Location) -ChildPath $PfxFileToUse))
	}

    # Validate that the PFX file exists
    if (-not (Test-Path -Path $PfxFileToUse -PathType Leaf)) {
        Write-Error "PFX file not found: $PfxFileToUse"
        exit 1
    }

    # Call the Export-PfxToPemKey function
    try {
        Write-Host "Exporting PFX file: $PfxFileToUse"
        $result = Export-PfxToPemKey -PfxPath $PfxFileToUse -WorkingDirectory $WorkingDirectory
        Write-Host "Certificate: $($result.CertificatePath)"
        Write-Host "Private Key: $($result.PrivateKeyPath)"
    } catch {
        Write-Error "Failed to export PFX file: $_"
        exit 1
    }
	return
}

# Apk files was not provided via cmdline
if (-not $ApkFile) {

	$ApkFiles = Get-ChildItem -Path $WorkingDirectory -File | Where-Object { $_.Extension -in ".apk",".xapk" }

	# No apkfiles in working directory
	if ($ApkFiles.Count -eq 0) {
		Write-Host "No APK file found in '$WorkingDirectory'"
		Write-Host "Download the REWE xapk in a browser from https://rewe.en.uptodown.com/android/download/1220537953-x"
		Write-Host "Read the buttons carefully. You want the xapk file, NOT the Uptodown installer or app store."
		Write-Host "Then place it there or pass it with -ApkFile <path>"
		exit 1
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
	exit 1
} else {
	Write-Host "Using APK file: $($ApkFile)"
}
Write-Host "------------------"

$ErrorActionPreference = "Stop"

# Open the APK file as a zip archive
Add-Type -AssemblyName System.IO.Compression.FileSystem
$Zip = [System.IO.Compression.ZipFile]::OpenRead($ApkFile)

try {
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
} finally {
    $Zip.Dispose()
}

Write-Host "------------------"

Export-PfxToPemKey -PfxPath $Dest -WorkingDirectory $WorkingDirectory
Write-Host "Done :)"
