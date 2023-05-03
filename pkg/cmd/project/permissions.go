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

var (
	PermissionTypes = []string{
		"PROJECT_READ",
		"REPO_CREATE",
		"PROJECT_WRITE",
		"PROJECT_ADMIN",
	}
)

type GivenPermissions struct {
	Groups []string `json:"groups,omitempty" yaml:"groups,omitempty"`
	Users  []string `json:"users,omitempty" yaml:"users,omitempty"`
}

type PermissionObjects struct {
	Permissions map[string]*GivenPermissions `json:"permissions" yaml:"permissions"`
}

type ProjectPermissions struct {
	Project map[string]*PermissionObjects `json:"" yaml:",inline"`
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
		pterm.Warning.Println("Not all Group Permissions fetched, try with a higher limit")
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
		pterm.Warning.Println("Not all User Permissions fetched, try with a higher limit")
	}

	return users.Values, nil
}

func GetProjectPermissions(baseUrl string, projectKey string, limit int, token string) (*ProjectPermissions, error) {
	projectPermissions := &ProjectPermissions{
		Project: map[string]*PermissionObjects{},
	}
	projectPermissions.Project[projectKey] = new(PermissionObjects)
	projectPermissions.Project[projectKey].Permissions = make(map[string]*GivenPermissions)

	for _, permission := range PermissionTypes {
		projectPermissions.Project[projectKey].Permissions[permission] = new(GivenPermissions)
	}

	projectGroupPermissions, err := getProjectGroupPermissions(baseUrl, projectKey, limit, token)
	if err != nil {
		return &ProjectPermissions{}, err
	}

	for _, gp := range projectGroupPermissions {
		projectPermissions.Project[projectKey].Permissions[gp.Permission].Groups = append(projectPermissions.Project[projectKey].Permissions[gp.Permission].Groups, gp.Group.Name)
	}

	projectUserPermissions, err := getProjectUserPermissions(baseUrl, projectKey, limit, token)
	if err != nil {
		return &ProjectPermissions{}, err
	}

	for _, up := range projectUserPermissions {
		projectPermissions.Project[projectKey].Permissions[up.Permission].Users = append(projectPermissions.Project[projectKey].Permissions[up.Permission].Users, up.User.Name)
	}

	return projectPermissions, nil
}

func PrettyFormatProjectPermissions(projectPermissions *ProjectPermissions) [][]string {
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

			// Dersom verken en gruppe eller en bruker har rettigheten så hopper vi over den
			if len(groups)+len(users) > 0 {
				data = append(data, []string{proj, permission, groups, users})
				proj = ""
			}
		}
	}
	return data
}
