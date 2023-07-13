package webhook

import (
	"bucketctl/pkg"
	"bucketctl/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listWebhooksCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString(types.ProjectKeyFlag) == "" {
			cmd.MarkFlagRequired(types.ProjectKeyFlag)
		}
		viper.BindPFlag(types.ProjectKeyFlag, cmd.Flags().Lookup(types.ProjectKeyFlag))
		viper.BindPFlag(types.RepoSlugFlag, cmd.Flags().Lookup(types.RepoSlugFlag))
		viper.BindPFlag(types.IncludeReposFlag, cmd.Flags().Lookup(types.IncludeReposFlag))
	},
	Use:   "list",
	Short: "List webhooks for given project or repo",
	RunE:  listWebhooks,
}

func init() {
	listWebhooksCmd.Flags().StringVarP(&key, types.ProjectKeyFlag, "k", "", "Project key")
	listWebhooksCmd.Flags().StringVarP(&repo, types.RepoSlugFlag, "r", "", "Repository slug. Leave empty to query project webhooks.")
	listWebhooksCmd.Flags().Bool(types.IncludeReposFlag, false, "Include repository permissions when querying project")
}

func listWebhooks(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString(types.BaseUrlFlag)
	projectKey := viper.GetString(types.ProjectKeyFlag)
	repoSlug := viper.GetString(types.RepoSlugFlag)
	limit := viper.GetInt(types.LimitFlag)
	token := viper.GetString(types.TokenFlag)
	includeRepos := viper.GetBool(types.IncludeReposFlag)

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

	return pkg.PrintData(projectWebhooksMap, prettyFormatProjectWebhooks)
}
