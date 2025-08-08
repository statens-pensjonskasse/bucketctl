package bitbucket

import (
	"encoding/json"
	"fmt"

	"git.spk.no/infra/bucketctl/pkg/api/bitbucket/types"
	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"git.spk.no/infra/bucketctl/pkg/common"
	"git.spk.no/infra/bucketctl/pkg/logger"
)

var (
	key  string
	repo string
)

func GetProjectBranchRestrictions(baseUrl string, projectKey string, limit int, token string) (*BranchRestrictions, error) {
	url := fmt.Sprintf("%s/rest/branch-permissions/latest/projects/%s/restrictions?limit=%d", baseUrl, projectKey, limit)

	restrictions, err := getBranchRestrictions(url, token)
	if err != nil {
		return nil, err
	}

	branchRestrictions := &BranchRestrictions{}
	for _, r := range restrictions {
		branchRestrictions.AddRestriction(r)
	}

	return branchRestrictions, nil
}

func GetProjectRepositoriesBranchRestrictions(baseUrl string, projectKey string, limit int, token string) (*RepositoriesProperties, error) {
	projectRepositories, err := GetLexicallySortedProjectRepositoriesNames(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	repositoriesProperties := new(RepositoriesProperties)
	for _, repoSlug := range projectRepositories {
		repositoryBranchRestrictions, err := GetRepositoryBranchRestrictions(baseUrl, projectKey, repoSlug, limit, token)
		if err != nil {
			return nil, err
		}
		*repositoriesProperties = append(*repositoriesProperties, &RepositoryProperties{RepoSlug: repoSlug, BranchRestrictions: repositoryBranchRestrictions})
	}
	return repositoriesProperties, nil
}

func GetRepositoryBranchRestrictions(baseUrl string, projectKey string, repoSlug string, limit int, token string) (*BranchRestrictions, error) {
	url := fmt.Sprintf("%s/rest/branch-permissions/latest/projects/%s/repos/%s/restrictions?limit=%d", baseUrl, projectKey, repoSlug, limit)
	restrictions, err := getBranchRestrictions(url, token)
	if err != nil {
		return nil, err
	}

	branchRestrictions := &BranchRestrictions{}
	for _, r := range restrictions {
		// We also get Project scoped restrictions which we want to ignore as they are implicit
		if r.Scope.Type != "REPOSITORY" {
			continue
		}
		branchRestrictions.AddRestriction(r)
	}

	return branchRestrictions, nil
}

func getBranchRestrictions(url string, token string) ([]*types.Restriction, error) {
	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var restrictions types.RestrictionResponse
	if err := json.Unmarshal(body, &restrictions); err != nil {
		return nil, err
	}

	if !restrictions.IsLastPage {
		logger.Warn("not all restrictions fetched, try with a higher limit")
	}

	return restrictions.Values, nil
}
