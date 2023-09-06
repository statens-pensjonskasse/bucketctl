package branch

import (
	"bucketctl/pkg/common"
	"bucketctl/pkg/types"
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
