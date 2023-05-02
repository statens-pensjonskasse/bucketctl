package project

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
	"strings"
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

func prettyFormatProjectPermissions(pSet *PermissionSet) [][]string {
	var data [][]string

	data = append(data, []string{"Permission", "Groups", "Users"})

	for permission, v := range pSet.Permissions {
		var users string
		for _, user := range v.Users {
			users += user + "\n"
		}
		users = strings.Trim(users, "\n")
		var groups string
		for _, group := range v.Groups {
			groups += group + "\n"
		}
		groups = strings.Trim(groups, "\n")

		// Dersom verken en gruppe eller en bruker har rettigheten så hopper vi over den
		if len(groups)+len(users) > 0 {
			data = append(data, []string{permission, groups, users})
		}
	}

	return data
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

	pkg.PrintData(permissionSet, prettyFormatProjectPermissions)
	return nil
}
