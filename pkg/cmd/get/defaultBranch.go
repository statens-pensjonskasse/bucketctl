package get

import (
	"git.spk.no/infra/bucketctl/pkg/api/bitbucket"
	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"git.spk.no/infra/bucketctl/pkg/common"
	"git.spk.no/infra/bucketctl/pkg/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listDefaultBranchCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString(common.ProjectKeyFlag) == "" {
			cmd.MarkFlagRequired(common.ProjectKeyFlag)
		}
		viper.BindPFlag(common.ProjectKeyFlag, cmd.Flags().Lookup(common.ProjectKeyFlag))
		viper.BindPFlag(common.RepoSlugFlag, cmd.Flags().Lookup(common.RepoSlugFlag))
	},
	Use:     "defaultBranch",
	Aliases: []string{"defaultBranches"},
	Run:     listDefaultBranches,
}

func listDefaultBranches(cmd *cobra.Command, args []string) {
	baseUrl := viper.GetString(common.BaseUrlFlag)
	projectKey := viper.GetString(common.ProjectKeyFlag)
	repoSlug := viper.GetString(common.RepoSlugFlag)
	limit := viper.GetInt(common.LimitFlag)
	token := viper.GetString(common.TokenFlag)

	projectConfig := ProjectConfigV1alpha1()
	projectConfig.Metadata.Name = projectKey

	if repoSlug == "" {
		defaultBranches, err := FetchDefaultBranches(baseUrl, projectKey, limit, token)
		cobra.CheckErr(err)
		projectConfig.Spec = *defaultBranches
	} else {
		repoDefaultBranch, err := bitbucket.GetRepositoryDefaultBranch(baseUrl, projectKey, repoSlug, token)
		cobra.CheckErr(err)
		projectConfig.Spec.ProjectKey = projectKey
		projectConfig.Spec.Repositories = &RepositoriesProperties{&RepositoryProperties{RepoSlug: repoSlug, DefaultBranch: &repoDefaultBranch.DisplayId}}
	}

	cobra.CheckErr(printer.PrintData(projectConfig, nil))
}

func FetchDefaultBranches(baseUrl string, projectKey string, limit int, token string) (*ProjectConfigSpec, error) {
	repositoriesProperties, err := bitbucket.GetRepositoriesDefaultBranch(baseUrl, projectKey, limit, token)
	if err != nil {
		return nil, err
	}

	mostCommonDefaultBranch := getMostCommonDefaultBranch(repositoriesProperties)
	spec := &ProjectConfigSpec{ProjectKey: projectKey, DefaultBranch: &mostCommonDefaultBranch, Repositories: repositoriesProperties}

	return FilterRepositoryDefaultBranchMatchingProjectDefaultBranch(spec), nil
}

func getMostCommonDefaultBranch(repositoriesProperties *RepositoriesProperties) string {
	defaultBranch := "N/A"

	if repositoriesProperties != nil {
		defaultBranchMap := make(map[string]int)
		for _, repo := range *repositoriesProperties {
			if _, exist := defaultBranchMap[*repo.DefaultBranch]; !exist {
				defaultBranchMap[*repo.DefaultBranch] = 0
			}
			defaultBranchMap[*repo.DefaultBranch]++
		}

		highestCount := 0
		for _, branchName := range common.GetLexicallySortedKeys(defaultBranchMap) {
			if count := defaultBranchMap[branchName]; count > highestCount {
				defaultBranch = branchName
				highestCount = count
			}
		}
	}

	return defaultBranch
}

func FilterRepositoryDefaultBranchMatchingProjectDefaultBranch(spec *ProjectConfigSpec) *ProjectConfigSpec {
	specCopy := spec.Copy()

	if spec.DefaultBranch != nil {
		for _, repo := range *specCopy.Repositories {
			if *repo.DefaultBranch == *spec.DefaultBranch {
				repo.DefaultBranch = nil
			}
		}
	}
	return specCopy
}
