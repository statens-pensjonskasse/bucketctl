package printer

import (
	"strings"

	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
	"git.spk.no/infra/bucketctl/pkg/common"
)

func PrettyFormatProjectsSettings(projectConfig *ProjectConfig) [][]string {
	projectKey := projectConfig.Spec.ProjectKey
	projectBranchRestrictions := projectConfig.Spec.BranchRestrictions
	projectRepositories := projectConfig.Spec.Repositories

	var data [][]string

	data = append(data, []string{"Project", "Repository", "Matcher Type", "Matching", "Restriction", "Exempt Groups", "Exempt Users"})
	projectSettings := prettyFormatRestrictions(projectKey, ALL, projectBranchRestrictions)
	data = append(data, projectSettings...)

	repoSettings := prettyFormatRepositoryProperties(projectKey, projectRepositories)
	data = append(data, repoSettings...)

	return data
}

func prettyFormatRepositoryProperties(projectKey string, repositoriesProperties *RepositoriesProperties) [][]string {
	var data [][]string
	if repositoriesProperties == nil {
		return data
	}

	repoSettingsMap := make(map[string]*BranchRestrictions, len(*repositoriesProperties))
	for _, r := range *repositoriesProperties {
		repoSettingsMap[r.RepoSlug] = r.BranchRestrictions
	}

	repositories := common.GetLexicallySortedKeys(repoSettingsMap)
	for _, repoSlug := range repositories {
		repoSettings := prettyFormatRestrictions(projectKey, repoSlug, repoSettingsMap[repoSlug])
		data = append(data, repoSettings...)
	}

	return data
}

func prettyFormatRestrictions(projectKey string, repoSlug string, branchRestrictions *BranchRestrictions) [][]string {
	var data [][]string
	if branchRestrictions == nil {
		return data
	}

	matcherRestrictionsMap := make(map[string]*BranchMatchers, len(*branchRestrictions))
	for _, b := range *branchRestrictions {
		matcherRestrictionsMap[b.Type] = b.BranchMatchers
	}

	for matcher, restriction := range matcherRestrictionsMap {
		for _, branchRestriction := range *restriction {
			restrictions := prettyFormatBranchRestrictions(projectKey, repoSlug, matcher, branchRestriction)
			data = append(data, restrictions...)
		}
	}

	return data
}

func prettyFormatBranchRestrictions(projectKey string, repoSlug string, matcher string, branchMatcherRestriction *BranchMatcher) [][]string {
	var data [][]string
	if branchMatcherRestriction == nil {
		return data
	}

	restrictionsMap := make(map[string]*Restriction)
	for _, r := range *branchMatcherRestriction.Restrictions {
		restrictionsMap[r.Type] = r
	}

	branch := branchMatcherRestriction.Matching
	restrictions := common.GetLexicallySortedKeys(restrictionsMap)
	for _, restriction := range restrictions {
		var users string
		for _, user := range restrictionsMap[restriction].ExemptUsers {
			users += user + "\n"
		}
		users = strings.Trim(users, "\n")

		var groups string
		for _, group := range restrictionsMap[restriction].ExemptGroups {
			groups += group + "\n"
		}
		groups = strings.Trim(groups, "\n")

		data = append(data, []string{projectKey, repoSlug, matcher, branch, restriction, groups, users})
		projectKey = ""
		repoSlug = ""
		matcher = ""
		branch = ""
	}

	return data
}
