package config

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sort"
)

var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g", "list", "l"},
	Short:   "Get config",
	Run:     getConfig,
}

func getConfig(cmd *cobra.Command, args []string) {
	var keys = viper.AllKeys()
	sort.Strings(keys)

	var data [][]string
	data = append(data, []string{"Key", "Value"})

	for _, key := range keys {
		row := []string{key, viper.GetString(key)}
		data = append(data, row)
	}

	pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}
