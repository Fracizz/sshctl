package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Fracizz/sshctl/internal/config"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate legacy ~/.sshfrac or ~/.invossh inventory to ~/.sshctl",
	Long: `Copy servers.json from a legacy directory into ~/.sshctl/servers.json.

Legacy sources (first match wins):
  ~/.sshfrac/servers.json
  ~/.invossh/servers.json

The legacy file is renamed to servers.json.bak after success.
Migration also runs automatically before other commands when the primary file is missing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		from, err := config.MigrateLegacy()
		if err != nil {
			return err
		}
		dest := config.PrimaryConfigPath()
		if from == "" {
			fmt.Printf("nothing to migrate (primary already exists: %s)\n", dest)
			return nil
		}
		fmt.Printf("migrated %s -> %s\n", from, dest)
		fmt.Printf("legacy backed up to %s.bak\n", from)
		return nil
	},
}
