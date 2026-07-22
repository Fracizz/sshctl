# Install sshctl to C:\Program Files\sshctl and add to machine PATH.
# Requires Administrator.
$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$bin = Join-Path $repoRoot "bin\sshctl.exe"
if (-not (Test-Path $bin)) {
    Write-Host "building $bin"
    Push-Location $repoRoot
    $env:GOROOT = 'C:\Program Files\Go'
    & 'C:\Program Files\Go\bin\go.exe' build -o bin\sshctl.exe .
    Pop-Location
}

& $bin install
if ($LASTEXITCODE -ne 0) {
    Write-Host "Tip: run this script in an elevated (Administrator) PowerShell."
    exit $LASTEXITCODE
}

Write-Host "Done. Open a new terminal and run: sshctl version"
