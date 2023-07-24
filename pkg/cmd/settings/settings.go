package settings

import (
	"bucketctl/pkg/cmd/repository"
	"bucketctl/pkg/common"
	"bucketctl/pkg/types"
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"strings"
)

type ProjectSettings struct {
	Restrictions map[string]*Restrictions       `json:"restrictions,omitempty" yaml:"restrictions,omitempty"`
	Repositories map[string]*RepositorySettings `json:"repositories,omitempty" yaml:"repositories,omitempty"`
}

type RepositorySettings struct {
	Restrictions map[string]*Restrictions `json:"restrictions,omitempty" yaml:"restrictions,omitempty"`
}

type Restrictions struct {
	Branches map[string]*BranchRestrictions `json:"branch,omitempty" yaml:",inline,inline,omitempty"`
}

type BranchRestrictions struct {
	Restrictions map[string]*Restriction `json:"restriction,omitempty" yaml:"restriction,inline,omitempty"`
}

type Restriction struct {
	id           int
	ExemptUsers  []string `json:"exempt-users,omitempty" yaml:"exempt-users,omitempty"`
	ExemptGroups []string `json:"exempt-groups,omitempty" yaml:"exempt-groups,omitempty"`
}

var (
	key  string
	repo string
)

var Cmd = &cobra.Command{
	Use:     "settings",
	Short:   "View and edit Bitbucket project and repository settings",
	Aliases: []string{"setting"},
}

func init() {
	Cmd.AddCommand(listSettingsCmd)
	Cmd.AddCommand(listAllSettingsCmd)
	Cmd.AddCommand(applySettingsCmd)
}

func getProjectRestrictions(baseUrl string, projectKey string, limit int, token string, includeRepos bool) (*ProjectSettings, error) {
	url := fmt.Sprintf("%s/rest/branch-permissions/latest/projects/%s/restrictions?limit=%d", baseUrl, projectKey, limit)

	restrictions, err := getRestrictions(url, token)
	if err != nil {
		return nil, err
	}

	projectRestrictions := make(map[string]*Restrictions)
	for _, r := range restrictions {
		if _, exists := projectRestrictions[r.Matcher.Type.Id]; !exists {
			projectRestrictions[r.Matcher.Type.Id] = new(Restrictions)
			projectRestrictions[r.Matcher.Type.Id].Branches = make(map[string]*BranchRestrictions)
		}
		projectRestrictions[r.Matcher.Type.Id].addRestriction(r)
	}

	projectSettings := &ProjectSettings{Restrictions: projectRestrictions}

	if includeRepos {
		projectRepositories, err := repository.GetProjectRepositories(baseUrl, projectKey, token, limit)
		if err != nil {
			return nil, err
		}
		projectSettings.Repositories = make(map[string]*RepositorySettings)
		for repoSlug := range projectRepositories {
			repoSettings, err := getRepositoryRestrictions(baseUrl, projectKey, repoSlug, limit, token)
			if err != nil {
				return nil, err
			}
			if len(repoSettings.Restrictions) > 0 {
				projectSettings.Repositories[repoSlug] = repoSettings
			}
		}

	}

	return projectSettings, nil
}

func getRepositoryRestrictions(baseUrl string, projectKey string, repoSlug string, limit int, token string) (*RepositorySettings, error) {
	url := fmt.Sprintf("%s/rest/branch-permissions/latest/projects/%s/repos/%s/restrictions?limit=%d", baseUrl, projectKey, repoSlug, limit)
	restrictions, err := getRestrictions(url, token)
	if err != nil {
		return nil, err
	}

	repoRestrictions := make(map[string]*Restrictions)
	for _, r := range restrictions {
		if r.Scope.Type != "REPOSITORY" {
			continue
		}
		if _, exists := repoRestrictions[r.Matcher.Type.Id]; !exists {
			repoRestrictions[r.Matcher.Type.Id] = new(Restrictions)
			repoRestrictions[r.Matcher.Type.Id].Branches = make(map[string]*BranchRestrictions)
		}
		repoRestrictions[r.Matcher.Type.Id].addRestriction(r)
	}

	return &RepositorySettings{Restrictions: repoRestrictions}, nil
}

func (restrictions *Restrictions) addRestriction(r *types.Restriction) {
	if _, exists := restrictions.Branches[r.Matcher.Id]; !exists {
		restrictions.Branches[r.Matcher.Id] = new(BranchRestrictions)
		restrictions.Branches[r.Matcher.Id].Restrictions = make(map[string]*Restriction)
	}
	var users []string
	for _, u := range r.Users {
		users = append(users, u.Name)
	}
	restrictions.Branches[r.Matcher.Id].Restrictions[r.Type] = &Restriction{
		id:           r.Id,
		ExemptUsers:  users,
		ExemptGroups: r.Groups,
	}
}

func getRestrictions(url string, token string) ([]*types.Restriction, error) {
	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var restrictions types.RestrictionResponse
	if err := json.Unmarshal(body, &restrictions); err != nil {
		return nil, err
	}

	if !restrictions.IsLastPage {
		pterm.Warning.Println("Not all restrictions fetched, try with a higher limit")
	}

	return restrictions.Values, nil
}

func prettyFormatProjectsSettings(projectSettingsMap map[string]*ProjectSettings) [][]string {
	var data [][]string

	data = append(data, []string{"Project", "Repository", "Matcher", "Branch", "Restriction", "Exempt Groups", "Exempt Users"})
	projects := common.GetLexicallySortedKeys(projectSettingsMap)
	for _, projectKey := range projects {
		projectSettings := prettyFormatRestrictions(projectKey, "ALL", projectSettingsMap[projectKey].Restrictions)
		data = append(data, projectSettings...)

		repoSettings := prettyFormatRepositorySettings(projectKey, projectSettingsMap[projectKey].Repositories)
		data = append(data, repoSettings...)
	}

	return data
}

func prettyFormatRepositorySettings(projectKey string, repoSettingsMap map[string]*RepositorySettings) [][]string {
	var data [][]string

	repositories := common.GetLexicallySortedKeys(repoSettingsMap)
	for _, repoSlug := range repositories {
		repoSettings := prettyFormatRestrictions(projectKey, repoSlug, repoSettingsMap[repoSlug].Restrictions)
		data = append(data, repoSettings...)
	}

	return data
}

func prettyFormatRestrictions(projectKey string, repoSlug string, restrictions map[string]*Restrictions) [][]string {
	var data [][]string

	for matcher, restriction := range restrictions {
		for branch, branchRestriction := range restriction.Branches {
			branchRestrictions := prettyFormatBranchRestrictions(projectKey, repoSlug, matcher, branch, branchRestriction)
			data = append(data, branchRestrictions...)
		}
	}

	return data
}

func prettyFormatBranchRestrictions(projectKey string, repoSlug string, matcher string, branch string, branchRestrictions *BranchRestrictions) [][]string {
	var data [][]string

	restrictions := common.GetLexicallySortedKeys(branchRestrictions.Restrictions)
	for _, restriction := range restrictions {
		var users string
		for _, user := range branchRestrictions.Restrictions[restriction].ExemptUsers {
			users += user + "\n"
		}
		users = strings.Trim(users, "\n")

		var groups string
		for _, group := range branchRestrictions.Restrictions[restriction].ExemptGroups {
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
