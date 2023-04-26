package project

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List Bitbucket projects",
	Run:     listProjects,
}

func listProjects(cmd *cobra.Command, args []string) {
	var baseUrl = viper.GetString("baseUrl")
	var limit = viper.GetInt("limit")

	var projects = pkg.GetProjects(baseUrl, limit)

	pkg.PrintProjects(projects.Values)
	if !projects.IsLastPage {
		pterm.Warning.Println("Not all projects fetched, try with a higher limit")
	}

}

func init() {
}
