package context

import (
	"bucketctl/pkg"
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
	if err := pkg.ReadConfigFile(contextFile, &config); err != nil {
		return err
	}

	return pkg.PrintData(config, prettyFormatContext)
}
