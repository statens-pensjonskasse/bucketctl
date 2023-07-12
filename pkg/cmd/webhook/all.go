package webhook

import (
	"bucketctl/pkg"
	"bucketctl/pkg/cmd/project"
	"bucketctl/pkg/types"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var listAllWebhooksCmd = &cobra.Command{
	Use:   "all",
	Short: "List all webhooks",
	RunE:  listAllWebhooks,
}

func listAllWebhooks(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString(types.BaseUrlFlag)
	limit := viper.GetInt(types.LimitFlag)
	token := viper.GetString(types.TokenFlag)

	webhooks, err := getAllWebhooks(baseUrl, limit, token)
	if err != nil {
		return err
	}

	return pkg.PrintData(webhooks, prettyFormatProjectWebhooks)
}

func getAllWebhooks(baseUrl string, limit int, token string) (map[string]*ProjectWebhooks, error) {
	projects, err := project.GetProjects(baseUrl, token, limit)
	if err != nil {
		return nil, err
	}

	allWebhooks := make(map[string]*ProjectWebhooks)
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(len(projects)).WithRemoveWhenDone(true).WithWriter(os.Stderr).Start()
	for projectKey := range projects {
		progressBar.Title = projectKey
		projectWebhooks, err := getProjectWebhooks(baseUrl, projectKey, limit, token, true)
		if err != nil {
			return nil, err
		}
		allWebhooks[projectKey] = projectWebhooks
		progressBar.Increment()
	}

	return allWebhooks, nil
}
