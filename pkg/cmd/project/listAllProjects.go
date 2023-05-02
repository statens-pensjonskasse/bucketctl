package project

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
	"os"
	"strconv"
)

func GetProjects(baseUrl string, limit int) (Projects, error) {
	url := fmt.Sprintf("%s/rest/api/1.0/projects/?limit=%d", baseUrl, limit)

	body, err := pkg.GetRequestBody(url, "")
	if err != nil {
		return Projects{}, err
	}

	var result Projects
	if err := json.Unmarshal(body, &result); err != nil {
		return Projects{}, err
	}

	return result, nil
}

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

	pkg.PrintData(projects.Values, prettyFormatProjects)

	if !projects.IsLastPage {
		pterm.Warning.Println("Not all projects fetched, try with a higher limit")
	}
}
