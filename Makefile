APP=sshctl
VERSION?=0.2.0
LDFLAGS=-s -w -X github.com/Fracizz/sshctl/cmd.Version=$(VERSION)

.PHONY: build dist tidy test clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(APP)$(shell go env GOEXE) .

tidy:
	go mod tidy

test:
	go test ./...

dist: tidy
	@rm -rf dist && mkdir -p dist
	@tmpdir=$$(mktemp -d); \
	trap 'rm -rf "$$tmpdir"' EXIT; \
	for pair in \
		"linux:amd64:sshctl-linux-amd64" \
		"linux:arm64:sshctl-linux-arm64" \
		"windows:amd64:sshctl-windows-amd64" \
		"windows:arm64:sshctl-windows-arm64" \
		"darwin:amd64:sshctl-darwin-amd64" \
		"darwin:arm64:sshctl-darwin-arm64"; do \
		IFS=: read -r goos goarch name <<< "$$pair"; \
		stage="$$tmpdir/$$name"; \
		mkdir -p "$$stage"; \
		if [ "$$goos" = "windows" ]; then \
			out="$$stage/sshctl.exe"; \
		else \
			out="$$stage/sshctl"; \
		fi; \
		GOOS=$$goos GOARCH=$$goarch CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o "$$out" .; \
		( cd "$$stage" && zip -q "$(CURDIR)/dist/$$name.zip" "$$(basename "$$out")" ); \
	done

clean:
	rm -rf bin dist
