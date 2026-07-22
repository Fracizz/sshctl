# Contributing

Thanks for contributing to sshctl.

## Development

```bash
go test ./...
go vet ./...
go build -o bin/sshctl .
```

Windows local build (amd64 only; syncs skill bins):

```powershell
powershell -File scripts/build.ps1
```

Optional lint (same as CI):

```bash
golangci-lint run
```

Multi-platform release zips (GitHub Releases only):

```powershell
powershell -File scripts/release.ps1   # dist/sshctl-*.zip
```

CI builds per-platform binaries as workflow artifacts on push to `main`.

## Pull requests

- Keep changes focused.
- Add/adjust tests for config, search, and crypto when behavior changes.
- Never commit `~/.sshctl/servers.json`, `~/.sshfrac/servers.json`, `bin/`, `dist/`, `skills/*/bin/`, or real credentials.
- Prefer English for new CLI help strings; README may stay bilingual.
- Run `go test ./...` before opening a PR.

## Code style

`gofmt` / `go fmt ./...`.
