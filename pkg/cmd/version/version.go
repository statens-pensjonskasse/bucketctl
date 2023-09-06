package version

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const version = "0.1.0"

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Print version number",
	Run: func(cmd *cobra.Command, args []string) {
		pterm.Println("bucketctl version", version)
	},
}
