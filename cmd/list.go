package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/Fracizz/invossh/internal/config"
)

var listSearch string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List servers from the config file",
	Example: `  invossh list
  invossh list -s 157
  invossh list --search 客制化`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := config.ResolvePath(cfgPath)
		f, err := config.Load(path)
		if err != nil {
			return err
		}
		servers := f.Servers
		if listSearch != "" {
			hits := f.Search(listSearch)
			servers = make([]config.Server, 0, len(hits))
			for _, s := range hits {
				servers = append(servers, *s)
			}
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tUSER\tHOST\tPORT\tOS\tDESCRIPTION")
		for _, s := range servers {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\n", s.Name, s.User, s.Host, s.Port, s.OS, s.Description)
		}
		return w.Flush()
	},
}

func init() {
	listCmd.Flags().StringVarP(&listSearch, "search", "s", "", "filter by name / host(IP) / description (case-insensitive contains)")
}
