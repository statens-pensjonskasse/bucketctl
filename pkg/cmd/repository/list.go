package repository

import (
	"bucketctl/pkg/common"
	"bucketctl/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	key string
)

var listRepositoriesCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString(types.ProjectKeyFlag) == "" {
			cmd.MarkFlagRequired(types.ProjectKeyFlag)
		}
		viper.BindPFlag(types.ProjectKeyFlag, cmd.Flags().Lookup(types.ProjectKeyFlag))
	},
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List repositories in a given project",
	RunE:    listRepositories,
}

func init() {
	listRepositoriesCmd.Flags().StringVarP(&key, types.ProjectKeyFlag, types.ProjectKeyFlagShorthand, "", "Project key")
}

func listRepositories(cmd *cobra.Command, args []string) error {
	var baseUrl = viper.GetString(types.BaseUrlFlag)
	var projectKey = viper.GetString(types.ProjectKeyFlag)
	var token = viper.GetString(types.TokenFlag)
	var limit = viper.GetInt(types.LimitFlag)

	repos, err := GetProjectRepositories(baseUrl, projectKey, token, limit)
	if err != nil {
		return err
	}

	return common.PrintData(repos, prettyFormatRepositories)
}
