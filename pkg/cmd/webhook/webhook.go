package webhook

import (
	"bucketctl/pkg/cmd/repository"
	"bucketctl/pkg/common"
	"bucketctl/pkg/types"
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"io"
	"sort"
	"strconv"
)

type ProjectWebhooks struct {
	Webhooks     []*types.Webhook               `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`
	Repositories map[string]*RepositoryWebhooks `json:"repositories,omitempty" yaml:"repositories,omitempty"`
}

type RepositoryWebhooks struct {
	Webhooks []*types.Webhook `yaml:"webhooks"`
}

var (
	key  string
	repo string
)

var Cmd = &cobra.Command{
	Use:     "webhook",
	Short:   "View and edit repository and project webhooks",
	Aliases: []string{"wh"},
}

func init() {
	Cmd.AddCommand(applyWebhooksCmd)
	Cmd.AddCommand(listAllWebhooksCmd)
	Cmd.AddCommand(listWebhooksCmd)
}

func getWebhook(url string, token string) (*types.Webhook, error) {
	resp, err := common.GetRequest(url, token)
	if err != nil {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var webhook types.Webhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		return nil, err
	}

	return &webhook, nil
}

func getProjectWebhook(baseUrl string, projectKey string, webhookId int, limit int, token string) (*types.Webhook, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/webhooks/%d?limit=%d", baseUrl, projectKey, webhookId, limit)
	return getWebhook(url, token)
}

func getRepositoryWebhook(baseUrl string, projectKey string, repoSlug string, webhookId int, limit int, token string) (*types.Webhook, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/webhooks/%d?limit=%d", baseUrl, projectKey, repoSlug, webhookId, limit)
	return getWebhook(url, token)
}

func getWebhooks(url string, token string) ([]*types.Webhook, error) {
	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var webhooks types.WebhooksResponse
	if err := json.Unmarshal(body, &webhooks); err != nil {
		return nil, err
	}

	if !webhooks.IsLastPage {
		pterm.Warning.Println("Not all webhooks fetched, try with a higher limit")
	}

	return webhooks.Values, nil
}

func getProjectWebhooks(baseUrl string, projectKey string, limit int, token string, includeRepos bool) (*ProjectWebhooks, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/webhooks?limit=%d", baseUrl, projectKey, limit)

	webhooks, err := getWebhooks(url, token)
	if err != nil {
		return nil, err
	}

	projectWebhooks := &ProjectWebhooks{Webhooks: webhooks}

	if includeRepos {
		projectRepositories, err := repository.GetProjectRepositories(baseUrl, projectKey, token, limit)
		if err != nil {
			return nil, err
		}
		projectWebhooks.Repositories = make(map[string]*RepositoryWebhooks)
		for repoSlug := range projectRepositories {
			repoWebhooks, err := getRepositoryWebhooks(baseUrl, projectKey, repoSlug, limit, token)
			if err != nil {
				return nil, err
			}
			if len(repoWebhooks.Webhooks) > 0 {
				projectWebhooks.Repositories[repoSlug] = repoWebhooks
			}
		}
	}

	return projectWebhooks, nil
}

func getRepositoryWebhooks(baseUrl string, projectKey string, repoSlug string, limit int, token string) (*RepositoryWebhooks, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/webhooks?limit=%d", baseUrl, projectKey, repoSlug, limit)

	webhooks, err := getWebhooks(url, token)
	if err != nil {
		return nil, err
	}

	return &RepositoryWebhooks{Webhooks: webhooks}, nil
}

func prettyFormatProjectWebhooks(projectWebhooksMap map[string]*ProjectWebhooks) [][]string {
	var data [][]string
	data = append(data, []string{"Project", "Repository", "ID", "Name", "Events", "URL", "Active", "Verify SSL"})

	projects := common.GetLexicallySortedKeys(projectWebhooksMap)
	for _, projectKey := range projects {
		formattedProjectWebhooks := prettyFormatWebhooks(projectKey, "PROJECT", projectWebhooksMap[projectKey].Webhooks)
		data = append(data, formattedProjectWebhooks...)

		// Sorter repoene alfabetisk
		repos := make([]string, 0, len(projectWebhooksMap[projectKey].Repositories))
		for r := range projectWebhooksMap[projectKey].Repositories {
			repos = append(repos, r)
		}
		sort.Strings(repos)
		for _, repoSlug := range repos {
			formattedRepoWebhooks := prettyFormatWebhooks(projectKey, repoSlug, projectWebhooksMap[projectKey].Repositories[repoSlug].Webhooks)
			data = append(data, formattedRepoWebhooks...)
		}
	}

	return data
}

func prettyFormatWebhooks(projectKey string, repoSlug string, webhooks []*types.Webhook) [][]string {
	var data [][]string

	if webhooks == nil {
		return data
	}

	for _, webhook := range webhooks {
		sort.Strings(webhook.Events)
		data = append(data, []string{projectKey, repoSlug, strconv.Itoa(webhook.Id), webhook.Name, webhook.Events[0], webhook.Url, strconv.FormatBool(webhook.Active), strconv.FormatBool(webhook.SslVerificationRequired)})
		for _, event := range webhook.Events[1:] {
			data = append(data, []string{"", "", "", "", event, "", "", ""})
		}
	}

	return data
}
