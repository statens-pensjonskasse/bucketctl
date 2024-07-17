package get

import (
	"git.spk.no/infra/bucketctl/pkg/api/bitbucket"
	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"git.spk.no/infra/bucketctl/pkg/common"
	"git.spk.no/infra/bucketctl/pkg/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listWebhooksCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString(common.ProjectKeyFlag) == "" {
			cmd.MarkFlagRequired(common.ProjectKeyFlag)
		}
		viper.BindPFlag(common.ProjectKeyFlag, cmd.Flags().Lookup(common.ProjectKeyFlag))
		viper.BindPFlag(common.RepoSlugFlag, cmd.Flags().Lookup(common.RepoSlugFlag))
	},
	Use:     "webhooks",
	Short:   "List webhooks for given permission or repo",
	Aliases: []string{"wh"},
	Run:     listWebhooks,
}

func listWebhooks(cmd *cobra.Command, args []string) {
	baseUrl := viper.GetString(common.BaseUrlFlag)
	projectKey := viper.GetString(common.ProjectKeyFlag)
	repoSlug := viper.GetString(common.RepoSlugFlag)
	limit := viper.GetInt(common.LimitFlag)
	token := viper.GetString(common.TokenFlag)

	projectConfig := ProjectConfigV1alpha1()
	projectConfig.Metadata.Name = projectKey

	if repoSlug == "" {
		webhooks, err := FetchWebhooks(baseUrl, projectKey, limit, token)
		cobra.CheckErr(err)
		projectConfig.Spec = *webhooks
	} else {
		repoWebhooks, err := bitbucket.GetRepositoryWebhooks(baseUrl, projectKey, repoSlug, limit, token)
		cobra.CheckErr(err)
		projectConfig.Spec.ProjectKey = projectKey
		projectConfig.Spec.Repositories = &RepositoriesProperties{&RepositoryProperties{RepoSlug: repoSlug, Webhooks: repoWebhooks}}
	}

	err := printer.PrintData(projectConfig, printer.PrettyFormatProjectWebhooks)
	cobra.CheckErr(err)
}

func FetchWebhooks(baseUrl string, projectKey string, limit int, token string) (*ProjectConfigSpec, error) {
	projectWebhooks, err := bitbucket.GetProjectWebhooks(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}
	repositoriesWebhooks, err := bitbucket.GetProjectRepositoriesWebhooks(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	return &ProjectConfigSpec{ProjectKey: projectKey, Webhooks: projectWebhooks, Repositories: repositoriesWebhooks}, nil
}
