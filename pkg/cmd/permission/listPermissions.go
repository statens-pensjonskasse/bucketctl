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
		grantedProjectPermissions, err := GetProjectPermissions(baseUrl, projectKey, limit, token, includeRepos)
		if err != nil {
			return err
		}

		projectPermissions := &GrantedProjectPermissions{
			Project: map[string]*ProjectPermissionSet{
				projectKey: grantedProjectPermissions,
			},
		}

		pkg.PrintData(projectPermissions, PrettyFormatProjectPermissions)
	} else {
		grantedRepoPermissions, err := getRepositoryPermissions(baseUrl, projectKey, repoSlug, limit, token)
		if err != nil {
			return err
		}

		repositoryPermissions := &GrantedRepositoryPermissions{
			Repository: map[string]*PermissionSet{
				repoSlug: grantedRepoPermissions,
			},
		}

		pkg.PrintData(repositoryPermissions, PrettyFormatRepositoryPermissions)
	}

	return nil
}
