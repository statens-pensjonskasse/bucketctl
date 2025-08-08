package bitbucket

import (
	"encoding/json"
	"fmt"

	"git.spk.no/infra/bucketctl/pkg/api/bitbucket/types"
	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"git.spk.no/infra/bucketctl/pkg/common"
)

const (
	RefPrefix = "refs/heads/"
)

func GetBranchingModel(baseUrl string, projectKey string, repoSlug string, token string) (*types.BranchingModel, error) {
	url := fmt.Sprintf("%s/rest/branch-utils/latest/projects/%s/repos/%s/branchmodel", baseUrl, projectKey, repoSlug)

	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var branchModel types.BranchingModel
	if err := json.Unmarshal(body, &branchModel); err != nil {
		return nil, err
	}

	return &branchModel, nil
}

func GetProjectBranchingModel(baseUrl string, projectKey string, token string) (*BranchingModel, error) {
	url := fmt.Sprintf("%s/rest/branch-utils/latest/projects/%s/branchmodel/configuration", baseUrl, projectKey)

	bitbucketBranchModel, err := getBranchingModel(url, token)
	if err != nil {
		return nil, err
	}

	return &BranchingModel{
		Development: bitbucketBranchModel.Development,
		Production:  bitbucketBranchModel.Production,
		Types:       FromBitbucketBranchModelTypes(bitbucketBranchModel.Types),
	}, nil
}

func GetRepositoryBranchingModel(baseUrl string, projectKey string, repoSlug string, token string) (*BranchingModel, error) {
	url := fmt.Sprintf("%s/rest/branch-utils/latest/projects/%s/repos/%s/branchmodel/configuration", baseUrl, projectKey, repoSlug)

	bitbucketBranchModel, err := getBranchingModel(url, token)
	if err != nil {
		return nil, err
	}

	if bitbucketBranchModel.Scope.Type != "REPOSITORY" {
		return nil, nil
	}

	return &BranchingModel{
		Development: bitbucketBranchModel.Development,
		Production:  bitbucketBranchModel.Production,
		Types:       FromBitbucketBranchModelTypes(bitbucketBranchModel.Types),
	}, nil
}

func getBranchingModel(url string, token string) (*types.BranchingModel, error) {
	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var branchModel types.BranchingModel
	if err := json.Unmarshal(body, &branchModel); err != nil {
		return nil, err
	}

	return &branchModel, nil
}

func GetRepositoriesBranchingModel(baseUrl string, projectKey string, limit int, token string) (*RepositoriesProperties, error) {
	projectRepositories, err := GetLexicallySortedProjectRepositoriesNames(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	repositoriesProperties := new(RepositoriesProperties)
	for _, repoSlug := range projectRepositories {
		repositoryBranchModel, err := GetRepositoryBranchingModel(baseUrl, projectKey, repoSlug, token)
		if err != nil {
			return nil, err
		}
		*repositoriesProperties = append(*repositoriesProperties, &RepositoryProperties{RepoSlug: repoSlug, BranchingModel: repositoryBranchModel})
	}
	return repositoriesProperties, nil
}
