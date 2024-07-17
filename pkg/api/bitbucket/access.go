package bitbucket

import (
	"encoding/json"
	"fmt"
	"git.spk.no/infra/bucketctl/pkg/api/bitbucket/types"
	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"git.spk.no/infra/bucketctl/pkg/common"
	"git.spk.no/infra/bucketctl/pkg/logger"
)

func isProjectPublic(baseUrl string, projectKey string, token string) (bool, error) {
	p, err := GetProject(baseUrl, projectKey, token)
	if err != nil {
		return false, err
	}
	return p.Public, err
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
		logger.Warn("not all Group permissions fetched, try with a higher limit")
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
		logger.Warn("not all User permissions fetched, try with a higher limit")
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

func GetProjectAccess(baseUrl string, projectKey string, limit int, token string) (*ProjectConfigSpec, error) {
	defaultPermission, err := getDefaultProjectPermission(baseUrl, projectKey, token)
	if err != nil {
		return nil, err
	}

	isPublic, err := isProjectPublic(baseUrl, projectKey, token)
	if err != nil {
		return nil, err
	}

	projectGroupPermissions, err := getProjectGroupPermissions(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	grantedPermissionsMap := make(map[string]*Entities)
	for _, groupWithPermission := range projectGroupPermissions {
		if _, exists := grantedPermissionsMap[groupWithPermission.Permission]; !exists {
			grantedPermissionsMap[groupWithPermission.Permission] = new(Entities)
		}
		grantedPermissionsMap[groupWithPermission.Permission].Groups = append(grantedPermissionsMap[groupWithPermission.Permission].Groups, groupWithPermission.Group.Name)
	}

	projectUserPermissions, err := getProjectUserPermissions(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	for _, userWithPermission := range projectUserPermissions {
		if _, exists := grantedPermissionsMap[userWithPermission.Permission]; !exists {
			grantedPermissionsMap[userWithPermission.Permission] = new(Entities)
		}
		grantedPermissionsMap[userWithPermission.Permission].Users = append(grantedPermissionsMap[userWithPermission.Permission].Users, userWithPermission.User.Name)
	}

	grantedPermissions := marshalPermissionsMapToList(grantedPermissionsMap)

	projectAccess := &ProjectConfigSpec{
		ProjectKey:        projectKey,
		Public:            &isPublic,
		DefaultPermission: &defaultPermission,
		Permissions:       grantedPermissions,
	}

	return projectAccess, nil
}

func GetProjectRepositoriesPermissions(baseUrl string, projectKey string, limit int, token string) (*RepositoriesProperties, error) {
	// Hent rettigheter for alle repositories i prosjektet
	projectRepositories, err := GetLexicallySortedProjectRepositoriesNames(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	repositoriesProperties := new(RepositoriesProperties)
	for _, repoSlug := range projectRepositories {
		repoPermissions, err := GetRepositoryPermissions(baseUrl, projectKey, repoSlug, limit, token)
		if err != nil {
			return nil, err
		}
		*repositoriesProperties = append(*repositoriesProperties, &RepositoryProperties{RepoSlug: repoSlug, Permissions: repoPermissions})
	}
	return repositoriesProperties, nil
}

func GetRepositoryPermissions(baseUrl string, projectKey string, repoSlug string, limit int, token string) (*Permissions, error) {
	repoGroupPermissions, err := getRepositoryGroupPermissions(baseUrl, projectKey, repoSlug, limit, token)
	if err != nil {
		return nil, err
	}

	grantedPermissionsMap := make(map[string]*Entities)
	for _, groupWithPermission := range repoGroupPermissions {
		if _, exists := grantedPermissionsMap[groupWithPermission.Permission]; !exists {
			grantedPermissionsMap[groupWithPermission.Permission] = new(Entities)
		}
		grantedPermissionsMap[groupWithPermission.Permission].Groups = append(grantedPermissionsMap[groupWithPermission.Permission].Groups, groupWithPermission.Group.Name)
	}

	repoUserPermissions, err := getRepositoryUserPermissions(baseUrl, projectKey, repoSlug, limit, token)
	if err != nil {
		return nil, err
	}

	for _, userWithPermission := range repoUserPermissions {
		if _, exists := grantedPermissionsMap[userWithPermission.Permission]; !exists {
			grantedPermissionsMap[userWithPermission.Permission] = new(Entities)
		}
		grantedPermissionsMap[userWithPermission.Permission].Users = append(grantedPermissionsMap[userWithPermission.Permission].Users, userWithPermission.User.Name)
	}

	grantedPermissions := marshalPermissionsMapToList(grantedPermissionsMap)

	return grantedPermissions, nil
}

func marshalPermissionsMapToList(permissionsMap map[string]*Entities) (permissionsList *Permissions) {
	permissionsList = new(Permissions)
	for _, permission := range common.GetLexicallySortedKeys(permissionsMap) {
		*permissionsList = append(*permissionsList, &Permission{
			Name: permission,
			Entities: &Entities{
				Groups: permissionsMap[permission].Groups,
				Users:  permissionsMap[permission].Users,
			},
		})
	}
	return permissionsList
}
