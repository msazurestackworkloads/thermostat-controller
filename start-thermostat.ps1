param(
    [bool] $RunAsProcess = $false,
    [bool] $BuildExe = $true
)

$goversion = "1.25.5"

if ($BuildExe) {
    Write-Host "Building thermostat controller"
    & "go$goversion" mod download
    & "go$goversion" build -o thermostat.exe
}

Write-Host "Running thermostat controller"
$cmd = ".\thermostat.exe"
if ($RunAsProcess) {
    $cmd = "Start-Process -FilePath $cmd"
}
& $cmd
exit !$?
