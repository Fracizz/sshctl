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
New-Item -ItemType Directory -Force -Path dist | Out-Null
New-Item -ItemType Directory -Force -Path bin | Out-Null
$skillBinDir = Join-Path "skills" (Join-Path "sshctl" "bin")
$skillRoot = Join-Path "skills" "sshctl"
$hasSkillDir = Test-Path $skillRoot
if ($hasSkillDir) {
    New-Item -ItemType Directory -Force -Path $skillBinDir | Out-Null
}

$targets = @(
    @{ GOOS = "linux"; GOARCH = "amd64"; Name = "sshctl-linux-amd64" },
    @{ GOOS = "linux"; GOARCH = "arm64"; Name = "sshctl-linux-arm64" },
    @{ GOOS = "windows"; GOARCH = "amd64"; Name = "sshctl-windows-amd64" },
    @{ GOOS = "windows"; GOARCH = "arm64"; Name = "sshctl-windows-arm64" },
    @{ GOOS = "darwin"; GOARCH = "amd64"; Name = "sshctl-darwin-amd64" },
    @{ GOOS = "darwin"; GOARCH = "arm64"; Name = "sshctl-darwin-arm64" }
)

$stagingRoot = Join-Path ([System.IO.Path]::GetTempPath()) ("sshctl-build-" + [guid]::NewGuid().ToString("n"))
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
        if (Test-Path $zipPath) {
            Remove-Item -Force $zipPath
        }
        Compress-Archive -Path $stagePath -DestinationPath $zipPath
        Write-Host "packaged $zipPath"

        if ($t.GOOS -eq "windows" -and $t.GOARCH -eq "amd64") {
            Copy-Item $stagePath (Join-Path bin "sshctl.exe") -Force
            if ($hasSkillDir) {
                Copy-Item $stagePath (Join-Path $skillBinDir "sshctl.exe") -Force
            }
        }
    }
} finally {
    Remove-Item Env:GOOS, Env:GOARCH -ErrorAction SilentlyContinue
    if (Test-Path $stagingRoot) {
        Remove-Item -Recurse -Force $stagingRoot
    }
}

if ($hasSkillDir) {
    Write-Host "done sshctl $Version (dist/*.zip, bin/sshctl.exe, skills/sshctl/bin/sshctl.exe)"
} else {
    Write-Host "done sshctl $Version (dist/*.zip only)"
}
