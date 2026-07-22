package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Fracizz/sshctl/internal/install"
)

var (
	installDir  string
	installPath bool
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install sshctl to a system directory and add it to PATH",
	Long: `Copy the current sshctl binary to a system directory.

Windows default: C:\Program Files\sshctl\sshctl.exe (requires Administrator)
Also adds that directory to the machine PATH when --path is set (default).

After install, open a new terminal and run: sshctl version`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dest, err := install.Install(installDir, installPath)
		if err != nil {
			return err
		}
		fmt.Printf("installed %s\n", dest)
		if installPath {
			fmt.Println("PATH updated (machine). Open a new terminal, then: sshctl version")
		}
		return nil
	},
}

func init() {
	installCmd.Flags().StringVar(&installDir, "dir", "", "install directory (default: "+install.DefaultDir()+")")
	installCmd.Flags().BoolVar(&installPath, "path", true, "add install directory to machine PATH (Windows)")
}
