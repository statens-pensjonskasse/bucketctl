package repository

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

var (
	key string
)

var Cmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("key", cmd.PersistentFlags().Lookup("key"))
	},
	Use:     "repository",
	Short:   "Bitbucket repository commands",
	Aliases: []string{"repo"},
}

func init() {
	Cmd.PersistentFlags().StringVarP(&key, "key", "k", "", "Project key")
	Cmd.MarkPersistentFlagRequired("key")
	Cmd.AddCommand(listRepositoriesCmd)
}

var listRepositoriesCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("key", cmd.PersistentFlags().Lookup("key"))
	},
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List Bitbucket repositories in a given project",
	RunE:    listRepositories,
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
