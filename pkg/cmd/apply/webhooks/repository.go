package webhooks

import (
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/common"
	"bucketctl/pkg/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
)

func findRepositoriesWebhookChanges(desired *RepositoriesProperties, actual *RepositoriesProperties) (toCreate *RepositoriesProperties, toUpdate *RepositoriesProperties, toDelete *RepositoriesProperties) {
	toCreate = new(RepositoriesProperties)
	toUpdate = new(RepositoriesProperties)
	toDelete = new(RepositoriesProperties)
	for repoSlug, repo := range GroupRepositories(desired, actual) {
		whToCreate, whToUpdate, whToDelete := FindWebhooksToChange(repo.Desired.Webhooks, repo.Actual.Webhooks)
		if whToCreate != nil && len(*whToCreate) > 0 {
			*toCreate = append(*toCreate, &RepositoryProperties{RepoSlug: repoSlug, Webhooks: whToCreate})
		}
		if whToUpdate != nil && len(*whToUpdate) > 0 {
			*toUpdate = append(*toUpdate, &RepositoryProperties{RepoSlug: repoSlug, Webhooks: whToUpdate})
		}
		if whToDelete != nil && len(*whToDelete) > 0 {
			*toDelete = append(*toDelete, &RepositoryProperties{RepoSlug: repoSlug, Webhooks: whToDelete})
		}
	}
	return toCreate, toUpdate, toDelete
}

func setRepositoriesWebhooks(baseUrl string, projectKey string, token string, toCreate *RepositoriesProperties, toUpdate *RepositoriesProperties, toDelete *RepositoriesProperties) error {
	for _, r := range *toDelete {
		if err := deleteRepositoryWebhooks(baseUrl, projectKey, r.RepoSlug, token, r.Webhooks); err != nil {
			return err
		}
	}
	for _, r := range *toUpdate {
		if err := updateRepositoryWebhooks(baseUrl, projectKey, r.RepoSlug, token, r.Webhooks); err != nil {
			return err
		}
	}
	for _, r := range *toCreate {
		if err := createRepositoryWebhooks(baseUrl, projectKey, r.RepoSlug, token, r.Webhooks); err != nil {
			return err
		}
	}

	return nil
}

func createRepositoryWebhooks(baseUrl string, projectKey string, repoSlug string, token string, webhooks *Webhooks) error {
	if webhooks != nil && len(*webhooks) > 0 {
		for _, w := range *webhooks {
			if err := createRepositoryWebhook(baseUrl, projectKey, repoSlug, token, w); err != nil {
				return err
			}
		}
	}
	return nil
}

func createRepositoryWebhook(baseUrl string, projectKey string, repoSlug string, token string, webhook *Webhook) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/webhooks", baseUrl, projectKey, repoSlug)
	payload, err := json.Marshal(ToBitbucketWebhook(webhook))
	if err != nil {
		return err
	}
	if _, err := common.PostRequest(url, token, bytes.NewReader(payload), nil); err != nil {
		return err
	}
	logger.Log("%s webhook %s in repository %s/%s", pterm.Green("🪝 Created"), webhook.Name, projectKey, repoSlug)
	return nil
}

func updateRepositoryWebhooks(baseUrl string, projectKey string, repoSlug string, token string, webhooks *Webhooks) error {
	if webhooks != nil && len(*webhooks) > 0 {
		for _, w := range *webhooks {
			if err := updateRepositoryWebhook(baseUrl, projectKey, repoSlug, token, w); err != nil {
				return err
			}
		}
	}
	return nil
}

func updateRepositoryWebhook(baseUrl string, projectKey string, repoSlug string, token string, webhook *Webhook) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/webhooks/%d", baseUrl, projectKey, repoSlug, webhook.Id)
	payload, err := json.Marshal(ToBitbucketWebhook(webhook))
	if err != nil {
		return err
	}
	if _, err := common.PutRequest(url, token, bytes.NewReader(payload), nil); err != nil {
		return err
	}
	logger.Log("%s webhook %s in repository %s/%s", pterm.Blue("♻️ Updated"), webhook.Name, projectKey, repoSlug)
	return nil
}

func deleteRepositoryWebhooks(baseUrl string, projectKey string, repoSlug string, token string, webhooks *Webhooks) error {
	if webhooks != nil && len(*webhooks) > 0 {
		for _, w := range *webhooks {
			if err := deleteRepositoryWebhook(baseUrl, projectKey, repoSlug, token, w); err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteRepositoryWebhook(baseUrl string, projectKey string, repoSlug string, token string, webhook *Webhook) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/webhooks/%d", baseUrl, projectKey, repoSlug, webhook.Id)
	if _, err := common.DeleteRequest(url, token, nil); err != nil {
		return err
	}
	logger.Log("%s webhook %s in repository %s/%s", pterm.Red("️🛑 Deleted"), webhook.Name, projectKey, repoSlug)
	return nil
}
