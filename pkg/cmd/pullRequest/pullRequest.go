package pullRequest

import (
	"encoding/json"
	"fmt"
	"git.spk.no/infra/bucketctl/pkg/api/bitbucket/types"
	"git.spk.no/infra/bucketctl/pkg/common"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var (
	key  string
	repo string
)

var Cmd = &cobra.Command{
	Use:     "pullRequest",
	Short:   "Pull request commands",
	Aliases: []string{"pr"},
}

func init() {
	Cmd.MarkFlagRequired(common.BaseUrlFlag)
	Cmd.PersistentFlags().StringVarP(&key, common.ProjectKeyFlag, common.ProjectKeyFlagShorthand, "", "Project key")
	Cmd.PersistentFlags().StringVarP(&repo, common.RepoSlugFlag, common.RepoSlugFlagShorthand, "", "Repository slug")

	Cmd.AddCommand(createCmd)
}

func getDefaultReviewers(baseUrl string, projectKey string, repoSlug string, token string) ([]*types.DefaultReviewers, error) {
	url := fmt.Sprintf("%s/rest/default-reviewers/latest/projects/%s/repos/%s/conditions", baseUrl, projectKey, repoSlug)

	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var reviewers []*types.DefaultReviewers
	if err := json.Unmarshal(body, &reviewers); err != nil {
		return nil, err
	}

	return reviewers, nil
}

func getRemoteOriginUrl(repoPath string) (string, error) {
	gitRepo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}

	name, err := gitRepo.Remote("origin")
	if err != nil {
		return "", err
	}

	return name.Config().URLs[0], nil
}

func getBranchName(repoPath string) (string, error) {
	gitRepo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}

	head, err := gitRepo.Head()
	if err != nil {
		return "", err
	}
	return head.Name().String(), nil
}

func getLastCommitMessage(repoPath string) (string, error) {
	gitRepo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}

	head, err := gitRepo.Head()
	if err != nil {
		return "", err
	}

	commit, err := gitRepo.CommitObject(head.Hash())

	return commit.Message, nil
}
