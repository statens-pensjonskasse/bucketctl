package config

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "config",
	Short: "Get or set config",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
	RootCmd.AddCommand(setCmd)
}
