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

var getConfigCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g", "list", "l"},
	Short:   "Get config",
	RunE:    getConfig,
}

var setConfigCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s"},
	Short:   "Set config",
	RunE:    setConfig,
}
