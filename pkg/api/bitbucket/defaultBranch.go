package bitbucket

import (
	"bucketctl/pkg/api/bitbucket/types"
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/common"
	"encoding/json"
	"fmt"
)

func GetRepositoryDefaultBranch(baseUrl string, projectKey string, repoSlug string, token string) (*types.Branch, error) {
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

func GetRepositoriesDefaultBranch(baseUrl string, projectKey string, limit int, token string) (*RepositoriesProperties, error) {
	projectRepositories, err := GetLexicallySortedProjectRepositoriesNames(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	repositoriesProperties := new(RepositoriesProperties)
	for _, repoSlug := range projectRepositories {
		repoDefaultBranch, err := GetDefaultBranch(baseUrl, projectKey, repoSlug, token)
		if err != nil {
			return nil, err
		}
		*repositoriesProperties = append(*repositoriesProperties, &RepositoryProperties{RepoSlug: repoSlug, DefaultBranch: &repoDefaultBranch.DisplayId})
	}
	return repositoriesProperties, nil
}
