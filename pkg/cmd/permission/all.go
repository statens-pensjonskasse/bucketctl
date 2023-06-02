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
	Use:   "all",
	Short: "List all permissions",
	RunE:  listAllPermissions,
}

func listAllPermissions(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString("baseUrl")
	limit := viper.GetInt("limit")
	token := viper.GetString("token")

	permissions, err := getAllPermissions(baseUrl, limit, token)
	if err != nil {
		return err
	}

	pkg.PrintData(permissions, PrettyFormatProjectPermissions)
	return nil
}

func getAllPermissions(baseUrl string, limit int, token string) (map[string]*ProjectPermissions, error) {
	projects, err := project.GetProjects(baseUrl, limit)
	if err != nil {
		return nil, err
	}

	allPermissions := make(map[string]*ProjectPermissions)
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(len(projects)).WithRemoveWhenDone(true).WithWriter(os.Stderr).Start()
	for _, proj := range projects {
		progressBar.Title = proj.Key
		projectPermissions, err := getProjectPermissions(baseUrl, proj.Key, limit, token, true)
		if err != nil {
			return nil, err
		}
		allPermissions[proj.Key] = projectPermissions
		progressBar.Increment()
	}

	return allPermissions, nil
}
