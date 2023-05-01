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
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(setCmd)
}

var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g", "list", "l"},
	Short:   "Get config",
	Run:     getConfig,
}

var setCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s"},
	Short:   "Set config",
	RunE:    setConfig,
}
