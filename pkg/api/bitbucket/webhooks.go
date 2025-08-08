package bitbucket

import (
	"encoding/json"
	"fmt"
	"io"

	"git.spk.no/infra/bucketctl/pkg/api/bitbucket/types"
	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"git.spk.no/infra/bucketctl/pkg/common"
	"git.spk.no/infra/bucketctl/pkg/logger"
)

func getWebhook(url string, token string) (*Webhook, error) {
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

	return FromBitbucketWebhook(&webhook), nil
}

func getProjectWebhook(baseUrl string, projectKey string, webhookId int, limit int, token string) (*Webhook, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/webhooks/%d?limit=%d", baseUrl, projectKey, webhookId, limit)
	return getWebhook(url, token)
}

func getRepositoryWebhook(baseUrl string, projectKey string, repoSlug string, webhookId int, limit int, token string) (*Webhook, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/webhooks/%d?limit=%d", baseUrl, projectKey, repoSlug, webhookId, limit)
	return getWebhook(url, token)
}

func getWebhooks(url string, token string) (*Webhooks, error) {
	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var webhooksResponse types.WebhooksResponse
	if err := json.Unmarshal(body, &webhooksResponse); err != nil {
		return nil, err
	}

	if !webhooksResponse.IsLastPage {
		logger.Warn("not all webhooks fetched, try with a higher limit")
	}

	webhooks := new(Webhooks)

	for _, wh := range webhooksResponse.Values {
		*webhooks = append(*webhooks, FromBitbucketWebhook(wh))
	}

	return webhooks, nil
}

func GetProjectWebhooks(baseUrl string, projectKey string, limit int, token string) (*Webhooks, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/webhooks?limit=%d", baseUrl, projectKey, limit)

	webhooks, err := getWebhooks(url, token)
	if err != nil {
		return nil, err
	}

	return webhooks, nil
}

func GetProjectRepositoriesWebhooks(baseUrl string, projectKey string, limit int, token string) (*RepositoriesProperties, error) {
	projectRepositories, err := GetLexicallySortedProjectRepositoriesNames(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	repositoriesProperties := new(RepositoriesProperties)
	for _, repoSlug := range projectRepositories {
		repoWebhooks, err := GetRepositoryWebhooks(baseUrl, projectKey, repoSlug, limit, token)
		if err != nil {
			return nil, err
		}
		*repositoriesProperties = append(*repositoriesProperties, &RepositoryProperties{RepoSlug: repoSlug, Webhooks: repoWebhooks})
	}

	return repositoriesProperties, nil
}

func GetRepositoryWebhooks(baseUrl string, projectKey string, repoSlug string, limit int, token string) (*Webhooks, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/webhooks?limit=%d", baseUrl, projectKey, repoSlug, limit)

	webhooks, err := getWebhooks(url, token)
	if err != nil {
		return nil, err
	}

	return webhooks, nil
}
