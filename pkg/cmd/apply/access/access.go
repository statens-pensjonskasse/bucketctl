package access

import (
	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"github.com/pterm/pterm"
	"strconv"
)

func FindAccessChanges(desired *ProjectConfigSpec, actual *ProjectConfigSpec) (toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) {
	projectPermissionsToCreate, projectPermissionsToDelete := FindPermissionsToChange(desired.Permissions, actual.Permissions)
	repositoriesPermissionsToCreate, repositoriesPermissionsToDelete := findRepositoriesPermissionsChanges(desired.Repositories, actual.Repositories)

	toCreate = &ProjectConfigSpec{
		ProjectKey:   desired.ProjectKey,
		Permissions:  projectPermissionsToCreate,
		Repositories: repositoriesPermissionsToCreate}
	toUpdate = &ProjectConfigSpec{
		ProjectKey:        desired.ProjectKey,
		Public:            UpdatePublicProperty(desired, actual),
		DefaultPermission: UpdateDefaultProjectPermissionProperty(desired, actual),
		Permissions:       new(Permissions),
		Repositories:      new(RepositoriesProperties)}
	toDelete = &ProjectConfigSpec{
		ProjectKey:   desired.ProjectKey,
		Permissions:  projectPermissionsToDelete,
		Repositories: repositoriesPermissionsToDelete}

	return toCreate, toUpdate, toDelete
}

func SetAccess(baseUrl string, projectKey string, token string, toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) error {
	if err := setProjectAccessProperties(baseUrl, projectKey, token, toUpdate); err != nil {
		return err
	}
	if err := setProjectPermissions(baseUrl, projectKey, token, toCreate.Permissions, toDelete.Permissions); err != nil {
		return err
	}
	if err := setRepositoriesPermissions(baseUrl, projectKey, token, toCreate.Repositories, toDelete.Repositories); err != nil {
		return err
	}
	return nil
}

func GetChangesAsText(toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) (changes []string) {
	changes = append(changes, changesToText(pterm.Green("⭐️ create"), toCreate)...)
	changes = append(changes, changesToText(pterm.Blue("🔧 change"), toUpdate)...)
	changes = append(changes, changesToText(pterm.Red("🛑 remove"), toDelete)...)
	return changes
}

func changesToText(action string, pcs *ProjectConfigSpec) (changes []string) {
	if pcs.Public != nil {
		changes = append(changes,
			pterm.Sprintf("%s public access to %s in project %s",
				action, strconv.FormatBool(*pcs.Public), pcs.ProjectKey))
	}
	if pcs.DefaultPermission != nil {
		changes = append(changes,
			pterm.Sprintf("%s default permission to %s in project %s",
				action, *pcs.DefaultPermission, pcs.ProjectKey))
	}
	if pcs.Permissions != nil {
		for _, p := range *pcs.Permissions {
			for _, u := range p.Entities.Users {
				changes = append(changes,
					pterm.Sprintf("%s %s permission for user %s in project %s",
						action, p.Name, u, pcs.ProjectKey))
			}
			for _, g := range p.Entities.Groups {
				changes = append(changes,
					pterm.Sprintf("%s %s permission for group %s in project %s",
						action, p.Name, g, pcs.ProjectKey))
			}
		}
	}
	if pcs.Repositories != nil && len(*pcs.Repositories) > 0 {
		for _, repo := range *pcs.Repositories {
			if repo.Permissions != nil {
				for _, p := range *repo.Permissions {
					for _, u := range p.Entities.Users {
						changes = append(changes,
							pterm.Sprintf("%s %s permission for user %s in repository %s/%s",
								action, p.Name, u, pcs.ProjectKey, repo.RepoSlug))
					}
					for _, g := range p.Entities.Groups {
						changes = append(changes,
							pterm.Sprintf("%s %s permission for group %s in repository %s/%s",
								action, p.Name, g, pcs.ProjectKey, repo.RepoSlug))
					}
				}
			}
		}
	}
	return changes
}
