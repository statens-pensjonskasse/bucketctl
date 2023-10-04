package webhooks

import (
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/common"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
)

func setProjectWebhooks(baseUrl string, projectKey string, token string, toCreate *Webhooks, toUpdate *Webhooks, toDelete *Webhooks) error {
	if err := createProjectWebhooks(baseUrl, projectKey, token, toCreate); err != nil {
		return err
	}
	if err := updateProjectWebhooks(baseUrl, projectKey, token, toUpdate); err != nil {
		return err
	}
	if err := deleteProjectWebhooks(baseUrl, projectKey, token, toDelete); err != nil {
		return err
	}
	return nil
}

func createProjectWebhooks(baseUrl string, projectKey string, token string, webhooks *Webhooks) error {
	for _, w := range *webhooks {
		if err := createProjectWebhook(baseUrl, projectKey, token, w); err != nil {
			return err
		}
	}
	return nil
}

func createProjectWebhook(baseUrl string, projectKey string, token string, webhook *Webhook) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/webhooks", baseUrl, projectKey)
	payload, err := json.Marshal(ToBitbucketWebhook(webhook))
	if err != nil {
		return err
	}
	if _, err := common.PostRequest(url, token, bytes.NewReader(payload), nil); err != nil {
		return err
	}
	pterm.Printfln("%s webhook %s in project %s", pterm.Green("🪝 Created"), webhook.Name, projectKey)
	return nil
}

func updateProjectWebhooks(baseUrl string, projectKey string, token string, webhooks *Webhooks) error {
	for _, w := range *webhooks {
		if err := updateProjectWebhook(baseUrl, projectKey, token, w); err != nil {
			return err
		}
	}
	return nil
}

func updateProjectWebhook(baseUrl string, projectKey string, token string, webhook *Webhook) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/webhooks/%d", baseUrl, projectKey, webhook.Id)
	payload, err := json.Marshal(ToBitbucketWebhook(webhook))
	if err != nil {
		return err
	}
	if _, err := common.PutRequest(url, token, bytes.NewReader(payload), nil); err != nil {
		return err
	}
	pterm.Printfln("%s webhook %s in project %s", pterm.Blue("♻️ Updated"), webhook.Name, projectKey)
	return nil
}

func deleteProjectWebhooks(baseUrl string, projectKey string, token string, webhooks *Webhooks) error {
	for _, w := range *webhooks {
		if err := deleteProjectWebhook(baseUrl, projectKey, token, w); err != nil {
			return err
		}
	}
	return nil
}

func deleteProjectWebhook(baseUrl string, projectKey string, token string, webhook *Webhook) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/webhooks/%d", baseUrl, projectKey, webhook.Id)
	if _, err := common.DeleteRequest(url, token, nil); err != nil {
		return err
	}
	pterm.Printfln("%s webhook %s in project %s", pterm.Red("🛑 Deleted"), webhook.Name, projectKey)
	return nil
}
