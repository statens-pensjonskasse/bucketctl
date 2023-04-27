package project

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

type projects struct {
	pkg.BitbucketResponse
	Values []Project `json:"values"`
}

var (
	key string
)

var Cmd = &cobra.Command{
	Use:     "project",
	Short:   "Bitbucket project commands",
	Aliases: []string{"proj"},
}

func init() {
	Cmd.PersistentFlags().StringVarP(&key, "key", "k", "", "Project key")
	viper.BindPFlag("key", Cmd.PersistentFlags().Lookup("key"))

	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(listPermissionsCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List Bitbucket projects",
	Run:     listProjects,
}

var listPermissionsCmd = &cobra.Command{
	Use: "permissions",
	Run: listPermissions,
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
		pterm.Error.Println(err.Error())
		os.Exit(1)
	}

	return result, nil
}
