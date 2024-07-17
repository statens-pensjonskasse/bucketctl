package printer

import (
	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"git.spk.no/infra/bucketctl/pkg/common"
	"sort"
	"strconv"
)

func PrettyFormatProjectWebhooks(projectConfig *ProjectConfig) [][]string {
	projectKey := projectConfig.Spec.ProjectKey
	projectWebhooks := projectConfig.Spec.Webhooks
	projectRepositories := projectConfig.Spec.Repositories

	var data [][]string
	data = append(data, []string{"Project", "Repository", "ID", "Type", "Events", "URL", "Active", "Verify SSL"})

	formattedProjectWebhooks := prettyFormatWebhooks(projectKey, "PROJECT", projectWebhooks)
	data = append(data, formattedProjectWebhooks...)

	formattedRepoWebhooks := prettyFormatRepositoriesWebhooks(projectKey, projectRepositories)
	data = append(data, formattedRepoWebhooks...)

	return data
}

func prettyFormatRepositoriesWebhooks(projectKey string, repositoriesWebhooks *RepositoriesProperties) [][]string {
	var data [][]string

	if repositoriesWebhooks == nil {
		return data
	}

	repoWebhooksMap := make(map[string]*Webhooks, len(*repositoriesWebhooks))
	for _, r := range *repositoriesWebhooks {
		repoWebhooksMap[r.RepoSlug] = r.Webhooks
	}

	for _, slug := range common.GetLexicallySortedKeys(repoWebhooksMap) {
		repoWebhooks := prettyFormatWebhooks(projectKey, slug, repoWebhooksMap[slug])
		data = append(data, repoWebhooks...)
	}

	return data
}

func prettyFormatWebhooks(projectKey string, repoSlug string, webhooks *Webhooks) [][]string {
	var data [][]string

	if webhooks == nil {
		return data
	}

	webhooksMap := make(map[string]*Webhook)
	for _, w := range *webhooks {
		webhooksMap[w.Name] = w
	}

	for _, webhook := range common.GetLexicallySortedKeys(webhooksMap) {
		sort.Strings(webhooksMap[webhook].Events)
		data = append(data, []string{projectKey, repoSlug, strconv.Itoa(webhooksMap[webhook].Id), webhooksMap[webhook].Name, webhooksMap[webhook].Events[0], webhooksMap[webhook].Url, strconv.FormatBool(webhooksMap[webhook].Active), strconv.FormatBool(webhooksMap[webhook].SslVerificationRequired)})
		for _, event := range webhooksMap[webhook].Events[1:] {
			data = append(data, []string{"", "", "", "", event, "", "", ""})
		}
	}

	return data
}
