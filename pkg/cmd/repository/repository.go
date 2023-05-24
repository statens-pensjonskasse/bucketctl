package repository

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gobit/pkg"
	"gobit/pkg/cmd/project"
)

type Repository struct {
	Id            int             `json:"id"`
	Name          string          `json:"name"`
	Slug          string          `json:"slug"`
	HierarchyId   string          `json:"hierarchyId"`
	ScmId         string          `json:"scmId"`
	State         string          `json:"state"`
	StatusMessage string          `json:"statusMessage"`
	Forkable      bool            `json:"forkable"`
	Public        bool            `json:"public"`
	Archived      bool            `json:"archived"`
	Project       project.Project `json:"project"`
}

type RepositoriesResponse struct {
	pkg.BitbucketResponse
	Values []Repository `json:"values"`
}

var Cmd = &cobra.Command{
	Use:     "repository",
	Short:   "Bitbucket repository commands",
	Aliases: []string{"repo"},
}

func init() {
	Cmd.AddCommand(listRepositoriesCmd)
}

func getRepositories(baseUrl string, projectKey string, limit int) ([]Repository, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos?limit=%d", baseUrl, projectKey, limit)

	body, err := pkg.GetRequestBody(url, "")
	if err != nil {
		return nil, err
	}

	var repoResponse RepositoriesResponse
	if err := json.Unmarshal(body, &repoResponse); err != nil {
		return nil, err
	}

	if !repoResponse.IsLastPage {
		pterm.Warning.Println("Not all projects fetched, try with a higher limit")
	}

	return repoResponse.Values, nil
}
