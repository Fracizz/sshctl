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
	@mkdir -p dist
	GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(APP)-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(APP)-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(APP)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(APP)-windows-arm64.exe .
	GOOS=darwin  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(APP)-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(APP)-darwin-arm64 .

clean:
	rm -rf bin dist
