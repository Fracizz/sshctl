package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Fracizz/sshfrac/internal/crypto"
)

var (
	cfgPath        string
	insecure       bool
	masterPassword string
	bindMachine    bool
	// Version is overwritten by -ldflags at build time.
	Version = "0.1.2"
)

var rootCmd = &cobra.Command{
	Use:   "sshfrac",
	Short: "AI-friendly SSH/SCP CLI with encrypted server inventory",
	Long: `sshfrac is a cross-platform SSH/SCP CLI designed primarily for AI agents.

Exit codes:
  0   success
  1   local runtime error (dial, decrypt, I/O, …)
  2   usage / config error
  N   remote command exit status (sshfrac exec only, when available)

Master password (recommended on shared machines):
  --master-password / SSHFRAC_MASTER_PASSWORD  → enc:v2 (Argon2id + AES-GCM)
  --bind-machine / SSHFRAC_BIND_MACHINE=1      → also bind v2 keys to this machine
  Without a master password, new secrets use legacy enc:v1 (machine-derived).

Shell completion:
  sshfrac completion bash|zsh|fish|powershell`,
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if masterPassword != "" {
			crypto.SetMasterPassword(masterPassword)
		}
		if cmd.Flags().Changed("bind-machine") {
			crypto.SetBindMachine(bindMachine)
		}
	},
}

// Execute runs the root command.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return err
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "path to servers JSON (default: ~/.sshfrac/servers.json or $SSHFRAC_CONFIG)")
	rootCmd.PersistentFlags().BoolVar(&insecure, "insecure", false, "skip SSH host key verification (unsafe; for lab only)")
	rootCmd.PersistentFlags().StringVar(&masterPassword, "master-password", "", "master password for enc:v2 (or set SSHFRAC_MASTER_PASSWORD)")
	rootCmd.PersistentFlags().BoolVar(&bindMachine, "bind-machine", false, "mix machine identity into enc:v2 KDF (or SSHFRAC_BIND_MACHINE=1)")

	rootCmd.SetVersionTemplate("sshfrac {{.Version}}\n")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(shellCmd)
	rootCmd.AddCommand(scpCmd)
	rootCmd.AddCommand(initCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sshfrac %s\n", Version)
	},
}
