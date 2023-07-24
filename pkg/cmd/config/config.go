package config

import (
	"bucketctl/pkg/cmd/config/context"
	"bucketctl/pkg/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &cobra.Command{
	Use:     "config",
	Short:   "Get or set config",
	Aliases: []string{"conf"},
}

func init() {
	Cmd.AddCommand(context.Cmd)
	Cmd.AddCommand(getConfigCmd)
	Cmd.AddCommand(setConfigCmd)
}

func prettyFormatConfig(settingMap map[string]interface{}) [][]string {
	var data [][]string
	data = append(data, []string{"Key", "Value"})

	keys := common.GetLexicallySortedKeys(settingMap)
	for _, key := range keys {
		row := []string{key, viper.GetString(key)}
		data = append(data, row)
	}

	return data
}
