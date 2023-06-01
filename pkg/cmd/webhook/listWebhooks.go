package webhook

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
)

func listWebhooks(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString("baseUrl")
	projectKey := viper.GetString("key")
	repoSlug := viper.GetString("repo")
	limit := viper.GetInt("limit")
	token := viper.GetString("token")

	webhooks, err := getRepositoryWebhooks(baseUrl, projectKey, repoSlug, limit, token)
	if err != nil {
		return err
	}

	projectWebhooks := map[string]*ProjectWebhooks{
		projectKey: {
			Repositories: map[string]*RepositoryWebhook{
				repoSlug: {
					Webhooks: webhooks,
				},
			},
		},
	}

	pkg.PrintData(projectWebhooks, PrettyFormatProjectWebhooks)

	return nil
}
