package permission

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
	"gobit/pkg/cmd/project"
	"os"
)

var listAllPermissionsCmd = &cobra.Command{
	Use:  "all",
	RunE: listAllPermissions,
}

func getAllPermissions(baseUrl string, limit int, token string) (map[string]ProjectPermissions, error) {
	projects, err := project.GetProjects(baseUrl, limit)
	if err != nil {
		return nil, err
	}

	allPermissions := make(map[string]ProjectPermissions)
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(len(projects)).WithRemoveWhenDone(true).WithWriter(os.Stderr).Start()
	for _, proj := range projects {
		progressBar.Title = proj.Key
		projectPermissions, err := GetProjectPermissions(baseUrl, proj.Key, limit, token, true)
		if err != nil {
			return nil, err
		}
		allPermissions[proj.Key] = projectPermissions
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
