package defaultBranch

import (
	"bytes"
	"encoding/json"
	"fmt"

	"git.spk.no/infra/bucketctl/pkg/api/bitbucket"
	"git.spk.no/infra/bucketctl/pkg/api/bitbucket/types"
	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"git.spk.no/infra/bucketctl/pkg/common"
	"git.spk.no/infra/bucketctl/pkg/logger"
	"github.com/pterm/pterm"
)

func setRepositoriesDefaultBranch(baseUrl string, projectKey string, token string, toUpdate *RepositoriesProperties) error {
	for _, r := range *toUpdate {
		branches, err := bitbucket.GetRepositoryBranches(baseUrl, projectKey, r.RepoSlug, token)
		if err != nil {
			return err
		}
		if !branches.ContainsBranchDisplayId(r.DefaultBranch) {
			logger.Warn("repository %s does not have branch called %s, can't update default branch to a non-existent branch", r.RepoSlug, *r.DefaultBranch)
			return nil
		}

		if err := updateRepositoryDefaultBranch(baseUrl, projectKey, r.RepoSlug, token, r.DefaultBranch); err != nil {
			return err
		}
	}
	return nil
}

func updateRepositoryDefaultBranch(baseUrl string, projectKey string, repoSlug string, token string, newDefaultBranch *string) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/default-branch", baseUrl, projectKey, repoSlug)
	return updateDefaultBranch(url, token, newDefaultBranch, "repository "+projectKey+"/"+repoSlug)
}

func updateDefaultBranch(url string, token string, defaultBranch *string, scope string) error {
	if defaultBranch != nil {

		payload, err := json.Marshal(&types.Branch{Id: *defaultBranch})
		if err != nil {
			return err
		}

		if _, err := common.PutRequest(url, token, bytes.NewReader(payload), nil); err != nil {
			return err
		}

		logger.Log("%s default branch to %s in %s", pterm.Blue("🍃 Updated"), *defaultBranch, scope)
	}
	return nil
}
