package project

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "project",
	Short:   "Bitbucket project commands",
	Aliases: []string{"p"},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
