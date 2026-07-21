package main

import (
	"os"

	"github.com/Fracizz/sshfrac/cmd"
	"github.com/Fracizz/sshfrac/internal/exitcode"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(exitcode.ExecFailed)
	}
}
