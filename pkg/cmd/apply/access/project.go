package access

import (
	"fmt"

	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"git.spk.no/infra/bucketctl/pkg/common"
	"git.spk.no/infra/bucketctl/pkg/logger"
	"github.com/pterm/pterm"
)

func setProjectPermissions(baseUrl string, projectKey string, token string, toCreate *Permissions, toDelete *Permissions) error {
	if err := removeProjectPermissions(baseUrl, projectKey, token, toDelete); err != nil {
		return err
	}
	if err := grantProjectPermissions(baseUrl, projectKey, token, toCreate); err != nil {
		return err
	}
	return nil
}

func removeProjectPermissions(baseUrl string, projectKey string, token string, permissions *Permissions) error {
	for _, permission := range *permissions {
		for _, user := range permission.Entities.Users {
			if err := removeUserProjectPermissions(baseUrl, projectKey, token, user); err != nil {
				return err
			}
			logger.Log("%s permissions for user '%s' in project '%s'", pterm.Red("🙅 Revoked"), user, projectKey)
		}
		for _, group := range permission.Entities.Groups {
			if err := removeGroupProjectPermissions(baseUrl, projectKey, token, group); err != nil {
				return err
			}
			logger.Log("%s permissions for group '%s' in project '%s'", pterm.Red("🚫 Revoked"), group, projectKey)
		}
	}
	return nil
}

func removeUserProjectPermissions(baseUrl string, projectKey string, token string, user string) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/permissions/users", baseUrl, projectKey)
	params := map[string]string{
		"name": user,
	}

	if _, err := common.DeleteRequest(url, token, params); err != nil {
		return err
	}
	return nil
}

func removeGroupProjectPermissions(baseUrl string, projectKey string, token string, group string) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/permissions/groups", baseUrl, projectKey)
	params := map[string]string{
		"name": group,
	}

	if _, err := common.DeleteRequest(url, token, params); err != nil {
		return err
	}
	return nil
}

func grantProjectPermissions(baseUrl string, projectKey string, token string, permissions *Permissions) error {
	for _, permission := range *permissions {
		for _, user := range permission.Entities.Users {
			if err := grantUserProjectPermission(baseUrl, projectKey, token, user, permission.Name); err != nil {
				return err
			}
			logger.Log("%s user '%s' permission '%s' for project %s", pterm.Green("🧑 Granted"), user, permission.Name, projectKey)
		}
		for _, group := range permission.Entities.Groups {
			if err := grantGroupProjectPermission(baseUrl, projectKey, token, group, permission.Name); err != nil {
				return err
			}
			logger.Log("%s group '%s' permission '%s' for project %s", pterm.Green("👥 Granted"), group, permission.Name, projectKey)
		}
	}
	return nil
}

func grantUserProjectPermission(baseUrl string, projectKey string, token string, user string, permission string) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/permissions/users", baseUrl, projectKey)
	params := map[string]string{
		"name":       user,
		"permission": permission,
	}

	if _, err := common.PutRequest(url, token, nil, params); err != nil {
		return err
	}
	return nil
}

func grantGroupProjectPermission(baseUrl string, projectKey string, token string, group string, permission string) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/permissions/groups", baseUrl, projectKey)
	params := map[string]string{
		"name":       group,
		"permission": permission,
	}

	if _, err := common.PutRequest(url, token, nil, params); err != nil {
		return err
	}
	return nil
}
