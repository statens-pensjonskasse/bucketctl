package get

import (
	"git.spk.no/infra/bucketctl/pkg/api/bitbucket"
	"git.spk.no/infra/bucketctl/pkg/common"
	"git.spk.no/infra/bucketctl/pkg/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listRepositoriesCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString(common.ProjectKeyFlag) == "" {
			cmd.MarkFlagRequired(common.ProjectKeyFlag)
		}
		viper.BindPFlag(common.ProjectKeyFlag, cmd.Flags().Lookup(common.ProjectKeyFlag))
	},
	Use:     "repositories",
	Aliases: []string{"repo", "repos"},
	Short:   "List repositories in a given project",
	RunE:    listRepositories,
}

func listRepositories(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString(common.BaseUrlFlag)
	projectKey := viper.GetString(common.ProjectKeyFlag)
	limit := viper.GetInt(common.LimitFlag)
	token := viper.GetString(common.TokenFlag)

	repos, err := bitbucket.GetProjectRepositoriesMap(baseUrl, projectKey, limit, token)
	if err != nil {
		return err
	}

	return printer.PrintData(repos, printer.PrettyFormatRepositories)
}
