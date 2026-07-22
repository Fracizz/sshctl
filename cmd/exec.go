package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/Fracizz/sshctl/internal/config"
	"github.com/Fracizz/sshctl/internal/shellquote"
	"github.com/Fracizz/sshctl/internal/sshx"
)

var execTimeout time.Duration

var execCmd = &cobra.Command{
	Use:   "exec <server> [--] <command...>",
	Short: "Run a remote command over SSH",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverQuery := args[0]
		remoteArgs := args[1:]
		if dash := cmd.ArgsLenAtDash(); dash >= 0 {
			remoteArgs = args[dash:]
			if dash == 0 {
				return fmt.Errorf("missing server name")
			}
			serverQuery = args[0]
		}
		if len(remoteArgs) == 0 {
			return fmt.Errorf("missing remote command; use: sshctl exec <server> -- <cmd>")
		}
		path := config.ResolvePath(cfgPath)
		f, err := config.Load(path)
		if err != nil {
			return err
		}
		s, err := f.Find(serverQuery)
		if err != nil {
			return err
		}
		client, err := sshx.Dial(s, sshx.DialOptions{Timeout: execTimeout, Insecure: insecure})
		if err != nil {
			return err
		}
		defer client.Close()
		// 多参数必须 shell quote，否则 bash -lc "cd x && y" 会被拼成 bash -lc cd x && y
		code, err := sshx.Run(client, shellquote.JoinRemoteCommand(remoteArgs))
		if err != nil {
			return err
		}
		if code != 0 {
			os.Exit(code)
		}
		return nil
	},
}

func init() {
	execCmd.Flags().DurationVar(&execTimeout, "timeout", 15*time.Second, "SSH dial timeout")
}
