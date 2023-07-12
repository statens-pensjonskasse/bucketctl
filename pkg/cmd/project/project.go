package project

import (
	"bucketctl/pkg"
	"bucketctl/pkg/types"
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"sort"
	"strconv"
)

type Project struct {
	Id          int    `json:"id,omitempty" yaml:"id"`
	Name        string `json:"name,omitempty" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description"`
	Public      bool   `json:"public,omitempty" yaml:"public"`
}

var Cmd = &cobra.Command{
	Use:     "project",
	Short:   "Project commands",
	Aliases: []string{"proj"},
}

func init() {
	Cmd.AddCommand(listProjectsCmd)
}

func GetProjects(baseUrl string, limit int) (map[string]*Project, error) {
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

	projects := make(map[string]*Project)
	for _, p := range projectsResponse.Values {
		projects[p.Key] = &Project{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Public:      p.Public,
		}
	}

	return projects, nil
}

func prettyFormatProjects(projectsMap map[string]*Project) [][]string {
	projects := make([]string, 0, len(projectsMap))
	for p := range projectsMap {
		projects = append(projects, p)
	}
	sort.Strings(projects)

	var data [][]string
	data = append(data, []string{"ID", types.ProjectKeyFlag, "Name", "Description"})
	for _, key := range projects {
		row := []string{strconv.Itoa(projectsMap[key].Id), key, projectsMap[key].Name, projectsMap[key].Description}
		data = append(data, row)
	}

	return data
}
