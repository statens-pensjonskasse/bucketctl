package get

import (
	"git.spk.no/infra/bucketctl/pkg/common"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get information about a given resource",
}

func init() {
	Cmd.MarkFlagRequired(common.BaseUrlFlag)
	Cmd.PersistentFlags().StringP(common.ProjectKeyFlag, common.ProjectKeyFlagShorthand, "", "Project key")
	Cmd.PersistentFlags().StringP(common.RepoSlugFlag, common.RepoSlugFlagShorthand, "", "Repository slug. Leave empty to query permission permissions.")

	Cmd.AddCommand(listAccessCmd)
	Cmd.AddCommand(listBranchingModelsCmd)
	Cmd.AddCommand(listBranchRestrictionsCmd)
	Cmd.AddCommand(listDefaultBranchCmd)
	Cmd.AddCommand(listWebhooksCmd)

	Cmd.AddCommand(listProjectsCmd)
	Cmd.AddCommand(listRepositoriesCmd)

	Cmd.AddCommand(getProjectConfigCmd)
}
