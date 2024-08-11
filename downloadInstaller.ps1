# Specify the GitHub repository owner and repository name
$repoOwner = "alabsi91"
$repoName = "win-tools-go"
$fileName = "installer.exe"

# Construct the URL to fetch the latest release information
$latestReleaseUrl = "https://api.github.com/repos/$repoOwner/$repoName/releases/latest"

# Fetch the latest release information from GitHub
$response = Invoke-RestMethod -Uri $latestReleaseUrl

# Get the tag name of the latest release
$latestTagName = $response.tag_name

# Construct the URL for the release asset
$assetUrl = "https://github.com/$repoOwner/$repoName/releases/download/$latestTagName/$fileName"

# Specify the local path where you want to save the downloaded installer
$localFilePath = "$env:TEMP\$fileName"

# Download the installer.exe from GitHub releases
Invoke-WebRequest -Uri $assetUrl -OutFile $localFilePath

# Check if the file was downloaded successfully
if (Test-Path $localFilePath) {
    # Launch the installer
    Start-Process -FilePath $localFilePath -Wait -ArgumentList "/S"
    Start-Process "win-tools"
} else {
    Write-Host "Failed to download the installer."
}
