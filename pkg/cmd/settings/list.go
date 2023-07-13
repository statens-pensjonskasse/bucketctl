package settings

import (
	"bucketctl/pkg"
	"bucketctl/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listSettingsCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString(types.ProjectKeyFlag) == "" {
			cmd.MarkFlagRequired(types.ProjectKeyFlag)
		}
		viper.BindPFlag(types.ProjectKeyFlag, cmd.Flags().Lookup(types.ProjectKeyFlag))
		viper.BindPFlag(types.RepoSlugFlag, cmd.Flags().Lookup(types.RepoSlugFlag))
		viper.BindPFlag(types.IncludeReposFlag, cmd.Flags().Lookup(types.IncludeReposFlag))
	},
	Use:   "list",
	Short: "List settings for given project or repository",
	RunE:  listSettings,
}

func init() {
	listSettingsCmd.Flags().StringVarP(&key, types.ProjectKeyFlag, "k", "", "Project key")
	listSettingsCmd.Flags().StringVarP(&repo, types.RepoSlugFlag, "r", "", "Repository slug. Leave empty to query project webhooks.")
	listSettingsCmd.Flags().Bool(types.IncludeReposFlag, false, "Include repository permissions when querying project")
}

func listSettings(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString(types.BaseUrlFlag)
	projectKey := viper.GetString(types.ProjectKeyFlag)
	repoSlug := viper.GetString(types.RepoSlugFlag)
	limit := viper.GetInt(types.LimitFlag)
	token := viper.GetString(types.TokenFlag)
	includeRepos := viper.GetBool(types.IncludeReposFlag)

	projectSettingsMap := make(map[string]*ProjectSettings)
	if repoSlug == "" {
		projectRestrictions, err := getProjectRestrictions(baseUrl, projectKey, limit, token, includeRepos)
		if err != nil {
			return err
		}
		projectSettingsMap[projectKey] = projectRestrictions
	} else {
		repoRestrictions, err := getRepositoryRestrictions(baseUrl, projectKey, repoSlug, limit, token)
		if err != nil {
			return err
		}
		projectSettingsMap[projectKey] = new(ProjectSettings)
		projectSettingsMap[projectKey].Repositories = make(map[string]*RepositorySettings)
		projectSettingsMap[projectKey].Repositories[repoSlug] = new(RepositorySettings)
		projectSettingsMap[projectKey].Repositories[repoSlug].Restrictions = repoRestrictions.Restrictions
	}

	return pkg.PrintData(projectSettingsMap, prettyFormatProjectsSettings)
}
