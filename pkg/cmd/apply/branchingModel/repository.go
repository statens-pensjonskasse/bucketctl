package branchingModel

import (
	"fmt"
	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"github.com/pterm/pterm"
)

func findRepositoryBranchingModelChanges(desired *RepositoriesProperties, actual *RepositoriesProperties) (toCreate *RepositoriesProperties, toUpdate *RepositoriesProperties, toDelete *RepositoriesProperties) {
	toCreate = new(RepositoriesProperties)
	toUpdate = new(RepositoriesProperties)
	toDelete = new(RepositoriesProperties)
	for repoSlug, repo := range GroupRepositories(desired, actual) {
		bmToCreate, bmToUpdate, bmToDelete := FindBranchingModelsToChange(repo.Desired.BranchingModel, repo.Actual.BranchingModel)
		if bmToCreate != nil {
			*toCreate = append(*toCreate, &RepositoryProperties{RepoSlug: repoSlug, BranchingModel: bmToCreate})
		}
		if bmToUpdate != nil {
			*toUpdate = append(*toUpdate, &RepositoryProperties{RepoSlug: repoSlug, BranchingModel: bmToUpdate})
		}
		if bmToDelete != nil {
			*toDelete = append(*toDelete, &RepositoryProperties{RepoSlug: repoSlug, BranchingModel: bmToDelete})
		}
	}
	return toCreate, toUpdate, toDelete
}

func setRepositoriesBranchingModels(baseUrl string, projectKey string, token string, toCreate *RepositoriesProperties, toUpdate *RepositoriesProperties, toDelete *RepositoriesProperties) error {
	for _, r := range *toDelete {
		if err := deleteRepositoryBranchingModel(baseUrl, projectKey, r.RepoSlug, token, r.BranchingModel); err != nil {
			return err
		}
	}
	for _, r := range *toCreate {
		if err := updateRepositoryBranchingModel(baseUrl, projectKey, r.RepoSlug, token, r.BranchingModel); err != nil {
			return err
		}
	}
	for _, r := range *toUpdate {
		if err := createRepositoryBranchingModel(baseUrl, projectKey, r.RepoSlug, token, r.BranchingModel); err != nil {
			return err
		}
	}

	return nil
}

func createRepositoryBranchingModel(baseUrl string, projectKey string, repoSlug string, token string, branchingModel *BranchingModel) error {
	url := fmt.Sprintf("%s/rest/branch-utils/latest/projects/%s/repos/%s/branchmodel/configuration", baseUrl, projectKey, repoSlug)
	return createBranchingModel(url, token, branchingModel, pterm.Green("🌱 Created"), "repository "+projectKey+"/"+repoSlug)
}

func updateRepositoryBranchingModel(baseUrl string, projectKey string, repoSlug string, token string, branchingModel *BranchingModel) error {
	url := fmt.Sprintf("%s/rest/branch-utils/latest/projects/%s/repos/%s/branchmodel/configuration", baseUrl, projectKey, repoSlug)
	return createBranchingModel(url, token, branchingModel, pterm.Blue("🌿 Updated"), "repository "+projectKey+"/"+repoSlug)
}

func deleteRepositoryBranchingModel(baseUrl string, projectKey string, repoSlug string, token string, branchingModel *BranchingModel) error {
	url := fmt.Sprintf("%s/rest/branch-utils/latest/projects/%s/repos/%s/branchmodel/configuration", baseUrl, projectKey, repoSlug)
	return deleteBranchingModel(url, token, branchingModel, "repository "+projectKey+"/"+repoSlug)

}
