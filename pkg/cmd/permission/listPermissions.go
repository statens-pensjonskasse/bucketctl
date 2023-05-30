package permission

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
)

var (
	key  string
	repo string
)

var listPermissionsCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("key", cmd.Flags().Lookup("key"))
		viper.BindPFlag("repo", cmd.Flags().Lookup("repo"))
		viper.BindPFlag("include-repos", cmd.Flags().Lookup("include-repos"))
	},
	Use:     "list",
	Aliases: []string{"l"},
	RunE:    listPermissions,
}

func init() {
	listPermissionsCmd.Flags().StringVarP(&key, "key", "k", "", "Project key")
	listPermissionsCmd.Flags().StringVarP(&repo, "repo", "r", "", "Repository slug. Leave empty to query project permissions.")
	listPermissionsCmd.Flags().Bool("include-repos", false, "Include repository permissions when querying project permissions")

	listPermissionsCmd.MarkFlagRequired("key")
}

func listPermissions(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString("baseUrl")
	projectKey := viper.GetString("key")
	repoSlug := viper.GetString("repo")
	limit := viper.GetInt("limit")
	token := viper.GetString("token")
	includeRepos := viper.GetBool("include-repos")

	if repoSlug == "" {
		projectPermissions, err := GetProjectPermissions(baseUrl, projectKey, limit, token, includeRepos)
		if err != nil {
			return err
		}

		projectPermissionsMap := make(map[string]ProjectPermissions)
		projectPermissionsMap[projectKey] = projectPermissions

		// TODO: Also print repo permissions
		pkg.PrintData(projectPermissionsMap, PrettyFormatProjectPermissions)
	} else {
		repositoryPermissions, err := getRepositoryPermissions(baseUrl, projectKey, repoSlug, limit, token)
		if err != nil {
			return err
		}

		repoPermissionsMap := make(map[string]RepositoryPermissions)
		repoPermissionsMap[repoSlug] = repositoryPermissions

		pkg.PrintData(repoPermissionsMap, PrettyFormatRepositoryPermissions)
	}

	return nil
}
