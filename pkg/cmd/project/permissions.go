package project

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
	"os"
	"strings"
)

type Group struct {
	Name string `json:"name"`
}

type Perm int

const (
	PROJECT_ADMIN = iota
	PROJECT_WRITE
	PROJECT_READ
)

type PermissionSet map[string][]string

type Permission struct {
	Permission string `json:"permission"`
	Group      Group  `json:"group"`
}

type permissions struct {
	pkg.BitbucketResponse
	Values []Permission `json:"values"`
}

func getProjectPermissions(baseUrl string, projectKey string, token string, limit int) (permissions, error) {
	url := fmt.Sprintf("%s/rest/api/1.0/projects/%s/permissions/groups?limit=%d", baseUrl, projectKey, limit)

	body, err := pkg.GetRequestBody(url, token)
	if err != nil {
		return permissions{}, err
	}

	var result permissions
	if err := json.Unmarshal(body, &result); err != nil {
		return permissions{}, err
	}

	return result, nil
}

func printPermissions(permissions []Permission) {
	var data [][]string

	data = append(data, []string{"Permission", "Groups"})

	groupedPermissions := make(PermissionSet)

	for _, p := range permissions {
		groupedPermissions[p.Permission] = append(groupedPermissions[p.Permission], p.Group.Name)
	}

	for key, _ := range groupedPermissions {
		var groups string
		for _, g := range groupedPermissions[key] {
			groups += g + "\n"
		}
		row := []string{key, strings.Trim(groups, "\n")}
		data = append(data, row)
	}

	pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}

func listPermissions(cmd *cobra.Command, args []string) {
	var baseUrl = viper.GetString("baseUrl")
	var projectKey = viper.GetString("key")
	var token = viper.GetString("token")
	var limit = viper.GetInt("limit")

	permissions, err := getProjectPermissions(baseUrl, projectKey, token, limit)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}

	printPermissions(permissions.Values)

	if !permissions.IsLastPage {
		pterm.Warning.Println("Not all permissions fetched, try with a higher limit")
	}
}
