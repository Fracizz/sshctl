# Local build: Windows amd64 only. Syncs into skill bins.
# Multi-platform release zips: push a v* tag → GitHub Actions (.github/workflows/release.yml).
$ErrorActionPreference = "Stop"
$Version = if ($env:VERSION) { $env:VERSION } else { "0.2.3" }
$ld = "-s -w -X github.com/Fracizz/sshctl/cmd.Version=$Version"
$go = "go"
if (Test-Path "C:\Program Files\Go\bin\go.exe") {
    $env:GOROOT = "C:\Program Files\Go"
    $go = "C:\Program Files\Go\bin\go.exe"
}

New-Item -ItemType Directory -Force -Path bin | Out-Null
$out = Join-Path bin "sshctl.exe"
Write-Host "building $out"
& $go build -ldflags $ld -o $out .

$targets = @(
    (Join-Path "skills" (Join-Path "sshctl" "bin")),
    (Join-Path $env:USERPROFILE ".claude\skills\sshctl\bin"),
    (Join-Path $env:USERPROFILE ".cursor\skills\sshctl\bin"),
    (Join-Path $env:USERPROFILE ".codex\skills\sshctl\bin")
)
foreach ($binDir in $targets) {
    $skillRoot = Split-Path $binDir -Parent
    if (-not (Test-Path $skillRoot)) {
        continue
    }
    New-Item -ItemType Directory -Force -Path $binDir | Out-Null
    $dest = Join-Path $binDir "sshctl.exe"
    Copy-Item $out $dest -Force
    Write-Host "copied $dest"
}

Write-Host "done sshctl $Version"
