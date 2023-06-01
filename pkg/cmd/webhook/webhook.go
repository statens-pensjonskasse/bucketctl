package webhook

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
	"gobit/pkg/types"
	"sort"
	"strconv"
)

type ProjectWebhooks struct {
	Webhooks     []types.Webhook               `yaml:"webhooks,omitempty"`
	Repositories map[string]*RepositoryWebhook `yaml:"respositories,omitempty"`
}

type RepositoryWebhook struct {
	Webhooks []types.Webhook `yaml:"webhooks"`
}

var (
	key  string
	repo string
)

var Cmd = &cobra.Command{
	Use:     "webhook",
	Short:   "Bitbucket webhook commands",
	Aliases: []string{"wh"},
}

func init() {
	Cmd.PersistentFlags().StringVarP(&key, "key", "k", "", "Project key")
	Cmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "Repository slug. Leave empty to query project permissions.")

	Cmd.MarkPersistentFlagRequired("key")

	Cmd.AddCommand(listWebhooksCmd)
}

var listWebhooksCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		Cmd.MarkPersistentFlagRequired("repo")
		viper.BindPFlag("key", Cmd.PersistentFlags().Lookup("key"))
		viper.BindPFlag("repo", Cmd.PersistentFlags().Lookup("repo"))
	},
	Use:   "list",
	Short: "List Webhooks for given repo",
	RunE:  listWebhooks,
}

func getRepositoryWebhooks(baseUrl string, projectKey string, repoSlug string, limit int, token string) ([]types.Webhook, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/webhooks?limit=%d", baseUrl, projectKey, repoSlug, limit)

	body, err := pkg.GetRequestBody(url, token)
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

func PrettyFormatProjectWebhooks(projectWebhooks map[string]*ProjectWebhooks) [][]string {
	var data [][]string
	data = append(data, []string{"Project", "Repository", "ID", "Name", "Events", "URL", "Active", "VerifySSL"})

	for projectKey, webhooks := range projectWebhooks {
		for repoSlug, repoWebhooks := range webhooks.Repositories {
			for _, webhook := range repoWebhooks.Webhooks {
				sort.Strings(webhook.Events)
				data = append(data, []string{projectKey, repoSlug, strconv.Itoa(webhook.Id), webhook.Name, webhook.Events[0], webhook.Url, strconv.FormatBool(webhook.Active), strconv.FormatBool(webhook.SslVerificationRequired)})
				for _, event := range webhook.Events[1:] {
					data = append(data, []string{"", "", "", "", event, "", "", ""})
				}
			}
		}
	}

	return data
}
