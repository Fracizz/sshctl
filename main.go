package main

import (
	"os"

	"github.com/Fracizz/sshctl/cmd"
	"github.com/Fracizz/sshctl/internal/exitcode"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(exitcode.ExecFailed)
	}
}
