package project

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
)

var applyPermissionsFromFile = &cobra.Command{
	Use:  "apply",
	RunE: applyPermissions,
}

func removeUserPermission(baseUrl string, projectKey string, token string, permission UserPermission) error {
	url := fmt.Sprintf("%s/rest/api/1.0/projects/%s", baseUrl, projectKey)

	params := map[string]string{
		"user": permission.User.Name,
	}

	_, err := pkg.DeleteRequestWithParams(url, token, params)
	if err != nil {
		return err
	}

	return nil
}

func applyPermissions(cmd *cobra.Command, args []string) error {
	var baseUrl = viper.GetString("baseUrl")
	var limit = viper.GetInt("limit")
	var token = viper.GetString("token")

	var desiredState *ProjectPermissions
	err := pkg.ReadConfigFile("permissions.yaml", &desiredState)
	if err != nil {
		return err
	}

	actualState := &ProjectPermissions{
		Project: map[string]*PermissionObjects{},
	}
	for projectKey := range desiredState.Project {
		projectPermission, err := GetProjectPermissions(baseUrl, projectKey, limit, token)
		if err != nil {
			return err
		}
		actualState.Project[projectKey] = new(PermissionObjects)
		actualState.Project[projectKey].Permissions = projectPermission.Project[projectKey].Permissions
	}

	//toBeRemoved := &ProjectPermissions{
	//	Project: map[string]*PermissionObjects{},
	//}
	//for projectKey := range desiredState.Project {
	//	toBeRemoved.Project[projectKey] = new(PermissionObjects)
	//	for _, permission := range PermissionTypes {

	//		desiredUsers := desiredState.Project[projectKey].Permissions[permission].Users
	//		actualUsers := actualState.Project[projectKey].Permissions[permission].Users
	//		for user := range actualUsers {
	//			if user
	//		}
	//	}
	//}

	pkg.PrintData(desiredState, PrettyFormatProjectPermissions)
	pkg.PrintData(actualState, PrettyFormatProjectPermissions)

	return nil
}
