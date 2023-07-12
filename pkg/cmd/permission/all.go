package permission

import (
	"bucketctl/pkg"
	"bucketctl/pkg/cmd/project"
	"bucketctl/pkg/types"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var listAllPermissionsCmd = &cobra.Command{
	Use:   "all",
	Short: "List all permissions",
	RunE:  listAllPermissions,
}

func listAllPermissions(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString(types.BaseUrlFlag)
	limit := viper.GetInt(types.LimitFlag)
	token := viper.GetString(types.TokenFlag)

	permissions, err := getAllPermissions(baseUrl, limit, token)
	if err != nil {
		return err
	}

	return pkg.PrintData(permissions, prettyFormatProjectPermissions)
}

func getAllPermissions(baseUrl string, limit int, token string) (map[string]*ProjectPermissions, error) {
	projects, err := project.GetProjects(baseUrl, token, limit)
	if err != nil {
		return nil, err
	}

	allPermissions := make(map[string]*ProjectPermissions)
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(len(projects)).WithRemoveWhenDone(true).WithWriter(os.Stderr).Start()
	for projectKey := range projects {
		progressBar.Title = projectKey
		projectPermissions, err := getProjectPermissions(baseUrl, projectKey, limit, token, true)
		if err != nil {
			return nil, err
		}
		allPermissions[projectKey] = projectPermissions
		progressBar.Increment()
	}

	return allPermissions, nil
}
