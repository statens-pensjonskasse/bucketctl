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

type Group struct {
	Name string `json:"name"`
}

type Permission struct {
	Permission string `json:"permission"`
	Group      Group  `json:"group"`
}

type permissions struct {
	pkg.Response
	Values []Permission `json:"values"`
}

func getProjectPermissions(baseUrl string, projectKey string, token string, limit int) permissions {
	url := fmt.Sprintf("%s/rest/api/1.0/projects/%s/permissions/groups?limit=%d", baseUrl, projectKey, limit)

	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		pterm.Error.Println(err.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var result permissions
	if err := json.Unmarshal(body, &result); err != nil {
		pterm.Error.Println(err.Error())
		os.Exit(1)
	}

	return result
}

func printPermissions(permissions []Permission) {
	var data [][]string

	data = append(data, []string{"Permission", "Group"})

	for _, perm := range permissions {
		row := []string{perm.Permission, perm.Group.Name}
		data = append(data, row)
	}

	pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}

func listPermissions(cmd *cobra.Command, args []string) {
	var baseUrl = viper.GetString("baseUrl")
	var projectKey = viper.GetString("key")
	var token = viper.GetString("token")
	var limit = viper.GetInt("limit")

	permissions := getProjectPermissions(baseUrl, projectKey, token, limit)

	printPermissions(permissions.Values)
	if !permissions.IsLastPage {
		pterm.Warning.Println("Not all permissions fetched, try with a higher limit")
	}
}
