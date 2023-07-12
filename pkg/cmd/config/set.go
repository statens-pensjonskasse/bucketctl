package config

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setConfigCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s"},
	Short:   "Set config",
	RunE:    setConfig,
}

func setConfig(cmd *cobra.Command, args []string) error {
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	getConfig(cmd, args)
	pterm.Info.Println("Saved config to:", viper.ConfigFileUsed())

	return nil
}
