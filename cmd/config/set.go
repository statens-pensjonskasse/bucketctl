package config

import (
	"github.com/pterm/pterm"
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
		pterm.Error.Println("Error saving config")
		os.Exit(1)
	}
	getConfig(cmd, args)
	pterm.Info.Println("Saved config to:", viper.ConfigFileUsed())
}
