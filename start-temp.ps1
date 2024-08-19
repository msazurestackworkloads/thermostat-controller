# param to run as process or exe and to build or just run
param(
    [string]$runasprocess = "false",
    [string]$buildexe = "true"
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
