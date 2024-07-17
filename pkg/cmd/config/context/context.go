package context

import (
	"errors"
	"git.spk.no/infra/bucketctl/pkg/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	context string
	key     string
	repo    string
)

var Cmd = &cobra.Command{
	Use:     "context",
	Aliases: []string{"ctx"},
	Short:   "Create, edit and delete contexts",
}

func init() {
	Cmd.PersistentFlags().StringVarP(&context, common.ContextFlag, common.ContextFlagShorthand, "", "Context to use")
	Cmd.MarkPersistentFlagRequired(common.ContextFlag)

	Cmd.PersistentFlags().StringVarP(&key, common.ProjectKeyFlag, common.ProjectKeyFlagShorthand, "", "Project key")
	Cmd.PersistentFlags().StringVarP(&repo, common.RepoSlugFlag, common.RepoSlugFlagShorthand, "", "Repository slug")

	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(updateCmd)
}

func getContextFilename(ctx string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(home, ".config", "bucketctl")
	for _, ext := range viper.SupportedExts {
		ctxFilename := filepath.Join(configPath, ctx+"."+ext)
		if _, err := os.Stat(ctxFilename); err == nil {
			return ctxFilename, nil
		}
	}

	return "", errors.New("context '" + ctx + "' not found in '" + configPath + "'")
}

func prettyFormatContext(contextMap map[string]string) [][]string {
	var data [][]string
	data = append(data, []string{"Key", "Value"})

	keys := common.GetLexicallySortedKeys(contextMap)
	for _, k := range keys {
		row := []string{k, contextMap[k]}
		data = append(data, row)
	}

	return data
}

func addEntryIfChanged(cmd *cobra.Command, stringMap *map[string]interface{}, flag string) {
	if cmd.Flags().Changed(flag) {
		if flag == common.OutputFlag && viper.GetString(common.OutputFlag) != "" {
			(*stringMap)[flag] = viper.GetString(common.OutputFlag)
			return
		}
		if _, err := cmd.Flags().GetString(flag); err == nil {
			(*stringMap)[flag], _ = cmd.Flags().GetString(flag)
			return
		}
		if _, err := cmd.Flags().GetInt(flag); err == nil {
			(*stringMap)[flag], _ = cmd.Flags().GetInt(flag)
			return
		}
		if _, err := cmd.Flags().GetBool(flag); err == nil {
			(*stringMap)[flag], _ = cmd.Flags().GetBool(flag)
			return
		}
	}
}

func removeEmptyEntries(stringMap *map[string]interface{}) {
	for k, v := range *stringMap {
		if v == "" || v == 0 || v == nil {
			delete(*stringMap, k)
		}
	}
}

func addEntriesFromCommandLine(cmd *cobra.Command, config *map[string]interface{}) error {
	flags := []string{
		common.BaseUrlFlag,
		common.LimitFlag,
		common.TokenFlag,
		common.OutputFlag,
		common.ProjectKeyFlag,
		common.RepoSlugFlag,
	}

	for _, flag := range flags {
		addEntryIfChanged(cmd, config, flag)
	}

	return nil
}
