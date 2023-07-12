package settings

import (
	"bucketctl/pkg"
	"bucketctl/pkg/cmd/project"
	"bucketctl/pkg/types"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var listAllSettingsCmd = &cobra.Command{
	Use:   "all",
	Short: "Get all settings",
	RunE:  listAllSettings,
}

func listAllSettings(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString(types.BaseUrlFlag)
	limit := viper.GetInt(types.LimitFlag)
	token := viper.GetString(types.TokenFlag)

	settings, err := getAllSettings(baseUrl, limit, token)
	if err != nil {
		return err
	}

	return pkg.PrintData(settings, nil)
}

func getAllSettings(baseUrl string, limit int, token string) (map[string]*ProjectSettings, error) {
	projects, err := project.GetProjects(baseUrl, limit)
	if err != nil {
		return nil, err
	}

	allSettings := make(map[string]*ProjectSettings)
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(len(projects)).WithRemoveWhenDone(true).WithWriter(os.Stderr).Start()
	for projectKey := range projects {
		progressBar.Title = projectKey
		projectSettings, err := getProjectRestrictions(baseUrl, projectKey, limit, token, true)
		if err != nil {
			return nil, err
		}
		allSettings[projectKey] = projectSettings
		progressBar.Increment()
	}

	return allSettings, nil
}
