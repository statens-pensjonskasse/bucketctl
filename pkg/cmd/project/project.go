package project

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gobit/pkg"
	"io"
	"net/http"
	"os"
)

type Project struct {
	Id          int    `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type Projects struct {
	pkg.BitbucketResponse
	Values []Project `json:"values"`
}

var Cmd = &cobra.Command{
	Use:     "project",
	Short:   "Bitbucket project commands",
	Aliases: []string{"proj"},
}

func init() {
	Cmd.AddCommand(PermissionsCmd)
	Cmd.AddCommand(listProjectsCmd)
}

var listProjectsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List Bitbucket projects",
	Run:     listProjects,
}

func getProject(baseUrl string, projectKey string, limit int) (Project, error) {
	url := fmt.Sprintf("%s/rest/api/1.0/projects/%s/?limit=%d", baseUrl, projectKey, limit)

	resp, err := http.Get(url)
	if err != nil {
		return Project{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var result Project
	if err := json.Unmarshal(body, &result); err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}

	return result, nil
}
