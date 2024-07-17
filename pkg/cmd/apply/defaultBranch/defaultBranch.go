package defaultBranch

import (
	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"github.com/pterm/pterm"
)

func FindDefaultBranchChanges(desired *ProjectConfigSpec, actual *ProjectConfigSpec) (toUpdate *ProjectConfigSpec) {
	return FindDefaultBranchesToChange(desired, actual)
}

func SetDefaultBranches(baseUrl string, projectKey string, token string, toUpdate *ProjectConfigSpec) error {
	if err := setRepositoriesDefaultBranch(baseUrl, projectKey, token, toUpdate.Repositories); err != nil {
		return err
	}
	return nil
}

func GetChangesAsText(toUpdate *ProjectConfigSpec) (changes []string) {
	changes = append(changes, changesToText(pterm.Blue("🍃 update"), toUpdate)...)
	return changes
}

func changesToText(action string, pcs *ProjectConfigSpec) (changes []string) {
	if pcs.Repositories != nil && len(*pcs.Repositories) > 0 {
		for _, repo := range *pcs.Repositories {
			if repo.DefaultBranch != nil {
				changes = append(changes,
					pterm.Sprintf("%s default branch to %s in repository %s/%s",
						action, *repo.DefaultBranch, pcs.ProjectKey, repo.RepoSlug))
			}
		}
	}
	return changes
}
