# Cross-build sshfrac for Linux / Windows / macOS and package release zips.
$ErrorActionPreference = "Stop"
$Version = if ($env:VERSION) { $env:VERSION } else { "0.1.1" }
$ld = "-s -w -X github.com/Fracizz/sshfrac/cmd.Version=$Version"
$go = if ($env:GOEXE) { "go" } else { "go" }
if (Test-Path "C:\Program Files\Go\bin\go.exe") {
    $env:GOROOT = "C:\Program Files\Go"
    $go = "C:\Program Files\Go\bin\go.exe"
}

if (Test-Path dist) {
    Remove-Item -Recurse -Force dist
}
New-Item -ItemType Directory -Force -Path dist | Out-Null

$targets = @(
    @{ GOOS = "linux"; GOARCH = "amd64"; Out = "dist/sshfrac-linux-amd64" },
    @{ GOOS = "linux"; GOARCH = "arm64"; Out = "dist/sshfrac-linux-arm64" },
    @{ GOOS = "windows"; GOARCH = "amd64"; Out = "dist/sshfrac-windows-amd64.exe" },
    @{ GOOS = "windows"; GOARCH = "arm64"; Out = "dist/sshfrac-windows-arm64.exe" },
    @{ GOOS = "darwin"; GOARCH = "amd64"; Out = "dist/sshfrac-darwin-amd64" },
    @{ GOOS = "darwin"; GOARCH = "arm64"; Out = "dist/sshfrac-darwin-arm64" }
)

foreach ($t in $targets) {
    $env:GOOS = $t.GOOS
    $env:GOARCH = $t.GOARCH
    $env:CGO_ENABLED = "0"
    Write-Host "building $($t.Out)"
    & $go build -ldflags $ld -o $t.Out .
}
Remove-Item Env:GOOS, Env:GOARCH -ErrorAction SilentlyContinue

Get-ChildItem dist -File | ForEach-Object {
    $bundleName = $_.BaseName
    $bundleDir = Join-Path dist $bundleName
    if (Test-Path $bundleDir) {
        Remove-Item -Recurse -Force $bundleDir
    }
    New-Item -ItemType Directory -Force -Path $bundleDir | Out-Null

    if ($_.Extension -eq ".exe") {
        Copy-Item $_.FullName (Join-Path $bundleDir "sshfrac.exe")
    } else {
        Copy-Item $_.FullName (Join-Path $bundleDir "sshfrac")
    }

    $zipPath = Join-Path dist "$bundleName.zip"
    if (Test-Path $zipPath) {
        Remove-Item -Force $zipPath
    }
    Compress-Archive -Path (Join-Path $bundleDir "*") -DestinationPath $zipPath
    Write-Host "packaged $zipPath"
}

Copy-Item dist/sshfrac-windows-amd64.exe bin/sshfrac.exe -Force
Write-Host "done sshfrac $Version"
