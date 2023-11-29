package apply

import (
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/cmd/apply/access"
	"bucketctl/pkg/cmd/apply/branchRestrictions"
	"bucketctl/pkg/cmd/apply/branchingModel"
	"bucketctl/pkg/cmd/apply/defaultBranch"
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
		viper.BindPFlag(common.CompareWithFileFlag, cmd.Flags().Lookup(common.CompareWithFileFlag))
		viper.BindPFlag(common.DryRunFlag, cmd.Flags().Lookup(common.DryRunFlag))
	},
	Use:   "apply",
	Short: "Apply bucketctl manifest",
	Run:   applyProjectConfig,
}

func init() {
	Cmd.MarkFlagRequired(common.BaseUrlFlag)
	Cmd.MarkFlagRequired(common.TokenFlag)
	Cmd.Flags().StringP(common.FilenameFlag, common.FilenameFlagShorthand, "", "")
	Cmd.Flags().String(common.CompareWithFileFlag, "", "Compare with file, implicitly dry-run")
	Cmd.Flags().BoolP(common.DryRunFlag, common.DryRunFlagShorthand, false, "Dry run")
}

func applyProjectConfig(cmd *cobra.Command, args []string) {
	baseUrl := viper.GetString(common.BaseUrlFlag)
	limit := viper.GetInt(common.LimitFlag)
	token := viper.GetString(common.TokenFlag)

	file := viper.GetString(common.FilenameFlag)

	compareFile := viper.GetString(common.CompareWithFileFlag)

	desired, err := readProjectConfig(file)
	cobra.CheckErr(err)

	var actual *ProjectConfigSpec
	if len(compareFile) > 0 {
		viper.Set(common.DryRunFlag, true)
		comparison, err := readProjectConfig(compareFile)
		cobra.CheckErr(err)
		actual = &comparison.Spec
	} else {
		actual, err = get.FetchProjectConfigSpec(baseUrl, desired.Spec.ProjectKey, limit, token)
		cobra.CheckErr(err)
	}

	toCreate, toUpdate, toDelete := findProjectConfigChanges(&desired.Spec, actual)

	dryRun := viper.GetBool(common.DryRunFlag)
	if dryRun {
		printChanges(toCreate, toUpdate, toDelete)
	} else {
		err := setProjectConfig(baseUrl, desired.Spec.ProjectKey, token, toCreate, toUpdate, toDelete)
		cobra.CheckErr(err)
	}
}

func findProjectConfigChanges(desired *ProjectConfigSpec, actual *ProjectConfigSpec) (toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) {
	accessToCreate, accessToUpdate, accessToDelete := access.FindAccessChanges(desired, actual)
	bmToCreate, bmToUpdate, bmToDelete := branchingModel.FindBranchingModelChanges(desired, actual)
	brToCreate, brToUpdate, brToDelete := branchRestrictions.FindBranchRestrictionChanges(desired, actual)
	defaultBranchesToUpdate := defaultBranch.FindDefaultBranchChanges(desired, actual)
	whToCreate, whToUpdate, whToDelete := webhooks.FindWebhooksChanges(desired, actual)

	uncombinedToCreate := &UncombinedProjectConfigSpecs{
		Access:             accessToCreate,
		BranchingModels:    bmToCreate,
		BranchRestrictions: brToCreate,
		DefaultBranches:    nil,
		Webhooks:           whToCreate,
	}

	uncombinedToUpdate := &UncombinedProjectConfigSpecs{
		Access:             accessToUpdate,
		BranchingModels:    bmToUpdate,
		BranchRestrictions: brToUpdate,
		DefaultBranches:    defaultBranchesToUpdate,
		Webhooks:           whToUpdate,
	}

	uncombinedToDelete := &UncombinedProjectConfigSpecs{
		Access:             accessToDelete,
		BranchingModels:    bmToDelete,
		BranchRestrictions: brToDelete,
		DefaultBranches:    nil,
		Webhooks:           whToDelete,
	}

	toCreate = CombineProjectConfigSpecs(uncombinedToCreate)
	toUpdate = CombineProjectConfigSpecs(uncombinedToUpdate)
	toDelete = CombineProjectConfigSpecs(uncombinedToDelete)

	return toCreate, toUpdate, toDelete
}

func setProjectConfig(baseUrl string, projectKey string, token string, toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) error {
	if err := access.SetAccess(baseUrl, projectKey, token, toCreate, toUpdate, toDelete); err != nil {
		return err
	}
	if err := branchingModel.SetBranchingModels(baseUrl, projectKey, token, toCreate, toUpdate, toDelete); err != nil {
		return err
	}
	if err := branchRestrictions.SetBranchRestrictions(baseUrl, projectKey, token, toCreate, toUpdate, toDelete); err != nil {
		return err
	}
	if err := defaultBranch.SetDefaultBranches(baseUrl, projectKey, token, toUpdate); err != nil {
		return err
	}
	if err := webhooks.SetWebhooks(baseUrl, projectKey, token, toCreate, toUpdate, toDelete); err != nil {
		return err
	}

	return nil
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
	var changes []string
	changes = append(changes, access.GetChangesAsText(toCreate, toUpdate, toDelete)...)
	changes = append(changes, branchingModel.GetChangesAsText(toCreate, toUpdate, toDelete)...)
	changes = append(changes, branchRestrictions.GetChangesAsText(toCreate, toUpdate, toDelete)...)
	changes = append(changes, defaultBranch.GetChangesAsText(toUpdate)...)
	changes = append(changes, webhooks.GetChangesAsText(toCreate, toUpdate, toDelete)...)

	if changes != nil && len(changes) > 0 {
		for _, change := range changes {
			pterm.Printfln(change)
		}
	} else {
		pterm.Printfln("No changes in project %s", pterm.Bold.Sprint(toCreate.ProjectKey))
	}

}
