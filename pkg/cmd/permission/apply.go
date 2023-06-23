package permission

import (
	"bucketctl/pkg"
	"bucketctl/pkg/cmd/repository"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sort"
)

var (
	fileName string
)

var applyPermissionsCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("file", cmd.Flags().Lookup("file"))
		viper.BindPFlag("include-repos", cmd.Flags().Lookup("include-repos"))
	},
	Use:  "apply",
	RunE: applyPermissions,
}

func init() {
	applyPermissionsCmd.Flags().StringVarP(&fileName, "file", "f", "", "Permissions file")
	applyPermissionsCmd.Flags().Bool("include-repos", false, "Include repositories")

	applyPermissionsCmd.MarkFlagRequired("file")
}

func applyPermissions(cmd *cobra.Command, args []string) error {
	file := viper.GetString("file")
	baseUrl := viper.GetString("baseUrl")
	limit := viper.GetInt("limit")
	token := viper.GetString("token")
	includeRepos := viper.GetBool("include-repos")

	// Les inn fil (yaml eller json) med ønskede tilganger
	var desiredPermissions map[string]*ProjectPermissions
	if err := pkg.ReadConfigFile(file, &desiredPermissions); err != nil {
		return err
	}

	projectKeys := make([]string, 0, len(desiredPermissions))
	for p := range desiredPermissions {
		projectKeys = append(projectKeys, p)
	}
	sort.Strings(projectKeys)
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(len(desiredPermissions)).WithRemoveWhenDone(true).Start()
	for _, projectKey := range projectKeys {
		progressBar.UpdateTitle(projectKey)
		// Finn gjeldende tilganger
		actualProjectPermissions, err := getProjectPermissions(baseUrl, projectKey, limit, token, includeRepos)
		if err != nil {
			return err
		}
		if includeRepos {
			// Fyller inn data om manglende repositories med tom Permissions struct.
			// Tom Permissions struct betyr ingen ekstra rettigheter på repo-nivå.
			allProjectRepositories, err := repository.GetProjectRepositories(baseUrl, projectKey, limit)
			if err != nil {
				return err
			}
			for repoSlug := range allProjectRepositories {
				if desiredPermissions[projectKey].Repositories == nil {
					desiredPermissions[projectKey].Repositories = make(map[string]*RepositoryPermissions)
				}
				if _, exists := desiredPermissions[projectKey].Repositories[repoSlug]; !exists {
					desiredPermissions[projectKey].Repositories[repoSlug] = &RepositoryPermissions{Permissions: &Permissions{}}
				}
				if actualProjectPermissions.Repositories == nil {
					actualProjectPermissions.Repositories = make(map[string]*RepositoryPermissions)
				}
				if _, exists := actualProjectPermissions.Repositories[repoSlug]; !exists {
					actualProjectPermissions.Repositories[repoSlug] = &RepositoryPermissions{Permissions: &Permissions{}}
				}
			}
		}

		// Finner tilganger i 'actualProjectPermissions' som ikke finnes i 'desiredProjectPermissions'. Disse tilgangene skal fjernes.
		permissionsToBeRemoved := findProjectPermissionsDifference(actualProjectPermissions, desiredPermissions[projectKey])

		//Finner tilganger i 'desiredProjectPermissions' som ikke finnes i 'actualProjectPermissions'. Disse tilgangene skal gis.
		permissionsToBeGranted := findProjectPermissionsDifference(desiredPermissions[projectKey], actualProjectPermissions)

		// Fjern alle ikke-ønskede prosjekt-tilganger
		if err := removeProjectPermissions(baseUrl, projectKey, token, permissionsToBeRemoved); err != nil {
			return err
		}

		// Gi ønskede prosjekt-tilganger
		if err := grantProjectPermissions(baseUrl, projectKey, token, permissionsToBeGranted); err != nil {
			return err
		}

		if includeRepos {
			// Fjern alle ikke-ønskede repo-tilganger for prosjektet
			for repoSlug, permissions := range permissionsToBeRemoved.Repositories {
				if err := removeRepositoryPermissions(baseUrl, projectKey, repoSlug, token, permissions); err != nil {
					return err
				}
			}
			// GI ønskede repo-tilganger for prosjektet
			for repoSlug, permissions := range permissionsToBeGranted.Repositories {
				if err := grantRepositoryPermissions(baseUrl, projectKey, repoSlug, token, permissions); err != nil {
					return err
				}
			}
		}
		progressBar.Increment()
	}

	return nil
}

func findProjectPermissionsDifference(permissionsA *ProjectPermissions, permissionsB *ProjectPermissions) (permissionsDifference *ProjectPermissions) {
	// Finner tilganger i 'permissionsA' som ikke finnes i 'permissionsB'.
	permissionsDifference = &ProjectPermissions{}
	permissionsDifference.Permissions = permissionsA.Permissions.getPermissionsDifference(permissionsB.Permissions)
	if permissionsA.Repositories != nil {
		permissionsDifference.Repositories = make(map[string]*RepositoryPermissions)
		for repoSlug, repo := range permissionsA.Repositories {
			permissionsDifference.Repositories[repoSlug] = new(RepositoryPermissions)
			if permissionsB.Repositories[repoSlug] == nil {
				permissionsDifference.Repositories[repoSlug].Permissions = repo.Permissions
			} else {
				permissionsDifference.Repositories[repoSlug].Permissions = repo.Permissions.getPermissionsDifference(permissionsB.Repositories[repoSlug].Permissions)
			}
		}
	}

	return permissionsDifference
}

// Finner det relative komplementet til setA i setB
func (setA *Permissions) getPermissionsDifference(setB *Permissions) *Permissions {
	difference := make(Permissions)

	for permission := range *setA {
		difference[permission] = new(Entities)
		entriesInA := (*setA)[permission]
		entriesInB, existsInB := (*setB)[permission]

		if !existsInB {
			// Dersom tilgangen ikke finnes i B så legger vi til alle entries i A i det relative komplementet
			difference[permission].Users = (*setA)[permission].Users
			difference[permission].Groups = (*setA)[permission].Groups
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
	return &difference
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

func removeProjectPermissions(baseUrl string, projectKey string, token string, projectPermissions *ProjectPermissions) error {
	for _, entity := range *projectPermissions.Permissions {
		for _, user := range entity.Users {
			if err := removeUserProjectPermissions(baseUrl, projectKey, token, user); err != nil {
				return err
			}
			pterm.Info.Println(pterm.Red("🙅 Revoked") + " permissions for user '" + user + "' for project '" + projectKey + "'")
		}
		for _, group := range entity.Groups {
			if err := removeGroupProjectPermissions(baseUrl, projectKey, token, group); err != nil {
				return err
			}
			pterm.Info.Println(pterm.Red("🚫 Revoked") + " permissions for group '" + group + "' for project '" + projectKey + "'")
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

func removeRepositoryPermissions(baseUrl string, projectKey string, repoSlug string, token string, repositoryPermissions *RepositoryPermissions) error {
	for _, entity := range *repositoryPermissions.Permissions {
		for _, user := range entity.Users {
			if err := removeUserRepositoryPermissions(baseUrl, projectKey, repoSlug, token, user); err != nil {
				return err
			}
			pterm.Info.Println(pterm.Red("🙅 Revoked") + " permissions for user '" + user + "' for repository '" + projectKey + "/" + repoSlug + "'")
		}
		for _, group := range entity.Groups {
			if err := removeGroupRepositoryPermissions(baseUrl, projectKey, repoSlug, token, group); err != nil {
				return err
			}
			pterm.Info.Println(pterm.Red("🚫 Revoked") + " permissions for group '" + group + "' for repository '" + projectKey + "/" + repoSlug + "'")
		}
	}
	return nil
}

func removeUserRepositoryPermissions(baseUrl string, projectKey string, repoSlug string, token string, user string) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/permissions/users", baseUrl, projectKey, repoSlug)
	params := map[string]string{
		"name": user,
	}

	if _, err := pkg.DeleteRequest(url, token, params); err != nil {
		return err
	}
	return nil
}

func removeGroupRepositoryPermissions(baseUrl string, projectKey string, repoSlug string, token string, group string) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/permissions/groups", baseUrl, projectKey, repoSlug)
	params := map[string]string{
		"name": group,
	}

	if _, err := pkg.DeleteRequest(url, token, params); err != nil {
		return err
	}
	return nil
}

func grantProjectPermissions(baseUrl string, projectKey string, token string, projectPermissions *ProjectPermissions) error {
	for permission, entity := range *projectPermissions.Permissions {
		for _, user := range entity.Users {
			if err := grantUserProjectPermission(baseUrl, projectKey, token, user, permission); err != nil {
				return err
			}
			pterm.Info.Println(pterm.Green("🧑 Granted") + " user '" + user + "' permission '" + permission + "' for project '" + projectKey + "'")
		}
		for _, group := range entity.Groups {
			if err := grantGroupProjectPermission(baseUrl, projectKey, token, group, permission); err != nil {
				return err
			}
			pterm.Info.Println(pterm.Green("👥 Granted") + " group '" + group + "' permission '" + permission + "' for project '" + projectKey + "'")
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

	if _, err := pkg.PutRequest(url, token, nil, params); err != nil {
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

	if _, err := pkg.PutRequest(url, token, nil, params); err != nil {
		return err
	}
	return nil
}

func grantRepositoryPermissions(baseUrl string, projectKey string, repoSlug string, token string, repositoryPermissions *RepositoryPermissions) error {
	for permission, entity := range *repositoryPermissions.Permissions {
		for _, user := range entity.Users {
			if err := grantUserRepositoryPermission(baseUrl, projectKey, repoSlug, token, user, permission); err != nil {
				return err
			}
			pterm.Info.Println(pterm.Green("🧑 Granted") + " user '" + user + "' permission '" + permission + "' for repository '" + projectKey + "/" + repoSlug + "'")
		}
		for _, group := range entity.Groups {
			if err := grantGroupRepositoryPermission(baseUrl, projectKey, repoSlug, token, group, permission); err != nil {
				return err
			}
			pterm.Info.Println(pterm.Green("👥 Granted") + " group '" + group + "' permission '" + permission + "' for repository '" + projectKey + "/" + repoSlug + "'")
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

	if _, err := pkg.PutRequest(url, token, nil, params); err != nil {
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

	if _, err := pkg.PutRequest(url, token, nil, params); err != nil {
		return err
	}
	return nil
}
