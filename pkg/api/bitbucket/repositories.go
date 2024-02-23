package bitbucket

import (
	types2 "bucketctl/pkg/api/bitbucket/types"
	"bucketctl/pkg/common"
	"bucketctl/pkg/logger"
	"encoding/json"
	"fmt"
)

func GetProjectRepositories(baseUrl string, projectKey string, limit int, token string) ([]*types2.Repository, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos?limit=%d", baseUrl, projectKey, limit)

	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var repoResponse types2.RepositoriesResponse
	if err := json.Unmarshal(body, &repoResponse); err != nil {
		return nil, err
	}

	if !repoResponse.IsLastPage {
		logger.Warn("Not all repositories fetched, try with a higher limit")
	}
	return repoResponse.Values, nil
}

func GetProjectRepositoriesMap(baseUrl string, projectKey string, limit int, token string) (map[string]*types2.Repository, error) {
	repositoriesList, err := GetProjectRepositories(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	repositories := make(map[string]*types2.Repository)
	for _, r := range repositoriesList {
		if r.Archived {
			continue
		}
		repositories[r.Slug] = &types2.Repository{
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

func GetRepository(baseUrl string, projectKey string, repoSlug string, token string) (*types2.Repository, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s", baseUrl, projectKey, repoSlug)

	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}
	var repoResponse types2.Repository
	if err := json.Unmarshal(body, &repoResponse); err != nil {
		return nil, err
	}
	return &types2.Repository{
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

func GetDefaultBranch(baseUrl string, projectKey string, repoSlug string, token string) (*types2.Ref, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/default-branch", baseUrl, projectKey, repoSlug)

	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var branch *types2.Ref
	if err := json.Unmarshal(body, &branch); err != nil {
		return nil, err
	}

	return branch, nil
}
