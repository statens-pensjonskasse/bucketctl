package access

import (
	"bucketctl/pkg/api/bitbucket"
	"bucketctl/pkg/api/bitbucket/types"
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/common"
	"fmt"
	"github.com/pterm/pterm"
	"strconv"
)

func setProjectAccessProperties(baseUrl string, projectKey string, token string, toUpdate *ProjectConfigSpec) error {
	if toUpdate.Public != nil {
		if err := setProjectPublicProperty(baseUrl, projectKey, *toUpdate.Public, token); err != nil {
			return err
		}
	}
	if toUpdate.DefaultPermission != nil {
		if err := setProjectDefaultPermission(baseUrl, projectKey, *toUpdate.DefaultPermission, token); err != nil {
			return err
		}
	}
	return nil
}

func setProjectDefaultPermission(baseUrl string, projectKey string, permission string, token string) error {
	if err := changeDefaultProjectPermission(baseUrl, projectKey, permission, true, token); err != nil {
		return err
	}
	// Revoke all other default permissions
	for _, p := range []string{"REPO_CREATE", "PROJECT_ADMIN", "PROJECT_WRITE", "PROJECT_READ"} {
		if p == permission {
			continue
		}
		if err := revokeDefaultProjectPermission(baseUrl, projectKey, p, token); err != nil {
			return err
		}
	}
	pterm.Printfln("%s default permission for project '%s' to '%s'", pterm.Blue("🖋️Changed"), projectKey, permission)
	return nil
}

func setProjectPublicProperty(baseUrl string, projectKey string, isPublic bool, token string) error {
	if _, err := bitbucket.UpdateProject(baseUrl, token, &types.Project{Key: projectKey, Public: isPublic}); err != nil {
		return err
	}
	var action string
	if isPublic {
		action = pterm.Green("🔓 Opened")
	} else {
		action = pterm.Red("🔒 Closed")
	}
	pterm.Printfln("%s public access for project '%s'", action, projectKey)
	return nil
}

func changeDefaultProjectPermission(baseUrl string, projectKey string, permission string, allow bool, token string) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/permissions/%s/all", baseUrl, projectKey, permission)
	params := map[string]string{
		"allow": strconv.FormatBool(allow),
		// Workaround for https://confluence.atlassian.com/cloudkb/xsrf-check-failed-when-calling-cloud-apis-826874382.html
		"Header X-Atlassian-Token": "no-check",
	}
	_, err := common.PostRequest(url, token, nil, params)
	if err != nil {
		return err
	}
	return nil
}

func revokeDefaultProjectPermission(baseUrl string, projectKey string, permission string, token string) error {
	return changeDefaultProjectPermission(baseUrl, projectKey, permission, false, token)
}
