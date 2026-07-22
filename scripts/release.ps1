# Cross-build release zips for GitHub Releases (not for daily local use).
$ErrorActionPreference = "Stop"
$Version = if ($env:VERSION) { $env:VERSION } else { "0.2.1" }
$ld = "-s -w -X github.com/Fracizz/sshctl/cmd.Version=$Version"
$go = "go"
if (Test-Path "C:\Program Files\Go\bin\go.exe") {
    $env:GOROOT = "C:\Program Files\Go"
    $go = "C:\Program Files\Go\bin\go.exe"
}

if (Test-Path dist) {
    Remove-Item -Recurse -Force dist
}
New-Item -ItemType Directory -Force -Path dist | Out-Null

$targets = @(
    @{ GOOS = "linux"; GOARCH = "amd64"; Name = "sshctl-linux-amd64" },
    @{ GOOS = "linux"; GOARCH = "arm64"; Name = "sshctl-linux-arm64" },
    @{ GOOS = "windows"; GOARCH = "amd64"; Name = "sshctl-windows-amd64" },
    @{ GOOS = "windows"; GOARCH = "arm64"; Name = "sshctl-windows-arm64" },
    @{ GOOS = "darwin"; GOARCH = "amd64"; Name = "sshctl-darwin-amd64" },
    @{ GOOS = "darwin"; GOARCH = "arm64"; Name = "sshctl-darwin-arm64" }
)

$stagingRoot = Join-Path ([System.IO.Path]::GetTempPath()) ("sshctl-release-" + [guid]::NewGuid().ToString("n"))
New-Item -ItemType Directory -Force -Path $stagingRoot | Out-Null
try {
    foreach ($t in $targets) {
        $env:GOOS = $t.GOOS
        $env:GOARCH = $t.GOARCH
        $env:CGO_ENABLED = "0"
        $releaseName = if ($t.GOOS -eq "windows") { "sshctl.exe" } else { "sshctl" }
        $stageDir = Join-Path $stagingRoot $t.Name
        New-Item -ItemType Directory -Force -Path $stageDir | Out-Null
        $stagePath = Join-Path $stageDir $releaseName
        Write-Host "building $stagePath"
        & $go build -ldflags $ld -o $stagePath .
        $zipPath = Join-Path dist "$($t.Name).zip"
        Compress-Archive -Path $stagePath -DestinationPath $zipPath
        Write-Host "packaged $zipPath"
    }
} finally {
    Remove-Item Env:GOOS, Env:GOARCH -ErrorAction SilentlyContinue
    if (Test-Path $stagingRoot) {
        Remove-Item -Recurse -Force $stagingRoot
    }
}

Write-Host "done release sshctl $Version -> dist/*.zip"
