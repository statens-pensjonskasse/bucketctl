package get

import (
	"git.spk.no/infra/bucketctl/pkg/api/bitbucket"
	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"git.spk.no/infra/bucketctl/pkg/common"
	"git.spk.no/infra/bucketctl/pkg/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listAccessCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString(common.ProjectKeyFlag) == "" {
			cmd.MarkFlagRequired(common.ProjectKeyFlag)
		}
		viper.BindPFlag(common.ProjectKeyFlag, cmd.Flags().Lookup(common.ProjectKeyFlag))
		viper.BindPFlag(common.RepoSlugFlag, cmd.Flags().Lookup(common.RepoSlugFlag))
	},
	Use:     "access",
	Short:   "List access policies for a given project or repository",
	Aliases: []string{"acc", "permissions"},
	Run:     listAccess,
}

func listAccess(cmd *cobra.Command, args []string) {
	baseUrl := viper.GetString(common.BaseUrlFlag)
	projectKey := viper.GetString(common.ProjectKeyFlag)
	repoSlug := viper.GetString(common.RepoSlugFlag)
	limit := viper.GetInt(common.LimitFlag)
	token := viper.GetString(common.TokenFlag)

	projectConfig := ProjectConfigV1alpha1()
	projectConfig.Metadata.Name = projectKey

	if repoSlug == "" {
		access, err := FetchPermissions(baseUrl, projectKey, limit, token)
		cobra.CheckErr(err)
		projectConfig.Spec = *access
	} else {
		permissions, err := bitbucket.GetRepositoryPermissions(baseUrl, projectKey, repoSlug, limit, token)
		cobra.CheckErr(err)
		projectConfig.Spec.ProjectKey = projectKey
		projectConfig.Spec.Repositories = &RepositoriesProperties{&RepositoryProperties{RepoSlug: repoSlug, Permissions: permissions}}
	}

	err := printer.PrintData(projectConfig, printer.PrettyFormatAccess)
	cobra.CheckErr(err)
}

func FetchPermissions(baseUrl string, projectKey string, limit int, token string) (*ProjectConfigSpec, error) {
	actualProjectAccess, err := bitbucket.GetProjectAccess(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}
	actualRepoAccess, err := bitbucket.GetProjectRepositoriesPermissions(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	return &ProjectConfigSpec{
		ProjectKey:        projectKey,
		Public:            actualProjectAccess.Public,
		DefaultPermission: actualProjectAccess.DefaultPermission,
		Permissions:       actualProjectAccess.Permissions,
		Repositories:      actualRepoAccess}, nil
}
