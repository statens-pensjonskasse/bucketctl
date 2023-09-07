package permission

import (
	project2 "bucketctl/pkg/cmd/project"
	"bucketctl/pkg/cmd/repository"
	"bucketctl/pkg/common"
	"bucketctl/pkg/types"
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"strings"
)

type Entities struct {
	Groups []string `json:"groups,omitempty" yaml:"groups,omitempty"`
	Users  []string `json:"users,omitempty" yaml:"users,omitempty"`
}

type Permissions map[string]*Entities

type ProjectPermissions struct {
	Public            bool                              `json:"public" yaml:"public"`
	DefaultPermission string                            `json:"default-permission" yaml:"default-permission"`
	Permissions       *Permissions                      `json:"permissions,omitempty" yaml:"permissions,omitempty"`
	Repositories      map[string]*RepositoryPermissions `json:"repositories,omitempty" yaml:"repositories,omitempty"`
}

type RepositoryPermissions struct {
	Permissions *Permissions `json:"permissions" yaml:"permissions"`
}

var Cmd = &cobra.Command{
	Use:     "permission",
	Short:   "Permission commands",
	Aliases: []string{"perm"},
}

func init() {
	Cmd.AddCommand(applyPermissionsCmd)
	Cmd.AddCommand(listAllPermissionsCmd)
	Cmd.AddCommand(listPermissionsCmd)
}

func isProjectPublic(baseUrl string, projectKey string, token string) (bool, error) {
	project, err := project2.GetProject(baseUrl, projectKey, token)
	if err != nil {
		return false, err
	}
	return project.Public, err
}

func getDefaultProjectPermission(baseUrl string, projectKey string, token string) (string, error) {
	for _, permission := range []string{"PROJECT_ADMIN", "PROJECT_WRITE", "PROJECT_READ"} {
		url := fmt.Sprintf("%s/rest/api/latest/projects/%s/permissions/%s/all", baseUrl, projectKey, permission)
		body, err := common.GetRequestBody(url, token)
		if err != nil {
			return "", err
		}
		var defaultProjectPermission types.DefaultProjectPermission
		if err := json.Unmarshal(body, &defaultProjectPermission); err != nil {
			return "", err
		}
		if defaultProjectPermission.Permitted {
			return permission, nil
		}
	}
	return "", nil
}

func getGroupPermissions(url string, token string) ([]*types.GroupPermission, error) {
	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var groups types.GroupPermissionsResponse
	if err := json.Unmarshal(body, &groups); err != nil {
		return nil, err
	}

	if !groups.IsLastPage {
		pterm.Warning.Println("Not all Group permissions fetched, try with a higher limit")
	}

	return groups.Values, nil
}

func getProjectGroupPermissions(baseUrl string, projectKey string, limit int, token string) ([]*types.GroupPermission, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/permissions/groups?limit=%d", baseUrl, projectKey, limit)
	return getGroupPermissions(url, token)
}

func getRepositoryGroupPermissions(baseUrl string, projectKey string, repoSlug string, limit int, token string) ([]*types.GroupPermission, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/permissions/groups?limit=%d", baseUrl, projectKey, repoSlug, limit)
	return getGroupPermissions(url, token)
}

func getUserPermissions(url string, token string) ([]*types.UserPermission, error) {
	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var users types.UserPermissionsResponse
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, err
	}

	if !users.IsLastPage {
		pterm.Warning.Println("Not all User permissions fetched, try with a higher limit")
	}

	return users.Values, nil
}

func getProjectUserPermissions(baseUrl string, projectKey string, limit int, token string) ([]*types.UserPermission, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/permissions/users?limit=%d", baseUrl, projectKey, limit)
	return getUserPermissions(url, token)
}

func getRepositoryUserPermissions(baseUrl string, projectKey string, repoSlug string, limit int, token string) ([]*types.UserPermission, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/permissions/users?limit=%d", baseUrl, projectKey, repoSlug, limit)
	return getUserPermissions(url, token)
}

func getProjectPermissions(baseUrl string, projectKey string, limit int, token string, includeRepos bool) (*ProjectPermissions, error) {
	projectGroupPermissions, err := getProjectGroupPermissions(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	grantedPermissions := make(Permissions)
	for _, groupWithPermission := range projectGroupPermissions {
		if _, exists := grantedPermissions[groupWithPermission.Permission]; !exists {
			grantedPermissions[groupWithPermission.Permission] = new(Entities)
		}
		grantedPermissions[groupWithPermission.Permission].Groups = append(grantedPermissions[groupWithPermission.Permission].Groups, groupWithPermission.Group.Name)
	}

	projectUserPermissions, err := getProjectUserPermissions(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	for _, userWithPermission := range projectUserPermissions {
		if _, exists := grantedPermissions[userWithPermission.Permission]; !exists {
			grantedPermissions[userWithPermission.Permission] = new(Entities)
		}
		grantedPermissions[userWithPermission.Permission].Users = append(grantedPermissions[userWithPermission.Permission].Users, userWithPermission.User.Name)
	}

	defaultPermission, err := getDefaultProjectPermission(baseUrl, projectKey, token)
	if err != nil {
		return nil, err
	}

	isPublic, err := isProjectPublic(baseUrl, projectKey, token)
	if err != nil {
		return nil, err
	}

	projectPermissions := ProjectPermissions{
		Public:            isPublic,
		DefaultPermission: defaultPermission,
		Permissions:       &grantedPermissions,
	}

	if includeRepos {
		// Hent rettigheter for alle repositories i prosjektet
		projectRepositories, err := repository.GetProjectRepositories(baseUrl, projectKey, token, limit)
		if err != nil {
			return nil, err
		}
		projectPermissions.Repositories = make(map[string]*RepositoryPermissions)
		for repoSlug := range projectRepositories {
			repoPerms, err := getRepositoryPermissions(baseUrl, projectKey, repoSlug, limit, token)
			if err != nil {
				return nil, err
			}
			if len(*repoPerms.Permissions) > 0 {
				projectPermissions.Repositories[repoSlug] = repoPerms
			}
		}
	}

	return &projectPermissions, nil
}

func getRepositoryPermissions(baseUrl string, projectKey string, repoSlug string, limit int, token string) (*RepositoryPermissions, error) {
	repoGroupPermissions, err := getRepositoryGroupPermissions(baseUrl, projectKey, repoSlug, limit, token)
	if err != nil {
		return nil, err
	}

	grantedPermissions := make(Permissions)
	for _, groupWithPermission := range repoGroupPermissions {
		if _, exists := grantedPermissions[groupWithPermission.Permission]; !exists {
			grantedPermissions[groupWithPermission.Permission] = new(Entities)
		}
		grantedPermissions[groupWithPermission.Permission].Groups = append(grantedPermissions[groupWithPermission.Permission].Groups, groupWithPermission.Group.Name)
	}

	repoUserPermissions, err := getRepositoryUserPermissions(baseUrl, projectKey, repoSlug, limit, token)
	if err != nil {
		return nil, err
	}

	for _, userWithPermission := range repoUserPermissions {
		if _, exists := grantedPermissions[userWithPermission.Permission]; !exists {
			grantedPermissions[userWithPermission.Permission] = new(Entities)
		}
		grantedPermissions[userWithPermission.Permission].Users = append(grantedPermissions[userWithPermission.Permission].Users, userWithPermission.User.Name)
	}

	return &RepositoryPermissions{Permissions: &grantedPermissions}, nil
}

func prettyFormatProjectPermissions(projectPermissionsMap map[string]*ProjectPermissions) [][]string {
	var data [][]string
	data = append(data, []string{"Project", "Repository", "Permission", "Groups", "Users"})

	projects := common.GetLexicallySortedKeys(projectPermissionsMap)
	for _, projectKey := range projects {
		formattedProjectPermissions := prettyFormatPermissions(projectKey, "ALL", projectPermissionsMap[projectKey].Permissions)
		data = append(data, formattedProjectPermissions...)

		formattedRepositoryPermissions := prettyFormatRepositoryPermissions(projectKey, projectPermissionsMap[projectKey].Repositories)
		data = append(data, formattedRepositoryPermissions...)
	}
	return data
}

func prettyFormatRepositoryPermissions(projectKey string, repositoryPermissionsMap map[string]*RepositoryPermissions) [][]string {
	var data [][]string

	repositories := common.GetLexicallySortedKeys(repositoryPermissionsMap)
	for _, repoSlug := range repositories {
		repoPermissions := prettyFormatPermissions(projectKey, repoSlug, repositoryPermissionsMap[repoSlug].Permissions)
		data = append(data, repoPermissions...)
	}
	return data
}

func prettyFormatPermissions(projectKey string, repoSlug string, permissions *Permissions) [][]string {
	var data [][]string

	if permissions == nil {
		return data
	}

	for permission, v := range *permissions {
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

		data = append(data, []string{projectKey, repoSlug, permission, groups, users})
	}
	return data
}
