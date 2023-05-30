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
	var desiredPermissions map[string]ProjectPermissions
	if err := pkg.ReadConfigFile(file, &desiredPermissions); err != nil {
		return err
	}

	progressBar, _ := pterm.DefaultProgressbar.WithTotal(len(desiredPermissions)).WithRemoveWhenDone(true).WithWriter(os.Stderr).Start()
	for projectKey, desiredState := range desiredPermissions {
		progressBar.Title = projectKey
		// Finn gjeldende tilganger
		actualProjectPermissions, err := GetProjectPermissions(baseUrl, projectKey, limit, token, includeRepos)
		if err != nil {
			return err
		}

		// Finn forskjeller i gjeldende og ønskede tilganger
		permissionsToBeRemoved, permissionsToBeGranted := findProjectPermissionDifference(desiredState, actualProjectPermissions)

		// Fjern alle prosjekt-tilganger som ikke er ønsket
		if err := removeProjectPermissions(baseUrl, projectKey, token, permissionsToBeRemoved); err != nil {
			return err
		}

		// Gi ønskede prosjekt-tilganger
		if err := grantProjectPermissions(baseUrl, projectKey, token, permissionsToBeGranted); err != nil {
			return err
		}

		progressBar.Increment()
	}

	return nil
}

func findProjectPermissionDifference(desiredState ProjectPermissions, actualState ProjectPermissions) (permissionsToBeRemoved ProjectPermissions, permissionsToBeGranted ProjectPermissions) {
	permissionsToBeRemoved = ProjectPermissions{}
	permissionsToBeGranted = ProjectPermissions{}
	// Finner tilganger i 'actualState' som ikke finnes i 'desiredState'. Disse tilgangene skal fjernes.
	permissionsToBeRemoved.Permissions = actualState.Permissions.getPermissionSetDifference(desiredState.Permissions)
	// Finner tilganger i 'desiredState' som ikke finnes i 'actualState'. Disse tilgangene skal gis.
	permissionsToBeGranted.Permissions = desiredState.Permissions.getPermissionSetDifference(actualState.Permissions)

	return permissionsToBeRemoved, permissionsToBeGranted
}

// Finner det relative komplementet til setA i setB
func (setA PermissionSet) getPermissionSetDifference(setB PermissionSet) PermissionSet {
	difference := make(PermissionSet)

	for permission := range setA {
		difference[permission] = new(Entities)
		entriesInA := setA[permission]
		entriesInB, existsInB := setB[permission]

		if !existsInB {
			// Dersom tilgangen ikke finnes i B så legger vi til alle entries i A i det relative komplementet
			difference[permission].Users = setA[permission].Users
			difference[permission].Groups = setA[permission].Groups
		} else {
			// Hvis tilgangen finnes i B så må sjekke alle elementene hver for seg
			for _, user := range entriesInA.Users {
				if !entriesInB.containsUser(user) {
					difference[permission].Users = append(difference[permission].Users, user)
				}
			}
			for _, group := range entriesInA.Groups {
				if !entriesInB.containsGroup(group) {
					difference[permission].Groups = append(difference[permission].Groups, group)
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

func removeProjectPermissions(baseUrl string, projectKey string, token string, permissionSet ProjectPermissions) error {
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

func grantProjectPermissions(baseUrl string, projectKey string, token string, permissionSet ProjectPermissions) error {
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
