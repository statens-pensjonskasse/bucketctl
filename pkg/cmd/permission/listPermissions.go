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

var ListPermissionsCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("key", cmd.Flags().Lookup("key"))
		viper.BindPFlag("repo", cmd.Flags().Lookup("repo"))
	},
	Use:     "list",
	Aliases: []string{"l"},
	RunE:    listProjectPermissions,
}

func init() {
	ListPermissionsCmd.Flags().StringVarP(&key, "key", "k", "", "Project key")
	ListPermissionsCmd.Flags().StringVarP(&repo, "repo", "r", "", "Repository slug")

	ListPermissionsCmd.MarkFlagRequired("key")
}

func listProjectPermissions(cmd *cobra.Command, args []string) error {
	var baseUrl = viper.GetString("baseUrl")
	var projectKey = viper.GetString("key")
	var repoSlug = viper.GetString("repo")
	var limit = viper.GetInt("limit")
	var token = viper.GetString("token")

	if repoSlug == "" {
		permissionSet, err := GetProjectPermissions(baseUrl, projectKey, limit, token)
		if err != nil {
			return err
		}
		pkg.PrintData(permissionSet, PrettyFormatProjectPermissions)
	} else {
		permissionSet, err := getRepositoryPermissions(baseUrl, projectKey, repoSlug, limit, token)
		if err != nil {
			return err
		}
		pkg.PrintData(permissionSet, PrettyFormatRepositoryPermissions)
	}

	return nil
}
