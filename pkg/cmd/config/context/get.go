package context

import (
	"bucketctl/pkg/common"
	"bucketctl/pkg/printer"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g", "list", "l"},
	Short:   "View context",
	RunE:    getContext,
}

func getContext(cmd *cobra.Command, args []string) error {
	contextFile, err := getContextFilename(context)
	if err != nil {
		return err
	}

	var config map[string]string
	if err := common.ReadConfigFile(contextFile, &config); err != nil {
		return err
	}

	return printer.PrintData(config, prettyFormatContext)
}
