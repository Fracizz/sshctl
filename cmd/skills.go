package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/Fracizz/sshctl/internal/skills"
)

var skillsSearch string

var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "List installed AI agent skills",
	Long: `Scan common Agent skill directories and list skills that contain SKILL.md.

Roots (when present):
  ~/.claude/skills
  ~/.cursor/skills
  ~/.codex/skills
  parent of skill next to this binary (.../bin/sshctl)`,
	Example: `  sshctl skills
  sshctl skills -s sshctl`,
	RunE: func(cmd *cobra.Command, args []string) error {
		exe, err := os.Executable()
		if err != nil {
			exe = ""
		}
		list, err := skills.Discover(exe)
		if err != nil {
			return err
		}
		list = skills.Filter(list, skillsSearch)
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "AGENT\tNAME\tPATH\tBIN")
		for _, s := range list {
			bin := "no"
			if s.HasBin {
				bin = "yes"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.Agent, s.Name, s.Path, bin)
		}
		return w.Flush()
	},
}

func init() {
	skillsCmd.Flags().StringVarP(&skillsSearch, "search", "s", "", "filter by name / directory / description (case-insensitive contains)")
}
