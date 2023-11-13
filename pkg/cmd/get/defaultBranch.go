package get

import (
	"bucketctl/pkg/api/bitbucket"
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/common"
	"bucketctl/pkg/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listDefaultBranchCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString(common.ProjectKeyFlag) == "" {
			cmd.MarkFlagRequired(common.ProjectKeyFlag)
		}
		viper.BindPFlag(common.ProjectKeyFlag, cmd.Flags().Lookup(common.ProjectKeyFlag))
		viper.BindPFlag(common.RepoSlugFlag, cmd.Flags().Lookup(common.RepoSlugFlag))
	},
	Use:     "defaultBranch",
	Aliases: []string{"defaultBranches"},
	Run:     listDefaultBranches,
}

func listDefaultBranches(cmd *cobra.Command, args []string) {
	baseUrl := viper.GetString(common.BaseUrlFlag)
	projectKey := viper.GetString(common.ProjectKeyFlag)
	repoSlug := viper.GetString(common.RepoSlugFlag)
	limit := viper.GetInt(common.LimitFlag)
	token := viper.GetString(common.TokenFlag)

	projectConfig := ProjectConfigV1alpha1()
	projectConfig.Metadata.Name = projectKey

	if repoSlug == "" {
		defaultBranches, err := FetchDefaultBranches(baseUrl, projectKey, limit, token)
		cobra.CheckErr(err)
		projectConfig.Spec = *defaultBranches
	} else {
		repoDefaultBranch, err := bitbucket.GetRepositoryDefaultBranch(baseUrl, projectKey, repoSlug, token)
		cobra.CheckErr(err)
		projectConfig.Spec.ProjectKey = projectKey
		projectConfig.Spec.Repositories = &RepositoriesProperties{&RepositoryProperties{RepoSlug: repoSlug, DefaultBranch: &repoDefaultBranch.DisplayId}}
	}

	cobra.CheckErr(printer.PrintData(projectConfig, nil))
}

func FetchDefaultBranches(baseUrl string, projectKey string, limit int, token string) (*ProjectConfigSpec, error) {
	repositoriesProperties, err := bitbucket.GetRepositoriesDefaultBranch(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}
	return &ProjectConfigSpec{ProjectKey: projectKey, Repositories: repositoriesProperties}, nil
}
