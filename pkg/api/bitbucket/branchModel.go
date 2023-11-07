package bitbucket

import (
	"bucketctl/pkg/api/bitbucket/types"
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/common"
	"encoding/json"
	"fmt"
)

const (
	RefPrefix = "refs/heads/"
)

func GetBranchModel(baseUrl string, projectKey string, repoSlug string, token string) (*types.BranchModel, error) {
	url := fmt.Sprintf("%s/rest/branch-utils/latest/projects/%s/repos/%s/branchmodel", baseUrl, projectKey, repoSlug)

	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var branchModel types.BranchModel
	if err := json.Unmarshal(body, &branchModel); err != nil {
		return nil, err
	}

	return &branchModel, nil
}

func getRepositoryDefaultBranch(baseUrl string, projectKey string, repoSlug string, token string) (*types.Branch, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/default-branch", baseUrl, projectKey, repoSlug)

	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var defaultBranch types.Branch
	if err := json.Unmarshal(body, &defaultBranch); err != nil {
		return nil, err
	}

	return &defaultBranch, nil
}

func GetProjectBranchModelConfiguration(baseUrl string, projectKey string, token string) (*BranchModel, error) {
	url := fmt.Sprintf("%s/rest/branch-utils/latest/projects/%s/branchmodel/configuration", baseUrl, projectKey)

	bitbucketBranchModel, err := getBranchModelConfiguration(url, token)
	if err != nil {
		return nil, err
	}

	return &BranchModel{
		Development: bitbucketBranchModel.Development,
		Production:  bitbucketBranchModel.Production,
		Types:       FromBitbucketBranchModelTypes(bitbucketBranchModel.Types),
	}, nil
}

func GetRepositoryBranchModelConfiguration(baseUrl string, projectKey string, repoSlug string, token string) (*BranchModel, error) {
	url := fmt.Sprintf("%s/rest/branch-utils/latest/projects/%s/repos/%s/branchmodel/configuration", baseUrl, projectKey, repoSlug)

	bitbucketBranchModel, err := getBranchModelConfiguration(url, token)
	if err != nil {
		return nil, err
	}

	defaultBranch, err := getRepositoryDefaultBranch(baseUrl, projectKey, repoSlug, token)
	if err != nil {
		return nil, err
	}

	if bitbucketBranchModel.Scope.Type != "REPOSITORY" {
		return &BranchModel{DefaultBranch: &defaultBranch.DisplayId}, nil
	}

	return &BranchModel{
		DefaultBranch: &defaultBranch.DisplayId,
		Development:   bitbucketBranchModel.Development,
		Production:    bitbucketBranchModel.Production,
		Types:         FromBitbucketBranchModelTypes(bitbucketBranchModel.Types),
	}, nil
}

func getBranchModelConfiguration(url string, token string) (*types.BranchModel, error) {
	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var branchModel types.BranchModel
	if err := json.Unmarshal(body, &branchModel); err != nil {
		return nil, err
	}

	return &branchModel, nil
}

func GetRepositoriesBranchModel(baseUrl string, projectKey string, limit int, token string) (*RepositoriesProperties, error) {
	projectRepositories, err := GetLexicallySortedProjectRepositoriesNames(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	repositoriesProperties := new(RepositoriesProperties)
	for _, repoSlug := range projectRepositories {
		repositoryBranchModel, err := GetRepositoryBranchModelConfiguration(baseUrl, projectKey, repoSlug, token)
		if err != nil {
			return nil, err
		}
		*repositoriesProperties = append(*repositoriesProperties, &RepositoryProperties{RepoSlug: repoSlug, BranchModel: repositoryBranchModel})
	}
	return repositoriesProperties, nil
}
