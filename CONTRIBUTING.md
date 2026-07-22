# Contributing

Thanks for contributing to sshctl.

## Development

```bash
go test ./...
go vet ./...
go build -o bin/sshctl .
```

Optional lint (same as CI):

```bash
golangci-lint run
```

Local build (single binary for your platform):

```bash
make build         # bin/sshctl (or bin/sshctl.exe on Windows)
# Windows:
powershell -File scripts/build.ps1   # bin/sshctl.exe + skills/sshctl/bin/sshctl.exe
```

Release binaries for all platforms (linux/windows/darwin × amd64/arm64) are built in [CI](.github/workflows/ci.yml) on push to `main` and attached as workflow artifacts. Download release zips from [GitHub Releases](https://github.com/Fracizz/sshctl/releases), not from a local `dist/` folder.

## Pull requests

- Keep changes focused.
- Add/adjust tests for config, search, and crypto when behavior changes.
- Never commit `~/.sshctl/servers.json`, `~/.sshfrac/servers.json`, `bin/`, `dist/`, or real credentials.
- Prefer English for new CLI help strings; README may stay bilingual.
- Run `go test ./...` before opening a PR.

## Code style

`gofmt` / `go fmt ./...`.
