param(
    [string]$RunAsProcess = "false",
    [string]$BuildExe = "true"
)

if ($BuildExe -eq "true") {
    Write-Host "Building thermostat controller"
    go1.22.5 mod download
    go1.22.5 build -o thermostat.exe
}

Write-Host "Running thermostat controller"
$cmd = ".\thermostat.exe"
if ($RunAsProcess -eq "true") {
    $cmd = "Start-Process -FilePath $cmd"
}
& $cmd
exit !$?
