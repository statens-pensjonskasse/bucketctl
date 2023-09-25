package access

import (
	. "bucketctl/pkg/api/v1alpha1"
	"github.com/pterm/pterm"
	"strconv"
)

func FindAccessChanges(desired *ProjectConfigSpec, actual *ProjectConfigSpec) (toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) {
	projectAccessToUpdate := FindProjectAccessToUpdate(desired.Access, actual.Access)
	projectPermissionsToCreate, projectPermissionsToDelete := FindPermissionsToChange(desired.Access.Permissions, actual.Access.Permissions)
	repositoriesPermissionsToCreate, repositoriesPermissionsToDelete := findRepositoriesPermissionsChanges(desired.Repositories, actual.Repositories)

	toCreate = &ProjectConfigSpec{ProjectKey: desired.ProjectKey, Access: &ProjectAccess{Permissions: projectPermissionsToCreate}, Repositories: repositoriesPermissionsToCreate}
	toUpdate = &ProjectConfigSpec{ProjectKey: desired.ProjectKey, Access: projectAccessToUpdate, Repositories: new(RepositoriesProperties)}
	toDelete = &ProjectConfigSpec{ProjectKey: desired.ProjectKey, Access: &ProjectAccess{Permissions: projectPermissionsToDelete}, Repositories: repositoriesPermissionsToDelete}

	return toCreate, toUpdate, toDelete
}

func SetAccess(baseUrl string, projectKey string, token string, toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) error {
	if err := setProjectAccessProperties(baseUrl, projectKey, token, toUpdate.Access); err != nil {
		return err
	}
	if err := setProjectPermissions(baseUrl, projectKey, token, toCreate.Access.Permissions, toDelete.Access.Permissions); err != nil {
		return err
	}
	if err := setRepositoriesPermissions(baseUrl, projectKey, token, toCreate.Repositories, toDelete.Repositories); err != nil {
		return err
	}
	return nil
}

func PrintAccessChanges(toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) {
	printAccess(pterm.Green("⭐️ create"), toCreate)
	printAccess(pterm.Blue("🔧 change"), toUpdate)
	printAccess(pterm.Red("❌  remove"), toDelete)
}

func printAccess(action string, pcs *ProjectConfigSpec) {
	if pcs.Access.Public != nil {
		pterm.Printfln("%s public access to %s in project %s",
			action, pterm.Bold.Sprint(strconv.FormatBool(*pcs.Access.Public)), pcs.ProjectKey)
	}
	if pcs.Access.DefaultPermission != nil {
		pterm.Printfln("%s default permission to %s in project %s",
			action, pterm.Bold.Sprint(*pcs.Access.DefaultPermission), pcs.ProjectKey)
	}
	if pcs.Access.Permissions != nil {
		for _, p := range *pcs.Access.Permissions {
			for _, u := range p.Entities.Users {
				pterm.Printfln("%s %s permission for user %s in project %s",
					action, pterm.Bold.Sprint(p.Name), pterm.Bold.Sprint(u), pcs.ProjectKey)
			}
			for _, g := range p.Entities.Groups {
				pterm.Printfln("%s %s permission for group %s in project %s",
					action, pterm.Bold.Sprint(p.Name), pterm.Bold.Sprint(g), pcs.ProjectKey)
			}
		}
	}
	if pcs.Repositories != nil && len(*pcs.Repositories) > 0 {
		for _, repo := range *pcs.Repositories {
			if repo.Permissions != nil {
				for _, p := range *repo.Permissions {
					for _, u := range p.Entities.Users {
						pterm.Printfln("%s %s permission for user %s in repository %s/%s",
							action, pterm.Bold.Sprint(p.Name), pterm.Bold.Sprint(u), pcs.ProjectKey, repo.RepoSlug)
					}
					for _, g := range p.Entities.Groups {
						pterm.Printfln("%s %s permission for group %s in repository %s/%s",
							action, pterm.Bold.Sprint(p.Name), pterm.Bold.Sprint(g), pcs.ProjectKey, repo.RepoSlug)
					}
				}
			}
		}
	}
}
