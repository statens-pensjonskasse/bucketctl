package permission

import (
	"bucketctl/pkg/common"
	"bucketctl/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	key  string
	repo string
)

var listPermissionsCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString(types.ProjectKeyFlag) == "" {
			cmd.MarkFlagRequired(types.ProjectKeyFlag)
		}
		viper.BindPFlag(types.ProjectKeyFlag, cmd.Flags().Lookup(types.ProjectKeyFlag))
		viper.BindPFlag(types.RepoSlugFlag, cmd.Flags().Lookup(types.RepoSlugFlag))
		viper.BindPFlag(types.IncludeReposFlag, cmd.Flags().Lookup(types.IncludeReposFlag))
	},
	Use:     "list",
	Short:   "List permissions for given project or repository",
	Aliases: []string{"l"},
	RunE:    listPermissions,
}

func init() {
	listPermissionsCmd.Flags().StringVarP(&key, types.ProjectKeyFlag, types.ProjectKeyFlagShorthand, "", "Project key")
	listPermissionsCmd.Flags().StringVarP(&repo, types.RepoSlugFlag, types.RepoSlugFlagShorthand, "", "Repository slug. Leave empty to query project permissions.")
	listPermissionsCmd.Flags().Bool(types.IncludeReposFlag, false, "Include repository permissions when querying project permissions")
}

func listPermissions(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString(types.BaseUrlFlag)
	projectKey := viper.GetString(types.ProjectKeyFlag)
	repoSlug := viper.GetString(types.RepoSlugFlag)
	limit := viper.GetInt(types.LimitFlag)
	token := viper.GetString(types.TokenFlag)
	includeRepos := viper.GetBool(types.IncludeReposFlag)

	projectPermissionsMap := make(map[string]*ProjectPermissions)
	if repoSlug == "" {
		projectPermissions, err := getProjectPermissions(baseUrl, projectKey, limit, token, includeRepos)
		if err != nil {
			return err
		}

		projectPermissionsMap[projectKey] = projectPermissions
	} else {
		permissions, err := getRepositoryPermissions(baseUrl, projectKey, repoSlug, limit, token)
		if err != nil {
			return err
		}

		projectPermissionsMap[projectKey] = new(ProjectPermissions)
		projectPermissionsMap[projectKey].Repositories = make(map[string]*RepositoryPermissions)
		projectPermissionsMap[projectKey].Repositories[repoSlug] = permissions
	}

	return common.PrintData(projectPermissionsMap, prettyFormatProjectPermissions)
}
