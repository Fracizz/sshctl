package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Fracizz/sshfrac/internal/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create starter config at ~/.sshfrac/servers.json (no real secrets)",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := config.ResolvePath(cfgPath)
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s already exists", path)
		}
		out := &config.File{}
		if err := out.Add(config.Server{
			Name:        "example-host",
			Description: "示例：请 delete 后用 sshfrac add 写入真实主机（勿保留 REPLACE_ME）",
			Host:        "192.0.2.10",
			Port:        22,
			User:        "root",
			Password:    "REPLACE_ME",
			OS:          "Linux",
		}); err != nil {
			return err
		}
		if err := config.Save(path, out); err != nil {
			return err
		}
		fmt.Printf("wrote %s\n", path)
		fmt.Println("next: sshfrac add --host <ip> --user root --password <secret> --desc \"...\"")
		fmt.Println("      then remove the example-host entry from the JSON if unused")
		return nil
	},
}
