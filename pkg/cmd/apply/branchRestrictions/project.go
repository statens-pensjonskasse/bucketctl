package branchRestrictions

import (
	"fmt"

	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"github.com/pterm/pterm"
)

func setProjectBranchRestrictions(baseUrl string, projectKey string, token string, toCreate *BranchRestrictions, toUpdate *BranchRestrictions, toDelete *BranchRestrictions) error {
	if err := createProjectBranchRestrictions(baseUrl, projectKey, token, toCreate); err != nil {
		return err
	}
	if err := updateProjectBranchRestrictions(baseUrl, projectKey, token, toUpdate); err != nil {
		return err
	}
	if err := deleteProjectBranchRestrictions(baseUrl, projectKey, token, toDelete); err != nil {
		return err
	}

	return nil
}

func createProjectBranchRestrictions(baseUrl string, projectKey string, token string, restrictions *BranchRestrictions) error {
	url := fmt.Sprintf("%s/rest/branch-permissions/latest/projects/%s/restrictions", baseUrl, projectKey)
	return createBranchRestrictions(url, token, restrictions, pterm.Green("🪵️Created"), "project "+projectKey)
}

func updateProjectBranchRestrictions(baseUrl string, projectKey string, token string, restrictions *BranchRestrictions) error {
	url := fmt.Sprintf("%s/rest/branch-permissions/latest/projects/%s/restrictions", baseUrl, projectKey)
	return createBranchRestrictions(url, token, restrictions, pterm.Blue("⚙️ Configured"), "project "+projectKey)
}

func deleteProjectBranchRestrictions(baseUrl string, projectKey string, token string, restrictions *BranchRestrictions) error {
	url := fmt.Sprintf("%s/rest/branch-permissions/latest/projects/%s/restrictions", baseUrl, projectKey)
	return deleteRestrictions(url, token, restrictions, "project '"+projectKey+"'")
}
