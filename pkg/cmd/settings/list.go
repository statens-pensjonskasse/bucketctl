package settings

import (
	"bucketctl/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listSettingsCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		Cmd.MarkPersistentFlagRequired("repo")
		viper.BindPFlag("key", cmd.Flags().Lookup("key"))
		viper.BindPFlag("repo", cmd.Flags().Lookup("repo"))
		viper.BindPFlag("include-repos", cmd.Flags().Lookup("include-repos"))
	},
	Use:   "list",
	Short: "List settings for given project/repo",
	RunE:  listSettings,
}

func init() {
	listSettingsCmd.Flags().StringVarP(&key, "key", "k", "", "Project key")
	listSettingsCmd.Flags().StringVarP(&repo, "repo", "r", "", "Repository slug. Leave empty to query project webhooks.")
	listSettingsCmd.Flags().Bool("include-repos", false, "Include repository permissions when querying project")

	listSettingsCmd.MarkFlagRequired("key")
}

func listSettings(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString("baseUrl")
	projectKey := viper.GetString("key")
	repoSlug := viper.GetString("repo")
	limit := viper.GetInt("limit")
	token := viper.GetString("token")
	includeRepos := viper.GetBool("include-repos")

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

	return pkg.PrintData(projectSettingsMap, nil)
}
