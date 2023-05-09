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

	var desiredState *GrantedProjectPermissions
	err := pkg.ReadConfigFile("permissions.yaml", &desiredState)
	if err != nil {
		return err
	}

	actualState := &GrantedProjectPermissions{
		Project: map[string]*PermissionSet{},
	}
	for projectKey := range desiredState.Project {
		projectPermission, err := GetProjectPermissions(baseUrl, projectKey, limit, token)
		if err != nil {
			return err
		}
		actualState.Project[projectKey] = new(PermissionSet)
		actualState.Project[projectKey].Permissions = projectPermission.Project[projectKey].Permissions
	}

	permissionsToBeRemoved, permissionsToBeGranted := findProjectPermissionDifference(desiredState, actualState)

	pkg.PrintData(permissionsToBeRemoved, PrettyFormatProjectPermissions)
	pkg.PrintData(permissionsToBeGranted, PrettyFormatProjectPermissions)

	return nil
}

func findProjectPermissionDifference(desiredState *GrantedProjectPermissions, actualState *GrantedProjectPermissions) (permissionsToBeRemoved *GrantedProjectPermissions, permissionsToBeGranted *GrantedProjectPermissions) {
	permissionsToBeRemoved = &GrantedProjectPermissions{
		Project: map[string]*PermissionSet{},
	}
	permissionsToBeGranted = &GrantedProjectPermissions{
		Project: map[string]*PermissionSet{},
	}
	for projectKey := range desiredState.Project {
		// Finner tilganger i 'actualState' som ikke finnes i 'desiredState'. Disse tilgangene skal fjernes.
		permissionsToBeRemoved.Project[projectKey] = actualState.Project[projectKey].getPermissionSetDifference(desiredState.Project[projectKey])
		// Finner tilganger i 'desiredState' som ikke finnes i 'actualState'. Disse tilgangene skal gis.
		permissionsToBeGranted.Project[projectKey] = desiredState.Project[projectKey].getPermissionSetDifference(actualState.Project[projectKey])
	}

	return permissionsToBeRemoved, permissionsToBeGranted
}

// Finner det relative komplementet til setA i setB
func (setA *PermissionSet) getPermissionSetDifference(setB *PermissionSet) *PermissionSet {
	difference := new(PermissionSet)
	difference.Permissions = make(map[string]*Entities)

	for _, permission := range PermissionTypes {
		difference.Permissions[permission] = new(Entities)
		permissionInB := setB.Permissions[permission]
		permissionInA := setA.Permissions[permission]
		for _, user := range permissionInA.Users {
			if !permissionInB.containsUser(user) {
				difference.Permissions[permission].Users = append(difference.Permissions[permission].Users, user)
			}
		}
		for _, group := range permissionInA.Groups {
			if !permissionInB.containsGroup(group) {
				difference.Permissions[permission].Groups = append(difference.Permissions[permission].Groups, group)
			}
		}
	}
	return difference
}

func (entities Entities) containsUser(user string) bool {
	for _, u := range entities.Users {
		if user == u {
			return true
		}
	}
	return false
}

func (entities Entities) containsGroup(group string) bool {
	for _, g := range entities.Groups {
		if group == g {
			return true
		}
	}
	return false
}
