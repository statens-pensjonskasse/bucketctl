package pullRequest

import (
	"bucketctl/pkg/api/bitbucket"
	"bucketctl/pkg/api/bitbucket/types"
	"bucketctl/pkg/common"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"strings"
)

var (
	title             string
	description       string
	reviewerUsernames *[]string
	fromBranch        string
	toBranch          string
)

var createCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag(common.ProjectKeyFlag, cmd.Flags().Lookup(common.ProjectKeyFlag))
		viper.BindPFlag(common.RepoSlugFlag, cmd.Flags().Lookup(common.RepoSlugFlag))
		viper.BindPFlag(common.NoDefaultReviewersFlag, cmd.Flags().Lookup(common.NoDefaultReviewersFlag))
	},
	Use:   "create",
	Short: "Create a pull request for current workdir repository towards the default branch",
	RunE:  createPullRequest,
}

func init() {
	createCmd.Flags().StringVarP(&title, common.PullRequestTitleFlag, common.PullRequestTitleFlagShorthand, "", "Pull request title")
	createCmd.Flags().StringVarP(&description, common.PullRequestDescriptionFlag, common.PullRequestDescriptionFlagShorthand, "", "Pull request description")
	createCmd.Flags().StringVar(&fromBranch, common.PullRequestFromBranchFlag, "", "From branch (current branch)")
	createCmd.Flags().StringVar(&toBranch, common.PullRequestToBranchFlag, "", "To branch (default branch)")
	createCmd.Flags().Bool(common.NoDefaultReviewersFlag, false, "Don't add default reviewers")
	reviewerUsernames = createCmd.Flags().StringSliceP(common.PullRequestReviewerFlag, common.PullRequestReviewerFlagShorthand, []string{}, "Reviewer username (e-mail)")
}

func createPullRequest(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString(common.BaseUrlFlag)
	projectKey := viper.GetString(common.ProjectKeyFlag)
	repoSlug := viper.GetString(common.RepoSlugFlag)
	token := viper.GetString(common.TokenFlag)

	// Find permission key and repo slug from origin url if not given
	if projectKey == "" || repoSlug == "" {
		remoteUrl, err := getRemoteOriginUrl("./")
		// Remote URL is on the format
		// ssh://<url>/<projectKey>/<repoSlug>.git
		// or
		// https://<url>/scm/<projectKey>/<repoSlug>.git
		remoteUrlSlice := strings.Split(remoteUrl, "/")
		cobra.CheckErr(err)
		if projectKey == "" {
			projectKey = remoteUrlSlice[len(remoteUrlSlice)-2]
		}
		if repoSlug == "" {
			repoSlug = strings.TrimSuffix(remoteUrlSlice[len(remoteUrlSlice)-1], ".git")
		}
	}

	branchModel, err := bitbucket.GetBranchingModel(baseUrl, projectKey, repoSlug, token)
	cobra.CheckErr(err)
	fromRef, err := getFromRef()
	cobra.CheckErr(err)
	toRef, err := getToRef(baseUrl, projectKey, repoSlug, token)
	cobra.CheckErr(err)

	if title == "" {
		title = getTitle(fromRef)
	}

	if description == "" {
		description, err = getDescription(baseUrl, projectKey, repoSlug, token, fromRef, toRef)
		cobra.CheckErr(err)
	}

	if !viper.GetBool(common.NoDefaultReviewersFlag) {
		// Add default reviewers for given source and target branch
		defaultReviewers, err := getDefaultReviewers(baseUrl, projectKey, repoSlug, token)
		cobra.CheckErr(err)
		defaultReviewerSlugs := getDefaultReviewerSlugsForSourceAndTargetRefs(defaultReviewers, branchModel, fromBranch, toBranch)
		for _, slug := range defaultReviewerSlugs {
			*reviewerUsernames = append(*reviewerUsernames, slug)
		}
	}

	uniqueReviewers := make(map[string]struct{})
	for _, username := range *reviewerUsernames {
		uniqueReviewers[username] = struct{}{}
	}
	var reviewers []*types.PullRequestParticipant
	for username := range uniqueReviewers {
		reviewers = append(reviewers, &types.PullRequestParticipant{User: &types.User{Name: username}})
	}

	r := &types.Repository{Slug: repoSlug, Project: &types.Project{Key: projectKey}}

	pullRequest := &types.PullRequest{
		Title:       title,
		Description: description,
		State:       "OPEN",
		Open:        true,
		Closed:      false,
		Locked:      false,
		FromRef:     &types.Ref{Id: fromRef, Repository: r},
		ToRef:       &types.Ref{Id: toRef, Repository: r},
		Reviewers:   reviewers,
	}

	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/pull-requests", baseUrl, projectKey, repoSlug)

	payload, err := json.Marshal(pullRequest)
	cobra.CheckErr(err)

	resp, err := common.PostRequest(url, token, bytes.NewReader(payload), nil)
	cobra.CheckErr(err)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	cobra.CheckErr(err)

	var prReply types.PullRequestInfo
	if err := json.Unmarshal(body, &prReply); err != nil {
		cobra.CheckErr(err)
	}

	pterm.Info.Println("🧑‍💻New pull request created. Go to", prReply.Links.Self[0].Href, "for details")

	return err
}

func getFromRef() (string, error) {
	if fromBranch == "" {
		currentBranch, err := getBranchName("./")
		if err != nil {
			return "", err
		}
		fromBranch = currentBranch
	}
	if !strings.HasPrefix(fromBranch, bitbucket.RefPrefix) {
		fromBranch = bitbucket.RefPrefix + fromBranch
	}
	return fromBranch, nil
}

func getToRef(baseUrl string, projectKey string, repoSlug string, token string) (string, error) {
	if toBranch == "" {
		defaultBranch, err := bitbucket.GetDefaultBranch(baseUrl, projectKey, repoSlug, token)
		if err != nil {
			return "", err
		}
		toBranch = defaultBranch.Id
	}
	if !strings.HasPrefix(toBranch, bitbucket.RefPrefix) {
		toBranch = bitbucket.RefPrefix + toBranch
	}
	return toBranch, nil
}

func getTitle(fromRef string) string {
	return strings.TrimPrefix(fromRef, bitbucket.RefPrefix)
}

func getDescription(baseUrl string, projectKey string, repoSlug string, token string, fromRef string, toRef string) (string, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/commits", baseUrl, projectKey, repoSlug)
	params := map[string]string{
		"limit":         "10",
		"merges":        "exclude",
		"ignoreMissing": "true",
		"since":         strings.TrimPrefix(toRef, bitbucket.RefPrefix),
		"until":         strings.TrimPrefix(fromRef, bitbucket.RefPrefix),
	}
	messages, err := getCommitMessagesBetween(url, token, params)
	if err != nil {
		return "", err
	}

	var combinedMessages string
	for _, m := range messages {
		combinedMessages += "* " + m + "\n"
	}

	return combinedMessages, nil
}

func getDefaultReviewerSlugsForSourceAndTargetRefs(defaultReviewers []*types.DefaultReviewers, branchModel *types.BranchingModel, sourceRef string, targetRef string) []string {
	var usernames []string
	if defaultReviewers == nil || len(defaultReviewers) == 0 {
		return usernames
	}

	for _, review := range defaultReviewers {
		if common.RefMatcher(review.SourceRefMatcher, branchModel, sourceRef) && common.RefMatcher(review.TargetRefMatcher, branchModel, targetRef) {
			for _, u := range review.Reviewers {
				usernames = append(usernames, u.Name)
			}
		}
	}
	return usernames
}

func getCommitMessagesBetween(url string, token string, params map[string]string) ([]string, error) {
	resp, err := common.HttpRequest("GET", url, nil, token, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var commitResp types.CommitsResponse
	if err := json.Unmarshal(body, &commitResp); err != nil {
		return nil, err
	}

	var commitMessages []string
	for _, c := range commitResp.Values {
		commitMessages = append(commitMessages, c.Message)
	}

	return commitMessages, nil
}
