package project

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"strconv"
)

func getProjects(baseUrl string, limit int) projects {
	url := fmt.Sprintf("%s/rest/api/1.0/projects/?limit=%d", baseUrl, limit)

	resp, err := http.Get(url)
	if err != nil {
		pterm.Error.Println(err.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var result projects
	if err := json.Unmarshal(body, &result); err != nil {
		pterm.Error.Println(err.Error())
		os.Exit(1)
	}

	return result
}

func printProjects(projects []Project) {
	var data [][]string

	data = append(data, []string{"ID", "Key", "Name", "Description"})

	for _, proj := range projects {
		row := []string{strconv.Itoa(proj.Id), proj.Key, proj.Name, proj.Description}
		data = append(data, row)
	}

	pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}

func listProjects(cmd *cobra.Command, args []string) {
	var baseUrl = viper.GetString("baseUrl")
	var limit = viper.GetInt("limit")

	var projects = getProjects(baseUrl, limit)

	printProjects(projects.Values)
	if !projects.IsLastPage {
		pterm.Warning.Println("Not all projects fetched, try with a higher limit")
	}
}
