package project

import (
	"bucketctl/pkg/common"
	"bucketctl/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listProjectsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List projects",
	RunE:    listProjects,
}

func listProjects(cmd *cobra.Command, args []string) error {
	var baseUrl = viper.GetString(types.BaseUrlFlag)
	var token = viper.GetString(types.TokenFlag)
	var limit = viper.GetInt(types.LimitFlag)

	projects, err := GetProjects(baseUrl, token, limit)
	if err != nil {
		return err
	}

	return common.PrintData(projects, prettyFormatProjects)
}
