package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GoBit version 0.0.1")
	},
}
