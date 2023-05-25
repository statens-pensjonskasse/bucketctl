package permission

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gobit/pkg"
	"gobit/pkg/cmd/repository"
	"gobit/pkg/types"
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

type ProjectPermissionSet struct {
	PermissionSet
	Repositories map[string]*PermissionSet `json:"repositories,omitempty" yaml:"repositories,omitempty"`
}

type GrantedProjectPermissions struct {
	Project map[string]*ProjectPermissionSet `json:"" yaml:",inline"`
}

type GrantedRepositoryPermissions struct {
	Repository map[string]*PermissionSet `json:"" yaml:",inline"`
}

var Cmd = &cobra.Command{
	Use:     "permissions",
	Short:   "Bitbucket project permission commands",
	Aliases: []string{"perm"},
}

func init() {
	Cmd.AddCommand(listPermissionsCmd)
	Cmd.AddCommand(listAllPermissionsCmd)
	Cmd.AddCommand(applyPermissionsFromFile)
}

func getGroupPermissions(url string, token string) ([]types.GroupPermission, error) {
	body, err := pkg.GetRequestBody(url, token)
	if err != nil {
		return []types.GroupPermission{}, err
	}

	var groups types.GroupPermissionsResponse
	if err := json.Unmarshal(body, &groups); err != nil {
		return []types.GroupPermission{}, err
	}

	if !groups.IsLastPage {
		pterm.Warning.Println("Not all Group permissions fetched, try with a higher limit")
	}

	return groups.Values, nil
}

func getProjectGroupPermissions(baseUrl string, projectKey string, limit int, token string) ([]types.GroupPermission, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/permissions/groups?limit=%d", baseUrl, projectKey, limit)
	return getGroupPermissions(url, token)
}

func getRepositoryGroupPermissions(baseUrl string, projectKey string, repoSlug string, limit int, token string) ([]types.GroupPermission, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/permissions/groups?limit=%d", baseUrl, projectKey, repoSlug, limit)
	return getGroupPermissions(url, token)
}

func getUserPermissions(url string, token string) ([]types.UserPermission, error) {
	body, err := pkg.GetRequestBody(url, token)
	if err != nil {
		return []types.UserPermission{}, err
	}

	var users types.UserPermissionsResponse
	if err := json.Unmarshal(body, &users); err != nil {
		return []types.UserPermission{}, err
	}

	if !users.IsLastPage {
		pterm.Warning.Println("Not all User permissions fetched, try with a higher limit")
	}

	return users.Values, nil
}

func getProjectUserPermissions(baseUrl string, projectKey string, limit int, token string) ([]types.UserPermission, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/permissions/users?limit=%d", baseUrl, projectKey, limit)
	return getUserPermissions(url, token)
}

func getRepositoryUserPermissions(baseUrl string, projectKey string, repoSlug string, limit int, token string) ([]types.UserPermission, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/permissions/users?limit=%d", baseUrl, projectKey, repoSlug, limit)
	return getUserPermissions(url, token)
}

func GetProjectPermissions(baseUrl string, projectKey string, limit int, token string) (*GrantedProjectPermissions, error) {
	projectPermissions := &GrantedProjectPermissions{
		Project: map[string]*ProjectPermissionSet{},
	}
	projectPermissions.Project[projectKey] = new(ProjectPermissionSet)

	projectGroupPermissions, err := getProjectGroupPermissions(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	grantedPermissions := make(map[string]*Entities)
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

	projectPermissions.Project[projectKey].Permissions = grantedPermissions

	// Hent repo-rettigheter
	projectRepositories, err := repository.GetProjectRepositories(baseUrl, projectKey, limit)
	if err != nil {
		return nil, err
	}
	projectPermissions.Project[projectKey].Repositories = make(map[string]*PermissionSet)
	for _, r := range projectRepositories {
		repoPerms, _ := getRepositoryPermissions(baseUrl, projectKey, r.Slug, limit, token)
		if len(repoPerms.Permissions) > 0 {
			projectPermissions.Project[projectKey].Repositories[r.Slug] = repoPerms
		}
	}

	return projectPermissions, nil
}

func getRepositoryPermissions(baseUrl string, projectKey string, repoSlug string, limit int, token string) (*PermissionSet, error) {
	repoGroupPermissions, err := getRepositoryGroupPermissions(baseUrl, projectKey, repoSlug, limit, token)
	if err != nil {
		return nil, err
	}

	grantedPermissions := make(map[string]*Entities)
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

	return &PermissionSet{Permissions: grantedPermissions}, nil
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

func PrettyFormatRepositoryPermissions(projectPermissions *GrantedRepositoryPermissions) [][]string {
	// Sorter repoene alfabetisk
	projects := make([]string, 0, len(projectPermissions.Repository))
	for k := range projectPermissions.Repository {
		projects = append(projects, k)
	}
	sort.Strings(projects)

	var data [][]string
	data = append(data, []string{"Project", "Permission", "Groups", "Users"})

	for _, k := range projects {
		proj := k
		for permission, v := range projectPermissions.Repository[k].Permissions {
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
