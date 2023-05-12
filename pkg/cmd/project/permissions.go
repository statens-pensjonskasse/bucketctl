package project

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gobit/pkg"
	"sort"
	"strings"
)

type Entities struct {
	Groups []string `json:"groups,omitempty" yaml:"groups,omitempty"`
	Users  []string `json:"users,omitempty" yaml:"users,omitempty"`
}

type PermissionSet struct {
	Permissions map[string]*Entities `json:"permissions" yaml:"permissions"`
}

type GrantedProjectPermissions struct {
	Project map[string]*PermissionSet `json:"" yaml:",inline"`
}

type groupPermissionsResponse struct {
	pkg.BitbucketResponse
	Values []GroupPermission `json:"values"`
}

type userPermissionsResponse struct {
	pkg.BitbucketResponse
	Values []UserPermission `json:"values"`
}

type GroupPermission struct {
	Group      pkg.Group `json:"group"`
	Permission string    `json:"permission"`
}

type UserPermission struct {
	User       pkg.User `json:"user"`
	Permission string   `json:"permission"`
}

var PermissionsCmd = &cobra.Command{
	Use:     "permissions",
	Short:   "Bitbucket project permission commands",
	Aliases: []string{"perm"},
}

func init() {
	PermissionsCmd.AddCommand(ListProjectPermissionsCmd)
	PermissionsCmd.AddCommand(listAllPermissionsCmd)
	PermissionsCmd.AddCommand(applyPermissionsFromFile)
}

func getProjectGroupPermissions(baseUrl string, projectKey string, limit int, token string) ([]GroupPermission, error) {
	groupPermissionsUrl := fmt.Sprintf("%s/rest/api/1.0/projects/%s/permissions/groups?limit=%d", baseUrl, projectKey, limit)

	body, err := pkg.GetRequestBody(groupPermissionsUrl, token)
	if err != nil {
		return []GroupPermission{}, err
	}

	var groups groupPermissionsResponse
	if err := json.Unmarshal(body, &groups); err != nil {
		return []GroupPermission{}, err
	}

	if !groups.IsLastPage {
		pterm.Warning.Println("Not all Group PermissionSet fetched, try with a higher limit")
	}

	return groups.Values, nil
}

func getProjectUserPermissions(baseUrl string, projectKey string, limit int, token string) ([]UserPermission, error) {
	userPermissionsUrl := fmt.Sprintf("%s/rest/api/1.0/projects/%s/permissions/users?limit=%d", baseUrl, projectKey, limit)

	body, err := pkg.GetRequestBody(userPermissionsUrl, token)
	if err != nil {
		return []UserPermission{}, err
	}

	var users userPermissionsResponse
	if err := json.Unmarshal(body, &users); err != nil {
		return []UserPermission{}, err
	}

	if !users.IsLastPage {
		pterm.Warning.Println("Not all User PermissionSet fetched, try with a higher limit")
	}

	return users.Values, nil
}

func GetProjectPermissions(baseUrl string, projectKey string, limit int, token string) (*GrantedProjectPermissions, error) {
	projectPermissions := &GrantedProjectPermissions{
		Project: map[string]*PermissionSet{},
	}
	projectPermissions.Project[projectKey] = new(PermissionSet)
	projectPermissions.Project[projectKey].Permissions = make(map[string]*Entities)

	projectGroupPermissions, err := getProjectGroupPermissions(baseUrl, projectKey, limit, token)
	if err != nil {
		return &GrantedProjectPermissions{}, err
	}

	for _, groupWithPermission := range projectGroupPermissions {
		if _, exists := projectPermissions.Project[projectKey].Permissions[groupWithPermission.Permission]; !exists {
			projectPermissions.Project[projectKey].Permissions[groupWithPermission.Permission] = new(Entities)
		}
		projectPermissions.Project[projectKey].Permissions[groupWithPermission.Permission].Groups = append(projectPermissions.Project[projectKey].Permissions[groupWithPermission.Permission].Groups, groupWithPermission.Group.Name)
	}

	projectUserPermissions, err := getProjectUserPermissions(baseUrl, projectKey, limit, token)
	if err != nil {
		return &GrantedProjectPermissions{}, err
	}

	for _, userWithPermission := range projectUserPermissions {
		if _, exists := projectPermissions.Project[projectKey].Permissions[userWithPermission.Permission]; !exists {
			projectPermissions.Project[projectKey].Permissions[userWithPermission.Permission] = new(Entities)
		}
		projectPermissions.Project[projectKey].Permissions[userWithPermission.Permission].Users = append(projectPermissions.Project[projectKey].Permissions[userWithPermission.Permission].Users, userWithPermission.User.Name)
	}

	return projectPermissions, nil
}

func PrettyFormatProjectPermissions(projectPermissions *GrantedProjectPermissions) [][]string {
	// Sorter prosjektene alfabetisk
	projects := make([]string, 0, len(projectPermissions.Project))
	for k := range projectPermissions.Project {
		projects = append(projects, k)
	}
	sort.Strings(projects)

	var data [][]string
	data = append(data, []string{"Project", "Permission", "Groups", "Users"})

	for _, k := range projects {
		proj := k
		for permission, v := range projectPermissions.Project[k].Permissions {
			var users string
			for _, user := range v.Users {
				users += user + "\n"
			}
			users = strings.Trim(users, "\n")
			var groups string
			for _, group := range v.Groups {
				groups += group + "\n"
			}
			groups = strings.Trim(groups, "\n")

			data = append(data, []string{proj, permission, groups, users})
			proj = ""
		}
	}
	return data
}
