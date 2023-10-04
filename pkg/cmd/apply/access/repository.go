package access

import (
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/common"
	"fmt"
	"github.com/pterm/pterm"
)

func findRepositoriesPermissionsChanges(desired *RepositoriesProperties, actual *RepositoriesProperties) (toCreate *RepositoriesProperties, toDelete *RepositoriesProperties) {
	toCreate = new(RepositoriesProperties)
	toDelete = new(RepositoriesProperties)
	for repoSlug, repo := range GroupRepositories(desired, actual) {
		permissionsToCreate, permissionsToDelete := FindPermissionsToChange(repo.Desired.Permissions, repo.Actual.Permissions)
		if len(*permissionsToCreate) > 0 {
			*toCreate = append(*toCreate, &RepositoryProperties{RepoSlug: repoSlug, Permissions: permissionsToCreate})
		}
		if len(*permissionsToDelete) > 0 {
			*toDelete = append(*toDelete, &RepositoryProperties{RepoSlug: repoSlug, Permissions: permissionsToDelete})
		}
	}
	return toCreate, toDelete
}

func setRepositoriesPermissions(baseUrl string, projectKey string, token string, toCreate *RepositoriesProperties, toDelete *RepositoriesProperties) error {
	for _, r := range *toDelete {
		if err := removeRepositoryPermissions(baseUrl, projectKey, r.RepoSlug, token, r.Permissions); err != nil {
			return err
		}
	}
	for _, r := range *toCreate {
		if err := grantRepositoryPermissions(baseUrl, projectKey, r.RepoSlug, token, r.Permissions); err != nil {
			return err
		}
	}
	return nil
}

func removeRepositoryPermissions(baseUrl string, projectKey string, repoSlug string, token string, permissions *Permissions) error {
	if permissions != nil && len(*permissions) > 0 {
		for _, permission := range *permissions {
			for _, user := range permission.Entities.Users {
				if err := removeUserRepositoryPermissions(baseUrl, projectKey, repoSlug, token, user); err != nil {
					return err
				}
				pterm.Printfln("%s permissions for user '%s' in repository '%s/%s'", pterm.Red("🙅 Revoked"), user, projectKey, repoSlug)
			}
			for _, group := range permission.Entities.Groups {
				if err := removeGroupRepositoryPermissions(baseUrl, projectKey, repoSlug, token, group); err != nil {
					return err
				}
				pterm.Printfln("%s permissions for group '%s' in repository '%s/%s'", pterm.Red("🚫 Revoked"), group, projectKey, repoSlug)
			}
		}
	}
	return nil
}

func removeUserRepositoryPermissions(baseUrl string, projectKey string, repoSlug string, token string, user string) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/permissions/users", baseUrl, projectKey, repoSlug)
	params := map[string]string{
		"name": user,
	}

	if _, err := common.DeleteRequest(url, token, params); err != nil {
		return err
	}
	return nil
}

func removeGroupRepositoryPermissions(baseUrl string, projectKey string, repoSlug string, token string, group string) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/permissions/groups", baseUrl, projectKey, repoSlug)
	params := map[string]string{
		"name": group,
	}

	if _, err := common.DeleteRequest(url, token, params); err != nil {
		return err
	}
	return nil
}

func grantRepositoryPermissions(baseUrl string, projectKey string, repoSlug string, token string, permissions *Permissions) error {
	if permissions != nil && len(*permissions) > 0 {
		for _, permission := range *permissions {
			for _, user := range permission.Entities.Users {
				if err := grantUserRepositoryPermission(baseUrl, projectKey, repoSlug, token, user, permission.Name); err != nil {
					return err
				}
				pterm.Printfln("%s user '%s' permission '%s' in repository '%s/%s'", pterm.Green("🧑 Granted"), user, permission.Name, projectKey, repoSlug)
			}
			for _, group := range permission.Entities.Groups {
				if err := grantGroupRepositoryPermission(baseUrl, projectKey, repoSlug, token, group, permission.Name); err != nil {
					return err
				}
				pterm.Printfln("%s group '%s' permission '%s' in repository '%s/%s'", pterm.Green("👥 Granted"), group, permission.Name, projectKey, repoSlug)
			}
		}
	}
	return nil
}

func grantUserRepositoryPermission(baseUrl string, projectKey string, repoSlug string, token string, user string, permission string) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/permissions/users", baseUrl, projectKey, repoSlug)
	params := map[string]string{
		"name":       user,
		"permission": permission,
	}

	if _, err := common.PutRequest(url, token, nil, params); err != nil {
		return err
	}
	return nil
}

func grantGroupRepositoryPermission(baseUrl string, projectKey string, reposSlug string, token string, group string, permission string) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/permissions/groups", baseUrl, projectKey, reposSlug)
	params := map[string]string{
		"name":       group,
		"permission": permission,
	}

	if _, err := common.PutRequest(url, token, nil, params); err != nil {
		return err
	}
	return nil
}
