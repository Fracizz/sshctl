package main

import (
	"os"

	"github.com/Fracizz/invossh/cmd"
	"github.com/Fracizz/invossh/internal/exitcode"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(exitcode.ExecFailed)
	}
}
