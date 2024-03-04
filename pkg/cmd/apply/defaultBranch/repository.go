package defaultBranch

import (
	"bucketctl/pkg/api/bitbucket/types"
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/common"
	"bucketctl/pkg/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
)

func setRepositoriesDefaultBranch(baseUrl string, projectKey string, token string, toUpdate *RepositoriesProperties) error {
	for _, r := range *toUpdate {
		if err := updateRepositoryDefaultBranch(baseUrl, projectKey, r.RepoSlug, token, r.DefaultBranch); err != nil {
			return err
		}
	}
	return nil
}

func updateRepositoryDefaultBranch(baseUrl string, projectKey string, repoSlug string, token string, defaultBranch *string) error {
	url := fmt.Sprintf("%s/projects/%s/repos/%s/default-branch", baseUrl, projectKey, repoSlug)
	return updateDefaultBranch(url, token, defaultBranch, "repository "+projectKey+"/"+repoSlug)
}

func updateDefaultBranch(url string, token string, defaultBranch *string, scope string) error {
	if defaultBranch != nil {

		payload, err := json.Marshal(&types.Branch{Id: *defaultBranch})
		if err != nil {
			return err
		}

		if resp, err := common.PutRequest(url, token, bytes.NewReader(payload), nil); err != nil {
			if resp.StatusCode == 404 {
				logger.Warn("%s not found, can't update default branch to %s", url, *defaultBranch)
				return nil
			}
			return err
		}

		logger.Log("%s default branch to %s in %s", pterm.Blue("🍃 Updated"), *defaultBranch, scope)
	}
	return nil
}
