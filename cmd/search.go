package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/Fracizz/invossh/internal/config"
)

var searchQuery string

var searchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"find"},
	Short:   "Case-insensitive contains search on name / IP / description",
	Example: `  invossh search -s 157
  invossh search --search 客制化
  invossh search -s 212`,
	RunE: func(cmd *cobra.Command, args []string) error {
		q := searchQuery
		if q == "" && len(args) > 0 {
			q = args[0]
		}
		if q == "" {
			return fmt.Errorf("missing query; use: invossh search -s <keyword>")
		}
		path := config.ResolvePath(cfgPath)
		f, err := config.Load(path)
		if err != nil {
			return err
		}
		hits := f.Search(q)
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tUSER\tHOST\tPORT\tOS\tDESCRIPTION")
		for _, s := range hits {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\n", s.Name, s.User, s.Host, s.Port, s.OS, s.Description)
		}
		if err := w.Flush(); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "# %d match(es) for %q\n", len(hits), q)
		return nil
	},
}

func init() {
	searchCmd.Flags().StringVarP(&searchQuery, "search", "s", "", "keyword matched against name / host(IP) / description (case-insensitive contains)")
}
