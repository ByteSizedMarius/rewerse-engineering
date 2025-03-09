param (
	[parameter(Mandatory=$false)]
	$WorkingDirectory,
	[parameter(Mandatory=$false)]
	$ApkFile,
    [parameter(Mandatory = $false, ParameterSetName = "DefaultPath")]
    [switch]$Pfx,
    [parameter(Mandatory = $true, ParameterSetName = "CustomPath")]
    [string]$PfxPath,
	[parameter(Mandatory = $false)]
    [switch]$Dl
)

# ——————————————————————————————————————————————————————————————————————————————
# Functions
# ——————————————————————————————————————————————————————————————————————————————

# APK download target version; no need to update unless certificates change
$targetVersion = "4.1.0"

function Get-ApkVersionUrl {
    try {
        $response = Invoke-WebRequest -UseBasicParsing -Uri "https://rewe.en.uptodown.com/android/versions"
        $htmlDoc = New-Object -ComObject "HTMLFile"
        $htmlDoc.IHTMLDocument2_write($response.Content)
        
        # Get version info
        $versionsDiv = $htmlDoc.getElementById("versions-items-list")
        if ($null -eq $versionsDiv) {
            Write-Error "Could not find versions list on the page"
            return $null
        }
        
        $versionDivs = $versionsDiv.getElementsByTagName("div")
        $targetVersionDiv = $null
        for ($i = 0; $i -lt $versionDivs.length; $i++) {
            $versionElement = $versionDivs[$i].getElementsByClassName("version")
            if ($versionElement.length -gt 0) {
                $versionText = $versionElement[0].innerText
				
                # Check if this is our target version
                if ($versionText -eq $targetVersion) {
                    $targetVersionDiv = $versionDivs[$i]
                    break
                }
            }
        }
        
        # Check if target version was found
        if ($null -eq $targetVersionDiv) {
            Write-Error "Target version $targetVersion not found on the page"
            return $null
        }
        
        return $targetVersionDiv.getAttribute("data-url")
    } catch {
        Write-Error "Error extracting download URL: $_"
        return $null
    } finally {
        # Release COM object
        if ($null -ne $htmlDoc) {
            [System.Runtime.InteropServices.Marshal]::ReleaseComObject($htmlDoc) | Out-Null
        }
        [System.GC]::Collect()
    }
}

function Get-ApkDownloadUrl {
    [CmdletBinding()]
    param()
    
    try {
        # Get the initial URL for the app version page
        $initialUrl = Get-ApkVersionUrl
        if ($null -eq $initialUrl) {
            Write-Error "Failed to get initial URL"
            return $null
        }
        Write-Verbose "Initial URL: $initialUrl"
        
        # Request the version page
        $versionPageResponse = Invoke-WebRequest -UseBasicParsing -Uri $initialUrl
        $htmlDoc = New-Object -ComObject "HTMLFile"
        $htmlDoc.IHTMLDocument2_write($versionPageResponse.Content)
        
        # Extract data-version
        $variantButton = $null
        $buttons = $htmlDoc.getElementsByTagName("button")
        for ($i = 0; $i -lt $buttons.length; $i++) {
            if ($buttons[$i].className -eq "button variants") {
                $variantButton = $buttons[$i]
                break
            }
        }
        if ($null -eq $variantButton) {
            Write-Error "Could not find variants button"
            return $null
        }
        $dataVersion = $variantButton.getAttribute("data-version")
        Write-Verbose "Data Version: $dataVersion"
        
        # Extract data-code
        $appNameH1 = $htmlDoc.getElementById("detail-app-name")
        if ($null -eq $appNameH1) {
            Write-Error "Could not find app name element"
            return $null
        }
        $dataCode = $appNameH1.getAttribute("data-code")
        Write-Verbose "Data Code: $dataCode"
        
        # Construct and request the variants URL
        $variantsUrl = "https://rewe.en.uptodown.com/app/$dataCode/version/$dataVersion/files"
        Write-Verbose "Variants URL: $variantsUrl"
        $variantsResponse = Invoke-WebRequest -UseBasicParsing -Uri $variantsUrl
        $jsonResponse = $variantsResponse.Content | ConvertFrom-Json
        $htmlContent = $jsonResponse.content
        $variantsHtmlDoc = New-Object -ComObject "HTMLFile"
        $variantsHtmlDoc.IHTMLDocument2_write($htmlContent)
        
        # Get all divs
        $allDivs = $variantsHtmlDoc.getElementsByTagName("div")
        Write-Verbose "Found $($allDivs.length) divs total"
        
        # Find first div with class "v-version" and extract URL
        $downloadUrl = $null
        for ($i = 0; $i -lt $allDivs.length; $i++) {
            if ($allDivs[$i].className -eq "v-version") {
                $versionDiv = $allDivs[$i]
                $version = $versionDiv.innerText
                
                # Get the HTML and extract the onclick using regex
                $html = $versionDiv.outerHTML
                if ($html -match 'onclick=[''"].*?location\.href=[''"]([^''"]+)[''"]') {
                    $downloadUrl = $Matches[1]
                    Write-Verbose "Found version: $version with URL: $downloadUrl"
                    break
                }
            }
        }
        
        return $downloadUrl
    } catch {
        Write-Error "Error extracting download URL: $_"
        return $null
    } finally {
        # Release COM objects
        if ($null -ne $htmlDoc) {
            [System.Runtime.InteropServices.Marshal]::ReleaseComObject($htmlDoc) | Out-Null
        }
        if ($null -ne $variantsHtmlDoc) {
            [System.Runtime.InteropServices.Marshal]::ReleaseComObject($variantsHtmlDoc) | Out-Null
        }
        [System.GC]::Collect()
    }
}

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

function Download-ReweApk {
    [CmdletBinding()]
    param (
        [Parameter(Mandatory = $true)]
        [string]$WorkingDirectory
    )
    
    # Ensure the working directory exists
    if (-not (Test-Path -Path $WorkingDirectory -PathType Container)) {
        New-Item -Path $WorkingDirectory -ItemType Directory -Force | Out-Null
        Write-Host "Created working directory: $WorkingDirectory"
    }
    
    # Function to get the APK download base URL
    # Note: This function should be defined elsewhere or adapted as needed
    $BaseUrl = Get-ApkDownloadUrl
    
    try {
        # Send a request to the URL, then extract the dynamic download link from the response
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
    
    return $ApkFile
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

if ($Dl) {
	Write-Host "Downloading APK from UpToDown.com"	
	$ApkFile = Download-ReweApk $WorkingDirectory
	Write-Host "Downloaded to $ApkFile"
	return
}
	
# Apk files was not provided via cmdline
if (-not $ApkFile) {
	
	Write-Host "No APK file specified. Looking in working directory..."
	$ApkFiles = Get-ChildItem -Path $WorkingDirectory -File | Where-Object { $_.Extension -in ".apk",".xapk" }
	
	# No apkfiles in working directory. download it.
	if ($ApkFiles.Count -eq 0) {
		Write-Host "No APK files found in the directory '$WorkingDirectory'"
		Write-Host "Downloading from UpToDown.com"
		
		$ApkFile = Download-ReweApk $WorkingDirectory
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
