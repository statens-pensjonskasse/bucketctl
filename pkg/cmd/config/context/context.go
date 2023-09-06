package context

import (
	"bucketctl/pkg/common"
	"bucketctl/pkg/types"
	"errors"
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
	Cmd.PersistentFlags().StringVarP(&context, types.ContextFlag, types.ContextFlagShorthand, "", "Context to use")
	Cmd.MarkPersistentFlagRequired(types.ContextFlag)

	Cmd.PersistentFlags().StringVarP(&key, types.ProjectKeyFlag, types.ProjectKeyFlagShorthand, "", "Project key")
	Cmd.PersistentFlags().StringVarP(&repo, types.RepoSlugFlag, types.RepoSlugFlagShorthand, "", "Repository slug")
	Cmd.PersistentFlags().Bool(types.IncludeReposFlag, false, "Include repository permissions when querying project permissions")

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
		if flag == types.OutputFlag && viper.GetString(types.OutputFlag) != "" {
			(*stringMap)[flag] = viper.GetString(types.OutputFlag)
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
		types.BaseUrlFlag,
		types.LimitFlag,
		types.TokenFlag,
		types.OutputFlag,
		types.ProjectKeyFlag,
		types.RepoSlugFlag,
		types.IncludeReposFlag,
	}

	for _, flag := range flags {
		addEntryIfChanged(cmd, config, flag)
	}

	return nil
}
