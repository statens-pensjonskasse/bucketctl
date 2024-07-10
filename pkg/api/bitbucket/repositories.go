package bitbucket

import (
	"bucketctl/pkg/api/bitbucket/types"
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/common"
	"bucketctl/pkg/logger"
	"encoding/json"
	"fmt"
)

func GetProjectRepositories(baseUrl string, projectKey string, limit int, token string) ([]*types.Repository, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos?limit=%d", baseUrl, projectKey, limit)

	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var repoResponse types.RepositoriesResponse
	if err := json.Unmarshal(body, &repoResponse); err != nil {
		return nil, err
	}

	if !repoResponse.IsLastPage {
		logger.Warn("not all repositories fetched, try with a higher limit")
	}
	return repoResponse.Values, nil
}

func GetProjectRepositoriesMap(baseUrl string, projectKey string, limit int, token string) (map[string]*types.Repository, error) {
	repositoriesList, err := GetProjectRepositories(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	repositories := make(map[string]*types.Repository)
	for _, r := range repositoriesList {
		if r.Archived {
			continue
		}
		repositories[r.Slug] = &types.Repository{
			Id:            r.Id,
			Name:          r.Name,
			Slug:          r.Slug,
			HierarchyId:   r.HierarchyId,
			ScmId:         r.ScmId,
			State:         r.State,
			StatusMessage: r.StatusMessage,
			Forkable:      r.Forkable,
			Public:        r.Public,
			Archived:      r.Archived,
		}
	}

	return repositories, nil
}

func GetLexicallySortedProjectRepositoriesNames(baseUrl string, projectKey string, limit int, token string) ([]string, error) {
	repositoriesMap, err := GetProjectRepositoriesMap(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}
	return common.GetLexicallySortedKeys(repositoriesMap), nil
}

func GetRepository(baseUrl string, projectKey string, repoSlug string, token string) (*types.Repository, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s", baseUrl, projectKey, repoSlug)

	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}
	var repoResponse types.Repository
	if err := json.Unmarshal(body, &repoResponse); err != nil {
		return nil, err
	}
	return &types.Repository{
		Id:            repoResponse.Id,
		Name:          repoResponse.Name,
		Slug:          repoResponse.Slug,
		HierarchyId:   repoResponse.HierarchyId,
		ScmId:         repoResponse.ScmId,
		State:         repoResponse.State,
		StatusMessage: repoResponse.StatusMessage,
		Forkable:      repoResponse.Forkable,
		Public:        repoResponse.Public,
		Archived:      repoResponse.Archived,
	}, nil
}

func GetDefaultBranch(baseUrl string, projectKey string, repoSlug string, token string) (*types.Ref, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/default-branch", baseUrl, projectKey, repoSlug)

	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var branch *types.Ref
	if err := json.Unmarshal(body, &branch); err != nil {
		return nil, err
	}

	return branch, nil
}

func GetRepositoryBranches(baseUrl string, projectKey string, repoSlug string, token string) (*Branches, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/branches", baseUrl, projectKey, repoSlug)

	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var branchesResponse types.BranchesResponse
	if err := json.Unmarshal(body, &branchesResponse); err != nil {
		return nil, err
	}

	if !branchesResponse.IsLastPage {
		logger.Warn("not all branches fetched, try with a higher limit")
	}

	branches := new(Branches)

	for _, branch := range branchesResponse.Values {
		*branches = append(*branches, FromBitbucketBranch(branch))
	}

	return branches, nil
}
