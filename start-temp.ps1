param(
    [string]$RunAsProcess = "false",
    [string]$BuildExe = "true"
)

if ($build -eq "true") {
    Write-Host "Building thermostat controller"
    go mod download
    go build -o thermostat.exe -v
}

Write-Host "Running thermostat controller"
$cmd = ".\thermostat.exe"
if ($process -eq "true") {
    $cmd = "Start-Process -FilePath $cmd"
}
Invoke-Expression $cmd
