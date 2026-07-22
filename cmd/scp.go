package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/Fracizz/sshctl/internal/config"
	"github.com/Fracizz/sshctl/internal/sshx"
)

var scpTimeout time.Duration

var scpCmd = &cobra.Command{
	Use:   "scp <src> <dst>",
	Short: "Copy files via SFTP (scp-compatible paths: server:path)",
	Long: `Copy files between local and remote hosts.

Examples:
  sshctl scp ./a.txt example-host:/tmp/a.txt
  sshctl scp example-host:/etc/hosts ./hosts
  sshctl scp ./dir example-host:/tmp/dir
`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		src, dst := args[0], args[1]
		srcRemote := strings.Contains(src, ":") && !isWindowsDrive(src)
		dstRemote := strings.Contains(dst, ":") && !isWindowsDrive(dst)
		if srcRemote == dstRemote {
			return fmt.Errorf("one of src/dst must be remote (server:path)")
		}

		path := config.ResolvePath(cfgPath)
		f, err := config.Load(path)
		if err != nil {
			return err
		}

		if srcRemote {
			serverQuery, remotePath, err := splitRemote(src)
			if err != nil {
				return err
			}
			s, err := f.Find(serverQuery)
			if err != nil {
				return err
			}
			client, err := sshx.Dial(s, sshx.DialOptions{Timeout: scpTimeout, Insecure: insecure})
			if err != nil {
				return err
			}
			defer client.Close()
			return sshx.Download(client, remotePath, dst)
		}

		serverQuery, remotePath, err := splitRemote(dst)
		if err != nil {
			return err
		}
		s, err := f.Find(serverQuery)
		if err != nil {
			return err
		}
		client, err := sshx.Dial(s, sshx.DialOptions{Timeout: scpTimeout, Insecure: insecure})
		if err != nil {
			return err
		}
		defer client.Close()
		return sshx.Upload(client, src, remotePath)
	},
}

func init() {
	scpCmd.Flags().DurationVar(&scpTimeout, "timeout", 15*time.Second, "SSH dial timeout")
}

func isWindowsDrive(p string) bool {
	return len(p) >= 2 && p[1] == ':' && ((p[0] >= 'a' && p[0] <= 'z') || (p[0] >= 'A' && p[0] <= 'Z'))
}

func splitRemote(spec string) (server, remotePath string, err error) {
	// Prefer last colon so names/IPv6-ish forms still work with path.
	idx := strings.LastIndex(spec, ":")
	if idx <= 0 || idx == len(spec)-1 {
		return "", "", fmt.Errorf("invalid remote path %q (want server:path)", spec)
	}
	return spec[:idx], spec[idx+1:], nil
}
