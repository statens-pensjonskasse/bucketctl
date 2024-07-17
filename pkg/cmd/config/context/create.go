package context

import (
	"errors"
	"git.spk.no/infra/bucketctl/pkg/common"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"path/filepath"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new context",
	Aliases: []string{"new"},
	RunE:    createContext,
}

func createContext(cmd *cobra.Command, args []string) error {
	cfgPath, err := common.GetConfigPath()
	if err != nil {
		return err
	}

	contextFile := filepath.Join(cfgPath, context+".yaml")
	if !common.FileNotExists(contextFile) {
		return errors.New("context '" + context + "' already exists")
	}

	if err := common.CreateFileIfNotExists(filepath.Join(cfgPath, context+".yaml"), 0600); err != nil {
		return err
	}

	config := make(map[string]interface{})
	if err := addEntriesFromCommandLine(cmd, &config); err != nil {
		return err
	}
	removeEmptyEntries(&config)

	yamlData, err := yaml.Marshal(&config)
	if err := common.WriteFile(contextFile, yamlData, 0600); err != nil {
		return err
	}

	pterm.Info.Println("🔧 Context '" + context + "' created")
	return nil
}
