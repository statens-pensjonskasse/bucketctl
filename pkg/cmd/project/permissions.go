package project

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

var (
	PermissionTypes = []string{
		"NONE",
		"PROJECT_READ",
		"REPO_CREATE",
		"PROJECT_WRITE",
		"PROJECT_ADMIN",
	}
)

type GivenPermissions struct {
	Groups []Group `yaml:"groups,omitempty"`
	Users  []User  `yaml:"users,omitempty"`
}

type PermissionSet struct {
	Permissions map[string]*GivenPermissions `xml:"permissions"`
}

type PSet map[string]*GivenPermissions

type Group struct {
	Name string `json:"name"`
}

type GroupPermission struct {
	Group      Group  `json:"group"`
	Permission string `json:"permission"`
}

type groupPermissions struct {
	pkg.BitbucketResponse
	Values []GroupPermission `json:"values"`
}

type User struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	Active       bool   `json:"active"`
	DisplayName  string `json:"displayName"`
	Id           int    `json:"id"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
}

type UserPermission struct {
	User       User   `json:"user"`
	Permission string `json:"permission"`
}

type userPermissions struct {
	pkg.BitbucketResponse
	Values []UserPermission `json:"values"`
}

func getProjectGroupPermissions(baseUrl string, projectKey string, token string, limit int) (groupPermissions, error) {
	groupPermissionsUrl := fmt.Sprintf("%s/rest/api/1.0/projects/%s/permissions/groups?limit=%d", baseUrl, projectKey, limit)

	body, err := pkg.GetRequestBody(groupPermissionsUrl, token)
	if err != nil {
		return groupPermissions{}, err
	}

	var groups groupPermissions
	if err := json.Unmarshal(body, &groups); err != nil {
		return groupPermissions{}, err
	}

	return groups, nil
}

func getProjectUserPermissions(baseUrl string, projectKey string, token string, limit int) (userPermissions, error) {
	userPermissionsUrl := fmt.Sprintf("%s/rest/api/1.0/projects/%s/permissions/users?limit=%d", baseUrl, projectKey, limit)

	body, err := pkg.GetRequestBody(userPermissionsUrl, token)
	if err != nil {
		return userPermissions{}, err
	}

	var users userPermissions
	if err := json.Unmarshal(body, &users); err != nil {
		return userPermissions{}, err
	}

	return users, nil
}

func printProjectPermissions(pSet *PermissionSet) {
	var data [][]string

	data = append(data, []string{"Permission", "Groups", "Users"})

	for permission, v := range pSet.Permissions {
		var users string
		for _, u := range v.Users {
			users += u.Name + "\n"
		}
		users = strings.Trim(users, "\n")
		var groups string
		for _, g := range v.Groups {
			groups += g.Name + "\n"
		}
		groups = strings.Trim(groups, "\n")

		if len(groups)+len(users) > 0 {
			data = append(data, []string{permission, groups, users})
		}
	}

	pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}

func listPermissions(cmd *cobra.Command, args []string) {
	var baseUrl = viper.GetString("baseUrl")
	var projectKey = viper.GetString("key")
	var token = viper.GetString("token")
	var limit = viper.GetInt("limit")

	projectGroupPermissions, err := getProjectGroupPermissions(baseUrl, projectKey, token, limit)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}

	projectUserPermissions, err := getProjectUserPermissions(baseUrl, projectKey, token, limit)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}

	pSet := &PermissionSet{
		Permissions: make(map[string]*GivenPermissions),
	}
	for _, permission := range PermissionTypes {
		pSet.Permissions[permission] = new(GivenPermissions)
	}

	for _, gp := range projectGroupPermissions.Values {
		pSet.Permissions[gp.Permission].Groups = append(pSet.Permissions[gp.Permission].Groups, gp.Group)
	}
	for _, up := range projectUserPermissions.Values {
		pSet.Permissions[up.Permission].Users = append(pSet.Permissions[up.Permission].Users, up.User)
	}

	printProjectPermissions(pSet)

	yamlData, err := yaml.Marshal(&pSet)
	if err != nil {
		pterm.Error.Println("Error while Marshaling. %v", err)
	}
	pterm.Println(string(yamlData))

	jsonData, err := json.MarshalIndent(&pSet, "", "  ")
	if err != nil {
		pterm.Error.Println("Error while Marshaling. %v", err)
	}
	pterm.Println(string(jsonData))

	if !projectGroupPermissions.IsLastPage {
		pterm.Warning.Println("Not all groupPermissions fetched, try with a higher limit")
	}

	if !projectUserPermissions.IsLastPage {
		pterm.Warning.Println("Not all projectUserPermissions fetched, try with a higher limit")
	}
}
