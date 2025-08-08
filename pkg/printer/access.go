package printer

import (
	"strings"

	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"git.spk.no/infra/bucketctl/pkg/common"
)

func PrettyFormatAccess(projectConfig *ProjectConfig) [][]string {
	projectKey := projectConfig.Spec.ProjectKey
	projectConfigSpec := projectConfig.Spec
	projectRepositories := projectConfig.Spec.Repositories

	var data [][]string
	data = append(data, []string{"Project", "Repository", "Permission", "Groups", "Users"})

	if projectConfigSpec.Public != nil && *projectConfigSpec.Public {
		data = append(data, []string{projectKey, ALL, "PROJECT_READ", ALL, ALL})
	}
	if projectConfigSpec.DefaultPermission != nil {
		data = append(data, []string{projectKey, ALL, *projectConfigSpec.DefaultPermission, AUTH, AUTH})
	}

	formattedProjectPermissions := prettyFormatPermissions(projectKey, ALL, projectConfigSpec.Permissions)
	data = append(data, formattedProjectPermissions...)

	formattedRepositoryPermissions := prettyFormatRepositoriesPermissions(projectKey, projectRepositories)
	data = append(data, formattedRepositoryPermissions...)
	return data
}

func prettyFormatRepositoriesPermissions(projectKey string, repositoriesProperties *RepositoriesProperties) [][]string {
	var data [][]string

	if repositoriesProperties == nil {
		return data
	}

	repoPermissionsMap := make(map[string]*Permissions, len(*repositoriesProperties))
	for _, r := range *repositoriesProperties {
		repoPermissionsMap[r.RepoSlug] = r.Permissions
	}

	for _, slug := range common.GetLexicallySortedKeys(repoPermissionsMap) {
		repoPermissions := prettyFormatPermissions(projectKey, slug, repoPermissionsMap[slug])
		data = append(data, repoPermissions...)
	}
	return data
}

func prettyFormatPermissions(projectKey string, repoSlug string, permissions *Permissions) [][]string {
	var data [][]string

	if permissions == nil {
		return data
	}

	permissionsEntitiesMap := make(map[string]*Entities, len(*permissions))
	for _, p := range *permissions {
		permissionsEntitiesMap[p.Name] = p.Entities
	}

	for _, permission := range common.GetLexicallySortedKeys(permissionsEntitiesMap) {
		var users string
		for _, user := range permissionsEntitiesMap[permission].Users {
			users += user + "\n"
		}
		users = strings.Trim(users, "\n")
		var groups string
		for _, group := range permissionsEntitiesMap[permission].Groups {
			groups += group + "\n"
		}
		groups = strings.Trim(groups, "\n")

		data = append(data, []string{projectKey, repoSlug, permission, groups, users})
	}
	return data
}
