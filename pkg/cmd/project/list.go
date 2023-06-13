package project

import (
	"bucketctl/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sort"
	"strconv"
)

func prettyFormatProjects(projectsMap map[string]*Project) [][]string {
	projects := make([]string, 0, len(projectsMap))
	for p := range projectsMap {
		projects = append(projects, p)
	}
	sort.Strings(projects)

	var data [][]string
	data = append(data, []string{"ID", "Key", "Name", "Description"})
	for _, key := range projects {
		row := []string{strconv.Itoa(projectsMap[key].Id), key, projectsMap[key].Name, projectsMap[key].Description}
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
