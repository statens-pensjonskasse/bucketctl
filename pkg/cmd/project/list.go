package project

import (
	"bucketctl/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
)

func prettyFormatProjects(projects map[string]*Project) [][]string {
	var data [][]string

	data = append(data, []string{"ID", "Key", "Name", "Description"})

	for key, proj := range projects {
		row := []string{strconv.Itoa(proj.Id), key, proj.Name, proj.Description}
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

	return pkg.PrintData(projects, prettyFormatProjects)
}
