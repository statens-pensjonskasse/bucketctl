package get

import (
	"git.spk.no/infra/bucketctl/pkg/api/bitbucket"
	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"git.spk.no/infra/bucketctl/pkg/common"
	"git.spk.no/infra/bucketctl/pkg/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listBranchingModelsCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString(common.ProjectKeyFlag) == "" {
			cmd.MarkFlagRequired(common.ProjectKeyFlag)
		}
		viper.BindPFlag(common.ProjectKeyFlag, cmd.Flags().Lookup(common.ProjectKeyFlag))
		viper.BindPFlag(common.RepoSlugFlag, cmd.Flags().Lookup(common.RepoSlugFlag))
	},
	Use:     "branchingModel",
	Short:   "List branch model for given project or repository",
	Aliases: []string{"bm", "branchModel"},
	Run:     listBranchingModels,
}

func listBranchingModels(cmd *cobra.Command, args []string) {
	baseUrl := viper.GetString(common.BaseUrlFlag)
	projectKey := viper.GetString(common.ProjectKeyFlag)
	repoSlug := viper.GetString(common.RepoSlugFlag)
	limit := viper.GetInt(common.LimitFlag)
	token := viper.GetString(common.TokenFlag)

	projectConfig := ProjectConfigV1alpha1()
	projectConfig.Metadata.Name = projectKey

	if repoSlug == "" {
		branchingModels, err := FetchBranchingModels(baseUrl, projectKey, limit, token)
		cobra.CheckErr(err)
		projectConfig.Spec = *branchingModels

	} else {
		repoBranchModel, err := bitbucket.GetRepositoryBranchingModel(baseUrl, projectKey, repoSlug, token)
		cobra.CheckErr(err)
		projectConfig.Spec.ProjectKey = projectKey
		projectConfig.Spec.Repositories = &RepositoriesProperties{&RepositoryProperties{RepoSlug: repoSlug, BranchingModel: repoBranchModel}}
	}

	cobra.CheckErr(printer.PrintData(projectConfig, nil))
}

func FetchBranchingModels(baseUrl string, projectKey string, limit int, token string) (*ProjectConfigSpec, error) {
	projectBranchModel, err := bitbucket.GetProjectBranchingModel(baseUrl, projectKey, token)
	if err != nil {
		return nil, err
	}
	repositoriesProperties, err := bitbucket.GetRepositoriesBranchingModel(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	return &ProjectConfigSpec{ProjectKey: projectKey, BranchingModel: projectBranchModel, Repositories: repositoriesProperties}, nil
}
