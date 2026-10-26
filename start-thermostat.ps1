param(
    [bool] $RunAsProcess = $false,
    [bool] $BuildExe = $true,
    [string] $GoCommand = "go"
)

$ErrorActionPreference = "Stop"
$exePath = Join-Path $PSScriptRoot "thermostat.exe"

if ($BuildExe) {
    Write-Host "Building thermostat controller"
    & $GoCommand -C $PSScriptRoot build -mod=readonly -o $exePath .
    if ($LASTEXITCODE -ne 0) {
        throw "Building thermostat controller failed with exit code $LASTEXITCODE."
    }
}

Write-Host "Running thermostat controller"
if ($RunAsProcess) {
    Start-Process -FilePath $exePath -WorkingDirectory $PSScriptRoot -PassThru
} else {
    Push-Location $PSScriptRoot
    try {
        & $exePath
        if ($LASTEXITCODE -ne 0) {
            throw "Thermostat controller failed with exit code $LASTEXITCODE."
        }
    } finally {
        Pop-Location
    }
}
