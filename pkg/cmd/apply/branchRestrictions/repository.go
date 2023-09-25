package branchRestrictions

import (
	. "bucketctl/pkg/api/v1alpha1"
	"fmt"
	"github.com/pterm/pterm"
)

func findRepositoriesBranchRestrictionsChanges(desired *RepositoriesProperties, actual *RepositoriesProperties) (toCreate *RepositoriesProperties, toUpdate *RepositoriesProperties, toDelete *RepositoriesProperties) {
	toCreate = new(RepositoriesProperties)
	toUpdate = new(RepositoriesProperties)
	toDelete = new(RepositoriesProperties)
	for repoSlug, repo := range GroupRepositories(desired, actual) {
		brToCreate, brToUpdate, brToDelete := FindBranchRestrictionsToChange(repo.Desired.BranchRestrictions, repo.Actual.BranchRestrictions)
		if len(*brToCreate) > 0 {
			*toCreate = append(*toCreate, &RepositoryProperties{RepoSlug: repoSlug, BranchRestrictions: brToCreate})
		}
		if len(*brToUpdate) > 0 {
			*toUpdate = append(*toUpdate, &RepositoryProperties{RepoSlug: repoSlug, BranchRestrictions: brToUpdate})
		}
		if len(*brToDelete) > 0 {
			*toDelete = append(*toDelete, &RepositoryProperties{RepoSlug: repoSlug, BranchRestrictions: brToDelete})
		}
	}
	return toCreate, toUpdate, toDelete
}

func setRepositoriesBranchRestrictions(baseUrl string, projectKey string, token string, toCreate *RepositoriesProperties, toUpdate *RepositoriesProperties, toDelete *RepositoriesProperties) error {
	for _, r := range *toDelete {
		if err := deleteRepositoryBranchRestrictions(baseUrl, projectKey, r.RepoSlug, token, r.BranchRestrictions); err != nil {
			return err
		}
	}
	for _, r := range *toUpdate {
		if err := updateRepositoryBranchRestrictions(baseUrl, projectKey, r.RepoSlug, token, r.BranchRestrictions); err != nil {
			return err
		}
	}
	for _, r := range *toCreate {
		if err := createRepositoryBranchRestrictions(baseUrl, projectKey, r.RepoSlug, token, r.BranchRestrictions); err != nil {
			return err
		}
	}
	return nil
}

func createRepositoryBranchRestrictions(baseUrl string, projectKey string, repoSlug string, token string, restrictions *BranchRestrictions) error {
	url := fmt.Sprintf("%s/rest/branch-permissions/latest/projects/%s/repos/%s/restrictions", baseUrl, projectKey, repoSlug)
	return createBranchRestrictions(url, token, restrictions, pterm.Green("🪵️ Created"), "repository "+projectKey+"/"+repoSlug)
}

func updateRepositoryBranchRestrictions(baseUrl string, projectKey string, repoSlug string, token string, restrictions *BranchRestrictions) error {
	url := fmt.Sprintf("%s/rest/branch-permissions/latest/projects/%s/repos/%s/restrictions", baseUrl, projectKey, repoSlug)
	return createBranchRestrictions(url, token, restrictions, pterm.Blue("⚙️ Configured"), "repository "+projectKey+"/"+repoSlug)
}

func deleteRepositoryBranchRestrictions(baseUrl string, projectKey string, repoSlug string, token string, restrictions *BranchRestrictions) error {
	url := fmt.Sprintf("%s/rest/branch-permissions/latest/projects/%s/repos/%s/restrictions", baseUrl, projectKey, repoSlug)
	return deleteRestrictions(url, token, restrictions, "repository "+projectKey+"/"+repoSlug)
}
