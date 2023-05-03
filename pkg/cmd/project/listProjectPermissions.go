package project

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
)

var (
	key string
)

var ListProjectPermissionsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	RunE:    listProjectPermissions,
}

func init() {
	ListProjectPermissionsCmd.Flags().StringVarP(&key, "key", "k", "", "Project key")
	ListProjectPermissionsCmd.MarkFlagRequired("key")
	viper.BindPFlag("key", ListProjectPermissionsCmd.Flags().Lookup("key"))
}

func listProjectPermissions(cmd *cobra.Command, args []string) error {
	var baseUrl = viper.GetString("baseUrl")
	var projectKey = viper.GetString("key")
	var limit = viper.GetInt("limit")
	var token = viper.GetString("token")

	permissionSet, err := GetProjectPermissions(baseUrl, projectKey, limit, token)
	if err != nil {
		return err
	}

	pkg.PrintData(permissionSet, PrettyFormatProjectPermissions)
	return nil
}
