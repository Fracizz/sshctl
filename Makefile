APP=sshctl
VERSION?=0.2.1
LDFLAGS=-s -w -X github.com/Fracizz/sshctl/cmd.Version=$(VERSION)

.PHONY: build tidy test clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(APP)$(shell go env GOEXE) .

tidy:
	go mod tidy

test:
	go test ./...

clean:
	rm -rf bin dist
