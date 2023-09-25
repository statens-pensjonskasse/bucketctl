package webhooks

import (
	. "bucketctl/pkg/api/v1alpha1"
	"github.com/pterm/pterm"
)

func FindWebhooksChanges(desired *ProjectConfigSpec, actual *ProjectConfigSpec) (toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) {
	projToCreate, projToUpdate, projToDelete := FindWebhooksToChange(desired.Webhooks, actual.Webhooks)
	repoToCreate, repoToUpdate, repoToDelete := findRepositoriesWebhookChanges(desired.Repositories, actual.Repositories)

	toCreate = &ProjectConfigSpec{ProjectKey: desired.ProjectKey, Webhooks: projToCreate, Repositories: repoToCreate}
	toUpdate = &ProjectConfigSpec{ProjectKey: desired.ProjectKey, Webhooks: projToUpdate, Repositories: repoToUpdate}
	toDelete = &ProjectConfigSpec{ProjectKey: desired.ProjectKey, Webhooks: projToDelete, Repositories: repoToDelete}

	return toCreate, toUpdate, toDelete
}

func SetWebhooks(baseUrl string, projectKey string, token string, toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) error {
	if err := setProjectWebhooks(baseUrl, projectKey, token, toCreate.Webhooks, toUpdate.Webhooks, toDelete.Webhooks); err != nil {
		return err
	}

	if err := setRepositoriesWebhooks(baseUrl, projectKey, token, toCreate.Repositories, toUpdate.Repositories, toDelete.Repositories); err != nil {
		return err
	}

	return nil
}

func PrintWebhookChanges(toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) {
	printWebhook(pterm.Green("🪝 create"), toCreate)
	printWebhook(pterm.Blue("🎣 update"), toUpdate)
	printWebhook(pterm.Red("🛑 delete"), toDelete)
}

func printWebhook(action string, pcs *ProjectConfigSpec) {
	if pcs.Webhooks != nil {
		for _, wh := range *pcs.Webhooks {
			pterm.Printfln("%s %s webhook in project %s",
				action, pterm.Bold.Sprint(wh.Name), pcs.ProjectKey)
		}
	}
	if pcs.Repositories != nil {
		for _, repo := range *pcs.Repositories {
			if repo.Webhooks != nil {
				for _, wh := range *repo.Webhooks {
					pterm.Printfln("%s %s webhook in repository %s/%s",
						action, pterm.Bold.Sprint(wh.Name), pcs.ProjectKey, repo.RepoSlug)
				}
			}
		}
	}

}
