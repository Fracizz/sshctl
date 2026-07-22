# Cross-build sshctl for Linux / Windows / macOS and package release zips.
$ErrorActionPreference = "Stop"
$Version = if ($env:VERSION) { $env:VERSION } else { "0.2.0" }
$ld = "-s -w -X github.com/Fracizz/sshctl/cmd.Version=$Version"
$go = "go"
if (Test-Path "C:\Program Files\Go\bin\go.exe") {
    $env:GOROOT = "C:\Program Files\Go"
    $go = "C:\Program Files\Go\bin\go.exe"
}

if (Test-Path dist) {
    Remove-Item -Recurse -Force dist
}
$binDir = Join-Path dist "bin"
New-Item -ItemType Directory -Force -Path $binDir | Out-Null
New-Item -ItemType Directory -Force -Path bin | Out-Null

$targets = @(
    @{ GOOS = "linux"; GOARCH = "amd64"; Name = "sshctl-linux-amd64" },
    @{ GOOS = "linux"; GOARCH = "arm64"; Name = "sshctl-linux-arm64" },
    @{ GOOS = "windows"; GOARCH = "amd64"; Name = "sshctl-windows-amd64" },
    @{ GOOS = "windows"; GOARCH = "arm64"; Name = "sshctl-windows-arm64" },
    @{ GOOS = "darwin"; GOARCH = "amd64"; Name = "sshctl-darwin-amd64" },
    @{ GOOS = "darwin"; GOARCH = "arm64"; Name = "sshctl-darwin-arm64" }
)

foreach ($t in $targets) {
    $env:GOOS = $t.GOOS
    $env:GOARCH = $t.GOARCH
    $env:CGO_ENABLED = "0"
    $outName = if ($t.GOOS -eq "windows") { "$($t.Name).exe" } else { $t.Name }
    $outPath = Join-Path $binDir $outName
    Write-Host "building $outPath"
    & $go build -ldflags $ld -o $outPath .
}
Remove-Item Env:GOOS, Env:GOARCH -ErrorAction SilentlyContinue

foreach ($t in $targets) {
    $bundleName = $t.Name
    $bundleDir = Join-Path dist $bundleName
    if (Test-Path $bundleDir) {
        Remove-Item -Recurse -Force $bundleDir
    }
    New-Item -ItemType Directory -Force -Path $bundleDir | Out-Null

    $srcName = if ($t.GOOS -eq "windows") { "$bundleName.exe" } else { $bundleName }
    $srcPath = Join-Path $binDir $srcName
    if ($t.GOOS -eq "windows") {
        Copy-Item $srcPath (Join-Path $bundleDir "sshctl.exe")
    } else {
        Copy-Item $srcPath (Join-Path $bundleDir "sshctl")
    }

    $zipPath = Join-Path dist "$bundleName.zip"
    if (Test-Path $zipPath) {
        Remove-Item -Force $zipPath
    }
    Compress-Archive -Path (Join-Path $bundleDir "*") -DestinationPath $zipPath
    Write-Host "packaged $zipPath"
}

Copy-Item (Join-Path $binDir "sshctl-windows-amd64.exe") (Join-Path bin "sshctl.exe") -Force
Write-Host "done sshctl $Version"
