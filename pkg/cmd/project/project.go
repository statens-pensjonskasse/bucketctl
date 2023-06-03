package project

import (
	"bucketctl/pkg"
	"bucketctl/pkg/types"
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "project",
	Short:   "Project commands",
	Aliases: []string{"proj"},
}

func init() {
	Cmd.AddCommand(listProjectsCmd)
}

var listProjectsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List projects",
	RunE:    listProjects,
}

func GetProjects(baseUrl string, limit int) ([]types.Project, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects?limit=%d", baseUrl, limit)

	body, err := pkg.GetRequestBody(url, "")
	if err != nil {
		return nil, err
	}

	var projectsResponse types.ProjectsResponse
	if err := json.Unmarshal(body, &projectsResponse); err != nil {
		return nil, err
	}

	if !projectsResponse.IsLastPage {
		pterm.Warning.Println("Not all projects fetched, try with a higher limit")
	}

	return projectsResponse.Values, nil
}
