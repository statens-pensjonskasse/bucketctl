package project

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
	"os"
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

	if !projects.IsLastPage {
		fmt.Fprintf(os.Stderr, "WARN: Not all projects fetched, try with a higher limit\n")
	}

	pkg.PrintProjects(projects.Values)
}

func init() {
}
