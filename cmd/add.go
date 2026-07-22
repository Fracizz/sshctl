package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Fracizz/sshctl/internal/config"
)

var (
	addName        string
	addDescription string
	addHost        string
	addPort        int
	addUser        string
	addPassword    string
	addOS          string
	addKeyFile     string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a server (password is encrypted on save)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if addHost == "" || addUser == "" {
			return fmt.Errorf("--host and --user are required")
		}
		if addName == "" {
			addName = addHost
		}
		path := config.ResolvePath(cfgPath)
		var f *config.File
		if _, err := os.Stat(path); os.IsNotExist(err) {
			f = &config.File{Servers: []config.Server{}}
		} else if err != nil {
			return err
		} else {
			loaded, err := config.Load(path)
			if err != nil {
				return err
			}
			f = loaded
		}
		updated, err := f.Add(config.Server{
			Name:        addName,
			Description: addDescription,
			Host:        addHost,
			Port:        addPort,
			User:        addUser,
			Password:    addPassword,
			OS:          addOS,
			KeyFile:     addKeyFile,
		})
		if err != nil {
			return err
		}
		if err := config.Save(path, f); err != nil {
			return err
		}
		verb := "added"
		if updated {
			verb = "updated"
		}
		fmt.Printf("%s %s (%s@%s) -> %s\n", verb, addName, addUser, addHost, path)
		return nil
	},
}

func init() {
	addCmd.Flags().StringVar(&addName, "name", "", "display name")
	addCmd.Flags().StringVar(&addDescription, "desc", "", "server description")
	addCmd.Flags().StringVar(&addHost, "host", "", "hostname or IP")
	addCmd.Flags().IntVar(&addPort, "port", 22, "SSH port")
	addCmd.Flags().StringVar(&addUser, "user", "", "SSH username")
	addCmd.Flags().StringVar(&addPassword, "password", "", "SSH password (encrypted at rest)")
	addCmd.Flags().StringVar(&addOS, "os", "Linux", "OS label")
	addCmd.Flags().StringVar(&addKeyFile, "key", "", "optional private key path")
}
