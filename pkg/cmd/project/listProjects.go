package project

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
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

func listProjects(cmd *cobra.Command, args []string) error {
	var baseUrl = viper.GetString("baseUrl")
	var limit = viper.GetInt("limit")

	projects, err := GetProjects(baseUrl, limit)
	if err != nil {
		return err
	}

	pkg.PrintData(projects, prettyFormatProjects)
	return nil
}
