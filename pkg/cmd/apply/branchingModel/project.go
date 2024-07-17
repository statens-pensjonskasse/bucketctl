package branchingModel

import (
	"fmt"
	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"github.com/pterm/pterm"
)

func setProjectBranchingModel(baseUrl string, projectKey string, token string, toUpdate *BranchingModel) error {
	return updateProjectBranchingModel(baseUrl, projectKey, token, toUpdate)
}

func updateProjectBranchingModel(baseUrl string, projectKey string, token string, branchingModel *BranchingModel) error {
	url := fmt.Sprintf("%s/rest/branch-utils/latest/projects/%s/branchmodel/configuration", baseUrl, projectKey)
	return createBranchingModel(url, token, branchingModel, pterm.Blue("🪴 Updated"), "project "+projectKey)
}
