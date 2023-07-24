package context

import (
	"bucketctl/pkg/common"
	"errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"path/filepath"
)

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete context",
	Aliases: []string{"del", "rm"},
	RunE:    deleteContext,
}

func deleteContext(cmd *cobra.Command, args []string) error {
	if context == "config" {
		return errors.New("can't delete base config")
	}

	cfgPath, err := common.GetConfigPath()
	if err != nil {
		return err
	}

	contextFile := filepath.Join(cfgPath, context+".yaml")
	if common.FileNotExists(contextFile) {
		return errors.New("context '" + context + "' doesn't exists")
	}

	if err := common.RemoveFile(contextFile); err != nil {
		return err
	}

	pterm.Info.Println("🗑️ Context '" + context + "' deleted")
	return nil
}
