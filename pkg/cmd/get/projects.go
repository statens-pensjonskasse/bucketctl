package get

import (
	"git.spk.no/infra/bucketctl/pkg/api/bitbucket"
	"git.spk.no/infra/bucketctl/pkg/common"
	"git.spk.no/infra/bucketctl/pkg/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listProjectsCmd = &cobra.Command{
	Use:     "projects",
	Aliases: []string{"proj"},
	Short:   "List projects",
	RunE:    listProjects,
}

func listProjects(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString(common.BaseUrlFlag)
	token := viper.GetString(common.TokenFlag)
	limit := viper.GetInt(common.LimitFlag)

	projects, err := bitbucket.GetProjects(baseUrl, limit, token)
	if err != nil {
		return err
	}

	return printer.PrintData(projects, printer.PrettyFormatProjects)
}
