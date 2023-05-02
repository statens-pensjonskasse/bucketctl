package project

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gobit/pkg"
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

type PermissionSet struct {
	Permissions map[string]*GivenPermissions
}

type ProjectPermissions struct {
	Project     string        `json:"project" yaml:"project,inline"`
	Permissions PermissionSet `json:"permissions" yaml:"permissions"`
}

type groupPermissions struct {
	pkg.BitbucketResponse
	Values []GroupPermission `json:"values"`
}

type userPermissions struct {
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

func GetProjectPermissions(baseUrl string, projectKey string, limit int, token string) (*PermissionSet, error) {
	permissionSet := &PermissionSet{
		Permissions: make(map[string]*GivenPermissions),
	}

	for _, permission := range PermissionTypes {
		permissionSet.Permissions[permission] = new(GivenPermissions)
	}

	projectGroupPermissions, err := getProjectGroupPermissions(baseUrl, projectKey, token, limit)
	if err != nil {
		return &PermissionSet{}, err
	}

	for _, gp := range projectGroupPermissions.Values {
		permissionSet.Permissions[gp.Permission].Groups = append(permissionSet.Permissions[gp.Permission].Groups, gp.Group.Name)
	}

	if !projectGroupPermissions.IsLastPage {
		pterm.Warning.Println("Not all Group Permissions fetched, try with a higher limit")
	}

	projectUserPermissions, err := getProjectUserPermissions(baseUrl, projectKey, token, limit)
	if err != nil {
		return &PermissionSet{}, err
	}

	for _, up := range projectUserPermissions.Values {
		permissionSet.Permissions[up.Permission].Users = append(permissionSet.Permissions[up.Permission].Users, up.User.Name)
	}

	if !projectUserPermissions.IsLastPage {
		pterm.Warning.Println("Not all User Permissions fetched, try with a higher limit")
	}

	return permissionSet, nil
}
