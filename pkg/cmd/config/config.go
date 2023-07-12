package config

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "config",
	Short:   "Get or set config",
	Aliases: []string{"conf"},
}

func init() {
	Cmd.AddCommand(getConfigCmd)
	Cmd.AddCommand(setConfigCmd)
}
