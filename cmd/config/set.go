package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var setCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s"},
	Short:   "Set config",
	Run:     setConfig,
}

func setConfig(cmd *cobra.Command, args []string) {
	if err := viper.WriteConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving config")
		os.Exit(1)
	}
	getConfig(cmd, args)
	fmt.Fprintln(os.Stdout, "Saved config to: ", viper.ConfigFileUsed())
}
