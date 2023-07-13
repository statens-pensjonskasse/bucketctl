package context

import (
	"bucketctl/pkg"
	"errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"path/filepath"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update context",
	Aliases: []string{"upd"},
	RunE:    updateContext,
}

func init() {

}

func updateContext(cmd *cobra.Command, args []string) error {
	cfgPath, err := pkg.GetConfigPath()
	if err != nil {
		return err
	}

	contextFile := filepath.Join(cfgPath, context+".yaml")
	if pkg.FileNotExists(contextFile) {
		return errors.New("context '" + context + "' doesn't exists")
	}

	var config map[string]interface{}
	if err := pkg.ReadConfigFile(contextFile, &config); err != nil {
		return err
	}

	if err := addEntriesFromCommandLine(cmd, &config); err != nil {
		return err
	}
	removeEmptyEntries(&config)

	yamlData, err := yaml.Marshal(&config)
	return pkg.WriteFile(contextFile, yamlData, 0600)
}
