package repository

import (
	"bucketctl/pkg"
	"bucketctl/pkg/types"
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type Repository struct {
	Id            int    `json:"id,omitempty" yaml:"id,omitempty"`
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	HierarchyId   string `json:"hierarchyId,omitempty" yaml:"hierarchyId,omitempty"`
	ScmId         string `json:"scmId,omitempty" yaml:"scmId,omitempty"`
	State         string `json:"state,omitempty" yaml:"state,omitempty"`
	StatusMessage string `json:"statusMessage,omitempty" yaml:"statusMessage,omitempty"`
	Forkable      bool   `json:"forkable,omitempty" yaml:"forkable,omitempty"`
	Public        bool   `json:"public,omitempty" yaml:"public,omitempty"`
	Archived      bool   `json:"archived,omitempty" yaml:"archived,omitempty"`
}

var Cmd = &cobra.Command{
	Use:     "repository",
	Short:   "Repository commands",
	Aliases: []string{"repo"},
}

func init() {
	Cmd.AddCommand(listRepositoriesCmd)
}

func GetProjectRepositories(baseUrl string, projectKey string, limit int) (map[string]*Repository, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos?limit=%d", baseUrl, projectKey, limit)

	body, err := pkg.GetRequestBody(url, "")
	if err != nil {
		return nil, err
	}

	var repoResponse types.RepositoriesResponse
	if err := json.Unmarshal(body, &repoResponse); err != nil {
		return nil, err
	}

	if !repoResponse.IsLastPage {
		pterm.Warning.Println("Not all repositories fetched, try with a higher limit")
	}
	repositories := make(map[string]*Repository)
	for _, r := range repoResponse.Values {
		repositories[r.Slug] = &Repository{
			Id:            r.Id,
			Name:          r.Name,
			HierarchyId:   r.HierarchyId,
			ScmId:         r.ScmId,
			State:         r.State,
			StatusMessage: r.StatusMessage,
			Forkable:      r.Forkable,
			Public:        r.Public,
			Archived:      r.Archived,
		}
	}

	return repositories, nil
}
