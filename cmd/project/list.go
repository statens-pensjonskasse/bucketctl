package project

import (
	"github.com/spf13/cobra"
	"gobit/pkg"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List Bitbucket projects",
	Run:     listProjects,
}

func listProjects(cmd *cobra.Command, args []string) {
	var baseUrl, _ = cmd.Flags().GetString("baseUrl")
	var limit, _ = cmd.Flags().GetInt("limit")

	var projects = pkg.GetProjects(baseUrl, limit)
	pkg.PrintProjects(projects.Values)
}

func init() {
}
