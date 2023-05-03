package project

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
	"os"
	"strconv"
)

func prettyFormatProjects(projects []Project) [][]string {
	var data [][]string

	data = append(data, []string{"ID", "Key", "Name", "Description"})

	for _, proj := range projects {
		row := []string{strconv.Itoa(proj.Id), proj.Key, proj.Name, proj.Description}
		data = append(data, row)
	}

	return data
}

func listProjects(cmd *cobra.Command, args []string) {
	var baseUrl = viper.GetString("baseUrl")
	var limit = viper.GetInt("limit")

	projects, err := GetProjects(baseUrl, limit)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}

	pkg.PrintData(projects, prettyFormatProjects)
}
