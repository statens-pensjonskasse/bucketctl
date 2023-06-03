package repository

import (
	"bucketctl/pkg"
	"bucketctl/pkg/types"
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "repository",
	Short:   "Repository commands",
	Aliases: []string{"repo"},
}

func init() {
	Cmd.AddCommand(listRepositoriesCmd)
}

func GetProjectRepositories(baseUrl string, projectKey string, limit int) ([]types.Repository, error) {
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

	return repoResponse.Values, nil
}
