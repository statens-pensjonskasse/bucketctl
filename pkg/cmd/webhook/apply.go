package webhook

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
)

var (
	fileName string
)

var applyWebhooksCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("file", cmd.Flags().Lookup("file"))
		viper.BindPFlag("include-repos", cmd.Flags().Lookup("include-repos"))
	},
	Use:  "apply",
	RunE: applyWebhooks,
}

func init() {
	applyWebhooksCmd.Flags().StringVarP(&fileName, "file", "f", "", "Webhooks file")
	applyWebhooksCmd.Flags().Bool("include-repos", false, "Include repositories")

	applyWebhooksCmd.MarkFlagRequired("file")
}

func applyWebhooks(cmd *cobra.Command, args []string) error {
	file := viper.GetString("file")
	baseUrl := viper.GetString("baseUrl")
	limit := viper.GetInt("limit")
	token := viper.GetString("token")
	includeRepos := viper.GetBool("include-repos")

	var desiredWebhooks map[string]*ProjectWebhooks
	if err := pkg.ReadConfigFile(file, &desiredWebhooks); err != nil {
		return err
	}

	for projectKey, desiredProjectWebhooks := range desiredWebhooks {

		actualWebhooks, err := getProjectWebhooks(baseUrl, projectKey, limit, token, includeRepos)
		if err != nil {
			return err
		}

		pterm.Info.Println(desiredProjectWebhooks)
		pterm.Info.Println(actualWebhooks)
	}

	return nil
}
