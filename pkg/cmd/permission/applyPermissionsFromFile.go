package permission

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
	"os"
)

var (
	fileName string
)

var applyPermissionsFromFile = &cobra.Command{
	Use:  "apply",
	RunE: applyPermissions,
}

func init() {
	applyPermissionsFromFile.Flags().StringVarP(&fileName, "file", "f", "", "Permission file")
	applyPermissionsFromFile.MarkFlagRequired("file")
	viper.BindPFlag("file", applyPermissionsFromFile.Flags().Lookup("file"))
}

func applyPermissions(cmd *cobra.Command, args []string) error {
	file := viper.GetString("file")
	baseUrl := viper.GetString("baseUrl")
	limit := viper.GetInt("limit")
	token := viper.GetString("token")
	includeRepos := viper.GetBool("include-repos")

	// Les inn fil (yaml eller json) med ønskede tilganger
	var desiredPermissions *GrantedProjectPermissions
	if err := pkg.ReadConfigFile(file, &desiredPermissions); err != nil {
		return err
	}

	// Finn aktuelle tilganger for prosjekter definert i ønskede tilganger
	actualPermissions := &GrantedProjectPermissions{
		Project: map[string]*ProjectPermissionSet{},
	}
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(len(desiredPermissions.Project)).WithRemoveWhenDone(true).WithWriter(os.Stderr).Start()
	for projectKey := range desiredPermissions.Project {
		progressBar.Title = projectKey

		projectPermission, err := GetProjectPermissions(baseUrl, projectKey, limit, token, includeRepos)
		if err != nil {
			return err
		}

		actualPermissions.Project[projectKey] = new(ProjectPermissionSet)
		actualPermissions.Project[projectKey].Permissions = projectPermission.Permissions

		progressBar.Increment()
	}

	permissionsToBeRemoved, permissionsToBeGranted := findProjectPermissionDifference(desiredPermissions, actualPermissions)

	// Fjern alle tilganger som ikke er ønsket
	progressBar, _ = pterm.DefaultProgressbar.WithTotal(len(permissionsToBeRemoved.Project)).WithRemoveWhenDone(true).WithWriter(os.Stderr).WithTitle("Fjerner tilganger").Start()
	for projectKey := range permissionsToBeRemoved.Project {
		if err := removeProjectPermissions(baseUrl, projectKey, token, permissionsToBeRemoved.Project[projectKey]); err != nil {
			return err
		}
		progressBar.Increment()
	}

	// Gi ønskede tilganger
	progressBar, _ = pterm.DefaultProgressbar.WithTotal(len(permissionsToBeGranted.Project)).WithRemoveWhenDone(true).WithWriter(os.Stderr).WithTitle("Gir tilganger").Start()
	for projectKey := range permissionsToBeGranted.Project {
		if err := grantProjectPermissions(baseUrl, projectKey, token, permissionsToBeGranted.Project[projectKey]); err != nil {
			return err
		}
		progressBar.Increment()
	}

	return nil
}

func findProjectPermissionDifference(desiredState *GrantedProjectPermissions, actualState *GrantedProjectPermissions) (permissionsToBeRemoved *GrantedProjectPermissions, permissionsToBeGranted *GrantedProjectPermissions) {
	permissionsToBeRemoved = &GrantedProjectPermissions{
		Project: map[string]*ProjectPermissionSet{},
	}
	permissionsToBeGranted = &GrantedProjectPermissions{
		Project: map[string]*ProjectPermissionSet{},
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
func (setA *ProjectPermissionSet) getPermissionSetDifference(setB *ProjectPermissionSet) *ProjectPermissionSet {
	difference := new(ProjectPermissionSet)
	difference.Permissions = make(map[string]*Entities)

	for permission := range setA.Permissions {
		difference.Permissions[permission] = new(Entities)
		entriesInA := setA.Permissions[permission]
		entriesInB, existsInB := setB.Permissions[permission]

		if !existsInB {
			// Dersom tilgangen ikke finnes i B så legger vi til alle entries i A i det relative komplementet
			difference.Permissions[permission].Users = setA.Permissions[permission].Users
			difference.Permissions[permission].Groups = setA.Permissions[permission].Groups
		} else {
			// Hvis tilgangen finnes i B så må sjekke alle elementene hver for seg
			for _, user := range entriesInA.Users {
				if !entriesInB.containsUser(user) {
					difference.Permissions[permission].Users = append(difference.Permissions[permission].Users, user)
				}
			}
			for _, group := range entriesInA.Groups {
				if !entriesInB.containsGroup(group) {
					difference.Permissions[permission].Groups = append(difference.Permissions[permission].Groups, group)
				}
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

func removeProjectPermissions(baseUrl string, projectKey string, token string, permissionSet *ProjectPermissionSet) error {
	for _, entity := range permissionSet.Permissions {
		for _, user := range entity.Users {
			if err := removeUserProjectPermissions(baseUrl, projectKey, token, user); err != nil {
				return err
			}
		}
		for _, group := range entity.Groups {
			if err := removeGroupProjectPermissions(baseUrl, projectKey, token, group); err != nil {
				return err
			}
		}
	}
	return nil
}

func removeUserProjectPermissions(baseUrl string, projectKey string, token string, user string) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/permissions/users", baseUrl, projectKey)
	params := map[string]string{
		"name": user,
	}

	if _, err := pkg.DeleteRequest(url, token, params); err != nil {
		return err
	}
	return nil
}

func removeGroupProjectPermissions(baseUrl string, projectKey string, token string, group string) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/permissions/groups", baseUrl, projectKey)
	params := map[string]string{
		"name": group,
	}

	if _, err := pkg.DeleteRequest(url, token, params); err != nil {
		return err
	}
	return nil
}

func grantProjectPermissions(baseUrl string, projectKey string, token string, permissionSet *ProjectPermissionSet) error {
	for permission, entity := range permissionSet.Permissions {
		for _, user := range entity.Users {
			if err := grantUserProjectPermission(baseUrl, projectKey, token, user, permission); err != nil {
				return err
			}
		}
		for _, group := range entity.Groups {
			if err := grantGroupProjectPermission(baseUrl, projectKey, token, group, permission); err != nil {
				return err
			}
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

	if _, err := pkg.PutRequest(url, token, params); err != nil {
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

	if _, err := pkg.PutRequest(url, token, params); err != nil {
		return err
	}
	return nil
}
