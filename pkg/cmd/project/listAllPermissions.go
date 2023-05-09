package project

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
	"os"
)

var listAllPermissionsCmd = &cobra.Command{
	Use:  "all",
	RunE: listAllPermissions,
}

func getAllPermissions(baseUrl string, limit int, token string) (*GrantedProjectPermissions, error) {
	projects, err := GetProjects(baseUrl, limit)
	if err != nil {
		return nil, err
	}

	allPermissions := &GrantedProjectPermissions{
		Project: map[string]*PermissionSet{},
	}
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(len(projects)).WithRemoveWhenDone(true).WithWriter(os.Stderr).Start()
	for _, proj := range projects {
		progressBar.Title = proj.Key
		projectPermissions, err := GetProjectPermissions(baseUrl, proj.Key, limit, token)
		if err != nil {
			return nil, err
		}
		allPermissions.Project[proj.Key] = projectPermissions.Project[proj.Key]
		progressBar.Increment()
	}

	return allPermissions, nil
}

func listAllPermissions(cmd *cobra.Command, args []string) error {
	var baseUrl = viper.GetString("baseUrl")
	var limit = viper.GetInt("limit")
	var token = viper.GetString("token")

	permissions, err := getAllPermissions(baseUrl, limit, token)
	if err != nil {
		return err
	}

	pkg.PrintData(permissions, PrettyFormatProjectPermissions)
	return nil
}
