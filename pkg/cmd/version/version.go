package version

import (
	"bucketctl/pkg/logger"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const version = "1.4.2"

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Print version number",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Log("🪣 %s version %s 🔧", pterm.Blue("bucketctl"), pterm.White(version))
	},
}
