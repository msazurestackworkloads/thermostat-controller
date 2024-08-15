# param to run as process or exe
param([string]$process = "false")

if (-not (Test-Path .\thermostat.exe)) {
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
