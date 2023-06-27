package webhook

import (
	"bucketctl/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listWebhooksCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		Cmd.MarkPersistentFlagRequired("repo")
		viper.BindPFlag("key", cmd.Flags().Lookup("key"))
		viper.BindPFlag("repo", cmd.Flags().Lookup("repo"))
		viper.BindPFlag("include-repos", cmd.Flags().Lookup("include-repos"))
	},
	Use:   "list",
	Short: "List Webhooks for given repo",
	RunE:  listWebhooks,
}

func init() {
	listWebhooksCmd.Flags().StringVarP(&key, "key", "k", "", "Project key")
	listWebhooksCmd.Flags().StringVarP(&repo, "repo", "r", "", "Repository slug. Leave empty to query project webhooks.")
	listWebhooksCmd.Flags().Bool("include-repos", false, "Include repository permissions when querying project")

	listWebhooksCmd.MarkFlagRequired("key")
}

func listWebhooks(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString("baseUrl")
	projectKey := viper.GetString("key")
	repoSlug := viper.GetString("repo")
	limit := viper.GetInt("limit")
	token := viper.GetString("token")
	includeRepos := viper.GetBool("include-repos")

	projectWebhooksMap := make(map[string]*ProjectWebhooks)
	if repoSlug == "" {
		projectWebhooks, err := getProjectWebhooks(baseUrl, projectKey, limit, token, includeRepos)
		if err != nil {
			return err
		}
		projectWebhooksMap[projectKey] = projectWebhooks
	} else {
		webhooks, err := getRepositoryWebhooks(baseUrl, projectKey, repoSlug, limit, token)
		if err != nil {
			return err
		}
		projectWebhooksMap[projectKey] = new(ProjectWebhooks)
		projectWebhooksMap[projectKey].Repositories = make(map[string]*RepositoryWebhooks)
		projectWebhooksMap[projectKey].Repositories[repoSlug] = webhooks
	}

	return pkg.PrintData(projectWebhooksMap, PrettyFormatProjectWebhooks)
}
