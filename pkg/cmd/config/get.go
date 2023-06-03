package config

import (
	"bucketctl/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sort"
)

func prettyFormatConfig(settings map[string]interface{}) [][]string {
	var data [][]string
	data = append(data, []string{"Key", "Value"})

	for key := range settings {
		row := []string{key, viper.GetString(key)}
		data = append(data, row)
	}

	return data
}

func getConfig(cmd *cobra.Command, args []string) error {
	keys := viper.AllKeys()
	sort.Strings(keys)

	var data [][]string
	data = append(data, []string{"Key", "Value"})

	for _, key := range keys {
		row := []string{key, viper.GetString(key)}
		data = append(data, row)
	}

	return pkg.PrintData(viper.AllSettings(), prettyFormatConfig)
}
