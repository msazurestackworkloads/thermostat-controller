# param to run as process or exe
param(
    [string]$process = "false",
    [string]$build = "true"
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
