## a basic kubernetes thermostat controller/operator

creates and reconciles CRs on a cluster provided a kube api server url

attempts to match current temperature to desired temperature by adjusting the current temperature by 1 each reconcile loop and deletes the CR once complete

## Build and run

Use Go 1.27.3 or later, as specified in `go.mod`:

```powershell
.\start-thermostat.ps1
```

The launcher builds and runs relative to its own directory, including when called
from another directory. It stops on build or controller failures instead of
running a stale executable or reporting success.

An E2E caller can reuse its selected Go toolchain without maintaining a second
version pin in the launcher:

```powershell
.\start-thermostat.ps1 -GoCommand go1.27.3
```

Use `-BuildExe $false` to run an existing binary, or `-RunAsProcess $true` to
start it in the background and return its process object.

```powershell
go test ./...
```
