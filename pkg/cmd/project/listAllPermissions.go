package project

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var listAllPermissionsCmd = &cobra.Command{
	Use:  "all",
	RunE: listAllPermissions,
}

func getAllPermissions(baseUrl string, limit int, token string) (map[string]*PermissionSet, error) {
	projects, err := GetProjects(baseUrl, limit)
	if err != nil {
		return nil, err
	}

	projectPermissions := make(map[string]*PermissionSet)
	for _, proj := range projects.Values {
		projectPermission, err := GetProjectPermissions(baseUrl, proj.Key, limit, token)
		if err != nil {
			return nil, err
		}
		projectPermissions[proj.Key] = new(PermissionSet)
		projectPermissions[proj.Key].Permissions = projectPermission.Permissions
	}

	return projectPermissions, nil
}

func listAllPermissions(cmd *cobra.Command, args []string) error {
	var baseUrl = viper.GetString("baseUrl")
	var limit = viper.GetInt("limit")
	var token = viper.GetString("token")

	permissions, err := getAllPermissions(baseUrl, limit, token)
	if err != nil {
		return err
	}

	yamlData, err := yaml.Marshal(permissions)
	if err != nil {
		pterm.Error.Println("Error while Marshaling to YAML. %v", err)
	}
	pterm.Println(string(yamlData))

	return nil
}
