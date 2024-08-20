param(
    [string]$RunAsProcess = "false",
    [string]$BuildExe = "true",
    [string]$BuildPath = "."
)

if ($BuildExe -eq "true") {
    Write-Host "Building thermostat controller"
    go1.22.5 mod download
    go1.22.5 build -o $BuildPath\thermostat.exe -v
}

Write-Host "Running thermostat controller"
$cmd = "$BuildPath\thermostat.exe"
if ($RunAsProcess -eq "true") {
    $cmd = "Start-Process -FilePath $cmd"
}
Invoke-Expression $cmd
