package webhook

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
	"gobit/pkg/cmd/project"
	"os"
)

var listAllWebhooksCmd = &cobra.Command{
	Use:   "all",
	Short: "List all webhooks",
	RunE:  listAllWebhooks,
}

func listAllWebhooks(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString("baseUrl")
	limit := viper.GetInt("limit")
	token := viper.GetString("token")

	webhooks, err := getAllWebhooks(baseUrl, limit, token)
	if err != nil {
		return err
	}

	pkg.PrintData(webhooks, PrettyFormatProjectWebhooks)
	return nil
}

func getAllWebhooks(baseUrl string, limit int, token string) (map[string]*ProjectWebhooks, error) {
	projects, err := project.GetProjects(baseUrl, limit)
	if err != nil {
		return nil, err
	}

	allWebhooks := make(map[string]*ProjectWebhooks)
	progessBar, _ := pterm.DefaultProgressbar.WithTotal(len(projects)).WithRemoveWhenDone(true).WithWriter(os.Stderr).Start()
	for _, proj := range projects {
		progessBar.Title = proj.Key
		projectWebhooks, err := getProjectWebhooks(baseUrl, proj.Key, limit, token, true)
		if err != nil {
			return nil, err
		}
		allWebhooks[proj.Key] = projectWebhooks
		progessBar.Increment()
	}

	return allWebhooks, nil
}
