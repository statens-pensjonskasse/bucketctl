package apply

import (
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/cmd/apply/access"
	"bucketctl/pkg/cmd/apply/branchRestrictions"
	"bucketctl/pkg/cmd/apply/webhooks"
	"bucketctl/pkg/cmd/get"
	"bucketctl/pkg/common"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired(common.FilenameFlag)
		viper.BindPFlag(common.FilenameFlag, cmd.Flags().Lookup(common.FilenameFlag))
		viper.BindPFlag(common.DryRunFlag, cmd.Flags().Lookup(common.DryRunFlag))
	},
	Use:   "apply",
	Short: "Apply bucketctl manifest",
	Run:   applyProjectConfig,
}

func init() {
	Cmd.Flags().StringP(common.FilenameFlag, common.FilenameFlagShorthand, "", "")
	Cmd.Flags().BoolP(common.DryRunFlag, common.DryRunFlagShorthand, false, "Dry run")
}

func applyProjectConfig(cmd *cobra.Command, args []string) {
	baseUrl := viper.GetString(common.BaseUrlFlag)
	limit := viper.GetInt(common.LimitFlag)
	token := viper.GetString(common.TokenFlag)

	file := viper.GetString(common.FilenameFlag)
	dryRun := viper.GetBool(common.DryRunFlag)

	desired, err := readProjectConfig(file)
	cobra.CheckErr(err)

	actual, err := getActualProjectConfig(baseUrl, desired.Spec.ProjectKey, limit, token)
	cobra.CheckErr(err)

	toCreate, toUpdate, toDelete := findProjectConfigChanges(&desired.Spec, actual)

	if dryRun {
		printChanges(toCreate, toUpdate, toDelete)
	} else {
		err := setProjectConfig(baseUrl, desired.Spec.ProjectKey, token, toCreate, toUpdate, toDelete)
		cobra.CheckErr(err)
	}
}

func findProjectConfigChanges(desired *ProjectConfigSpec, actual *ProjectConfigSpec) (toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) {
	accessToCreate, accessToUpdate, accessToDelete := access.FindAccessChanges(desired, actual)
	brToCreate, brToUpdate, brToDelete := branchRestrictions.FindBranchRestrictionChanges(desired, actual)
	whToCreate, whToUpdate, whToDelete := webhooks.FindWebhooksChanges(desired, actual)

	toCreate = CombineProjectConfigSpecs(accessToCreate, brToCreate, whToCreate)
	toUpdate = CombineProjectConfigSpecs(accessToUpdate, brToUpdate, whToUpdate)
	toDelete = CombineProjectConfigSpecs(accessToDelete, brToDelete, whToDelete)

	return toCreate, toUpdate, toDelete
}

func setProjectConfig(baseUrl string, projectKey string, token string, toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) error {
	if err := access.SetAccess(baseUrl, projectKey, token, toCreate, toUpdate, toDelete); err != nil {
		return err
	}
	if err := branchRestrictions.SetBranchRestrictions(baseUrl, projectKey, token, toCreate, toUpdate, toDelete); err != nil {
		return err
	}
	if err := webhooks.SetWebhooks(baseUrl, projectKey, token, toCreate, toUpdate, toDelete); err != nil {
		return err
	}

	return nil
}

func getActualProjectConfig(baseUrl string, projectKey string, limit int, token string) (*ProjectConfigSpec, error) {
	actualAccess, err := get.FetchAccess(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}
	actualBranchRestrictions, err := get.FetchBranchRestrictions(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}
	actualWebhooks, err := get.FetchWebhooks(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	return CombineProjectConfigSpecs(actualAccess, actualBranchRestrictions, actualWebhooks), nil
}

func readProjectConfig(file string) (*ProjectConfig, error) {
	var projectConfig ProjectConfig
	if err := common.ReadConfigFile(file, &projectConfig); err != nil {
		return nil, err
	}
	if err := projectConfig.Validate(); err != nil {
		return nil, err
	}

	return &projectConfig, nil
}

func printChanges(toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) {
	pterm.Println("--- Planned changes ---")
	pterm.Println("Access:")
	access.PrintAccessChanges(toCreate, toUpdate, toDelete)
	pterm.Println("Branch restrictions:")
	branchRestrictions.PrintBranchRestrictionChanges(toCreate, toUpdate, toDelete)
	pterm.Println("Webhooks:")
	webhooks.PrintWebhookChanges(toCreate, toUpdate, toDelete)
}
