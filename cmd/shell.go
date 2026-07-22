package cmd

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/Fracizz/sshctl/internal/config"
	"github.com/Fracizz/sshctl/internal/sshx"
)

var shellTimeout time.Duration

var shellCmd = &cobra.Command{
	Use:   "shell <server>",
	Short: "Open an interactive SSH shell",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := config.ResolvePath(cfgPath)
		f, err := config.Load(path)
		if err != nil {
			return err
		}
		s, err := f.Find(args[0])
		if err != nil {
			return err
		}
		client, err := sshx.Dial(s, sshx.DialOptions{Timeout: shellTimeout, Insecure: insecure})
		if err != nil {
			return err
		}
		defer client.Close()
		return sshx.Shell(client)
	},
}

func init() {
	shellCmd.Flags().DurationVar(&shellTimeout, "timeout", 15*time.Second, "SSH dial timeout")
}
