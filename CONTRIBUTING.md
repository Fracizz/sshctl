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

Cross-build:

```bash
make dist
# Windows:
powershell -File scripts/build.ps1
```

## Pull requests

- Keep changes focused.
- Add/adjust tests for config, search, and crypto when behavior changes.
- Never commit `~/.sshctl/servers.json`, `bin/`, `dist/`, or real credentials.
- Prefer English for new CLI help strings; README may stay bilingual.
- Run `go test ./...` before opening a PR.

## Code style

`gofmt` / `go fmt ./...`.
