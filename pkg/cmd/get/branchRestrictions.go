package get

import (
	"git.spk.no/infra/bucketctl/pkg/api/bitbucket"
	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"git.spk.no/infra/bucketctl/pkg/common"
	"git.spk.no/infra/bucketctl/pkg/printer"
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
	Use:     "branchRestrictions",
	Short:   "List settings for given project or repository",
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

	cobra.CheckErr(printer.PrintData(projectConfig, printer.PrettyFormatProjectsSettings))
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
