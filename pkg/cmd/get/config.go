package get

import (
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/common"
	"bucketctl/pkg/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getProjectConfigCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString(common.ProjectKeyFlag) == "" {
			cmd.MarkFlagRequired(common.ProjectKeyFlag)
		}
		viper.BindPFlag(common.ProjectKeyFlag, cmd.Flags().Lookup(common.ProjectKeyFlag))
		viper.BindPFlag(common.RepoSlugFlag, cmd.Flags().Lookup(common.RepoSlugFlag))
	},
	Use:     "project-config",
	Short:   "Get complete project config",
	Aliases: []string{"project", "config", "pc"},
	Run:     getProjectConfig,
}

func getProjectConfig(cmd *cobra.Command, args []string) {
	baseUrl := viper.GetString(common.BaseUrlFlag)
	projectKey := viper.GetString(common.ProjectKeyFlag)
	limit := viper.GetInt(common.LimitFlag)
	token := viper.GetString(common.TokenFlag)

	projectConfig := ProjectConfigV1alpha1()
	projectConfig.Metadata.Name = projectKey

	projectConfigSpec, err := FetchProjectConfigSpec(baseUrl, projectKey, limit, token)
	cobra.CheckErr(err)
	projectConfig.Spec = *projectConfigSpec

	err = printer.PrintData(projectConfig, nil)
	cobra.CheckErr(err)
}

func FetchProjectConfigSpec(baseUrl string, projectKey string, limit int, token string) (*ProjectConfigSpec, error) {
	actualAccess, err := FetchPermissions(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}
	actualBranchingModels, err := FetchBranchingModels(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}
	actualBranchRestrictions, err := FetchBranchRestrictions(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}
	actualDefaultBranches, err := FetchDefaultBranches(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}
	actualWebhooks, err := FetchWebhooks(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	projectConfigSpec := CombineProjectConfigSpecs(&UncombinedProjectConfigSpecs{
		Access:             actualAccess,
		DefaultBranches:    actualDefaultBranches,
		BranchingModels:    actualBranchingModels,
		BranchRestrictions: actualBranchRestrictions,
		Webhooks:           actualWebhooks,
	})
	projectConfigSpec.ProjectKey = projectKey
	return projectConfigSpec, nil
}
