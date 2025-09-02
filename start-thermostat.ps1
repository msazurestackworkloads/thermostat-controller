param(
    [bool] $RunAsProcess = $false,
    [bool] $BuildExe = $true
)

if ($BuildExe) {
    Write-Host "Building thermostat controller"
    go1.24.1 mod download
    go1.24.1 build -o thermostat.exe
}

Write-Host "Running thermostat controller"
$cmd = ".\thermostat.exe"
if ($RunAsProcess) {
    $cmd = "Start-Process -FilePath $cmd"
}
& $cmd
exit !$?
