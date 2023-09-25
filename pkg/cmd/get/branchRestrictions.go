package get

import (
	"bucketctl/pkg/api/bitbucket"
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/common"
	"bucketctl/pkg/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listBranchRestrictionsCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString(common.ProjectKeyFlag) == "" {
			cmd.MarkFlagRequired(common.ProjectKeyFlag)
		}
		viper.BindPFlag(common.ProjectKeyFlag, cmd.Flags().Lookup(common.ProjectKeyFlag))
		viper.BindPFlag(common.RepoSlugFlag, cmd.Flags().Lookup(common.RepoSlugFlag))
	},
	Use:     "branch-restrictions",
	Short:   "List settings for given permission or repository",
	Aliases: []string{"br", "restrictions"},
	Run:     listBranchRestrictions,
}

func listBranchRestrictions(cmd *cobra.Command, args []string) {
	baseUrl := viper.GetString(common.BaseUrlFlag)
	projectKey := viper.GetString(common.ProjectKeyFlag)
	repoSlug := viper.GetString(common.RepoSlugFlag)
	limit := viper.GetInt(common.LimitFlag)
	token := viper.GetString(common.TokenFlag)

	projectConfig := ProjectConfigV1alpha1()
	projectConfig.Metadata.Name = projectKey

	if repoSlug == "" {
		branchRestrictions, err := FetchBranchRestrictions(baseUrl, projectKey, limit, token)
		cobra.CheckErr(err)
		projectConfig.Spec = *branchRestrictions
	} else {
		repoBranchRestrictions, err := bitbucket.GetRepositoryBranchRestrictions(baseUrl, projectKey, repoSlug, limit, token)
		cobra.CheckErr(err)
		projectConfig.Spec.ProjectKey = projectKey
		projectConfig.Spec.Repositories = &RepositoriesProperties{&RepositoryProperties{RepoSlug: repoSlug, BranchRestrictions: repoBranchRestrictions}}
	}

	err := printer.PrintData(projectConfig, printer.PrettyFormatProjectsSettings)
	cobra.CheckErr(err)
}

func FetchBranchRestrictions(baseUrl string, projectKey string, limit int, token string) (*ProjectConfigSpec, error) {
	projectBranchRestrictions, err := bitbucket.GetProjectBranchRestrictions(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}
	repositoriesProperties, err := bitbucket.GetProjectRepositoriesBranchRestrictions(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	return &ProjectConfigSpec{ProjectKey: projectKey, BranchRestrictions: projectBranchRestrictions, Repositories: repositoriesProperties}, nil
}
