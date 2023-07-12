package repository

import (
	"bucketctl/pkg"
	"bucketctl/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	key string
)

var listRepositoriesCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag(types.ProjectKeyFlag, cmd.Flags().Lookup(types.ProjectKeyFlag))
	},
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List repositories in a given project",
	RunE:    listRepositories,
}

func init() {
	listRepositoriesCmd.Flags().StringVarP(&key, types.ProjectKeyFlag, "k", "", "Project key")
	listRepositoriesCmd.MarkFlagRequired(types.ProjectKeyFlag)
}

func listRepositories(cmd *cobra.Command, args []string) error {
	var baseUrl = viper.GetString(types.BaseUrlFlag)
	var projectKey = viper.GetString(types.ProjectKeyFlag)
	var limit = viper.GetInt(types.LimitFlag)

	repos, err := GetProjectRepositories(baseUrl, projectKey, limit)
	if err != nil {
		return err
	}

	return pkg.PrintData(repos, prettyFormatRepositories)
}
