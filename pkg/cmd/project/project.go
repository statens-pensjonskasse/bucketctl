package project

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gobit/pkg"
)

type Project struct {
	Id          int    `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type ProjectsResponse struct {
	pkg.BitbucketResponse
	Values []Project `json:"values"`
}

var Cmd = &cobra.Command{
	Use:     "project",
	Short:   "Bitbucket project commands",
	Aliases: []string{"proj"},
}

func init() {
	Cmd.AddCommand(listProjectsCmd)
}

var listProjectsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List Bitbucket projects",
	Run:     listProjects,
}

func GetProjects(baseUrl string, limit int) ([]Project, error) {
	url := fmt.Sprintf("%s/rest/api/1.0/projects/?limit=%d", baseUrl, limit)

	body, err := pkg.GetRequestBody(url, "")
	if err != nil {
		return []Project{}, err
	}

	var projectsResponse ProjectsResponse
	if err := json.Unmarshal(body, &projectsResponse); err != nil {
		return []Project{}, err
	}

	if !projectsResponse.IsLastPage {
		pterm.Warning.Println("Not all projects fetched, try with a higher limit")
	}

	return projectsResponse.Values, nil
}
