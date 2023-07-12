package settings

import (
	"bucketctl/pkg"
	"bucketctl/pkg/types"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
)

var (
	filename string
)

var applySettingsCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag(types.FilenameFlag, cmd.Flags().Lookup(types.FilenameFlag))
		viper.BindPFlag(types.IncludeReposFlag, cmd.Flags().Lookup(types.IncludeReposFlag))
	},
	Use:   "apply",
	Short: "Apply a settings configuration file",
	RunE:  applySettings,
}

func init() {
	applySettingsCmd.Flags().StringVarP(&filename, types.FilenameFlag, "f", "", "Settings file")
	applySettingsCmd.Flags().Bool(types.IncludeReposFlag, false, "Include repositories")

	applySettingsCmd.MarkFlagRequired(types.FilenameFlag)
}

func applySettings(cmd *cobra.Command, args []string) error {
	file := viper.GetString(types.FilenameFlag)
	baseUrl := viper.GetString(types.BaseUrlFlag)
	limit := viper.GetInt(types.LimitFlag)
	token := viper.GetString(types.TokenFlag)
	includeRepos := viper.GetBool(types.IncludeReposFlag)

	var desiredSettings map[string]*ProjectSettings
	if err := pkg.ReadConfigFile(file, &desiredSettings); err != nil {
		return err
	}

	projectKeys := pkg.GetLexicallySortedKeys(desiredSettings)
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(len(desiredSettings)).WithRemoveWhenDone(true).Start()
	for _, projectKey := range projectKeys {
		progressBar.UpdateTitle(projectKey)
		// Finn gjeldende settings
		actualProjectSettings, err := getProjectRestrictions(baseUrl, projectKey, limit, token, includeRepos)
		if err != nil {
			return err
		}

		toUpdate := findRestrictionsToUpdate(desiredSettings[projectKey].Restrictions, actualProjectSettings.Restrictions)
		toDelete := findRestrictionsToDelete(desiredSettings[projectKey].Restrictions, actualProjectSettings.Restrictions)

		if err := updateProjectRestriction(baseUrl, projectKey, token, toUpdate); err != nil {
			return err
		}
		if err := deleteProjectRestrictions(baseUrl, projectKey, token, toDelete); err != nil {
			return err
		}

		progressBar.Increment()
	}

	return nil
}

func findRestrictionsToUpdate(desired map[string]*Restrictions, actual map[string]*Restrictions) (toUpdate map[string]*Restrictions) {
	toUpdate = make(map[string]*Restrictions)

	for matcherType, matcherRestriction := range desired {
		if _, exists := actual[matcherType]; !exists {
			if toUpdate[matcherType] == nil {
				toUpdate[matcherType] = new(Restrictions)
			}
			toUpdate[matcherType] = matcherRestriction
		} else {
			for branchType, branchRestriction := range desired[matcherType].Branches {
				if _, exists := actual[matcherType].Branches[branchType]; !exists {
					if toUpdate[matcherType] == nil {
						toUpdate[matcherType] = new(Restrictions)
						toUpdate[matcherType].Branches = make(map[string]*BranchRestrictions)
					}
					toUpdate[matcherType].Branches[branchType] = branchRestriction
				} else {
					for restrictionType, restriction := range desired[matcherType].Branches[branchType].Restrictions {
						var restrictionToAdd *Restriction = nil
						if _, exists := actual[matcherType].Branches[branchType].Restrictions[restrictionType]; !exists {
							restrictionToAdd = restriction
						} else {
							actualExemptions := actual[matcherType].Branches[branchType].Restrictions[restrictionType]
							desiredExemptions := desired[matcherType].Branches[branchType].Restrictions[restrictionType]
							// Dersom unntakene er forskjellige må vi oppdatere
							if !(pkg.SlicesContainsSameElements(desiredExemptions.ExemptUsers, actualExemptions.ExemptUsers) && pkg.SlicesContainsSameElements(desiredExemptions.ExemptGroups, actualExemptions.ExemptGroups)) {
								restrictionToAdd = restriction
							}
						}
						if restrictionToAdd != nil {
							if toUpdate[matcherType] == nil {
								toUpdate[matcherType] = new(Restrictions)
								toUpdate[matcherType].Branches = make(map[string]*BranchRestrictions)
							}
							if toUpdate[matcherType].Branches[branchType] == nil {
								toUpdate[matcherType].Branches[branchType] = new(BranchRestrictions)
								toUpdate[matcherType].Branches[branchType].Restrictions = make(map[string]*Restriction)
							}
							toUpdate[matcherType].Branches[branchType].Restrictions[restrictionType] = restriction
						}
					}
				}
			}
		}
	}

	return toUpdate
}

func findRestrictionsToDelete(desired map[string]*Restrictions, actual map[string]*Restrictions) (toDelete map[string]*Restrictions) {
	toDelete = make(map[string]*Restrictions)

	for matcherType, matcherRestriction := range actual {
		// Sjekk om matcher-typen finnes
		if _, exists := desired[matcherType]; !exists {
			if toDelete[matcherType] == nil {
				toDelete[matcherType] = new(Restrictions)
			}
			// Slett dersom matcher-typen ikke finnes
			toDelete[matcherType] = matcherRestriction
		} else {
			// Hvis matcher-typen finnes må vi sjekke grenene for restriksjoner
			for branchType, branchRestriction := range actual[matcherType].Branches {
				if _, exists := desired[matcherType].Branches[branchType]; !exists {
					if toDelete[matcherType] == nil {
						toDelete[matcherType] = new(Restrictions)
						toDelete[matcherType].Branches = make(map[string]*BranchRestrictions)
					}
					// Sletter dersom grenen ikke skal ha restriksjoner
					toDelete[matcherType].Branches[branchType] = branchRestriction
				} else {
					// Hvis grenen skal ha restriksjoner må vi sjekke hvilke den skal ha
					for restrictionType, restriction := range actual[matcherType].Branches[branchType].Restrictions {
						if _, exists := desired[matcherType].Branches[branchType].Restrictions[restrictionType]; !exists {
							if toDelete[matcherType] == nil {
								toDelete[matcherType] = new(Restrictions)
								toDelete[matcherType].Branches = make(map[string]*BranchRestrictions)
							}
							if toDelete[matcherType].Branches[branchType] == nil {
								toDelete[matcherType].Branches[branchType] = new(BranchRestrictions)
								toDelete[matcherType].Branches[branchType].Restrictions = make(map[string]*Restriction)
							}
							// Dersom grenen ikke skal ha en angitte restriksjonen skal den fjernes
							toDelete[matcherType].Branches[branchType].Restrictions[restrictionType] = restriction
						}
					}
				}
			}
		}
	}

	return toDelete
}

func updateProjectRestriction(baseUrl string, projectKey string, token string, restrictions map[string]*Restrictions) error {
	url := fmt.Sprintf("%s/rest/branch-permissions/latest/projects/%s/restrictions", baseUrl, projectKey)
	return updateRestrictions(url, token, restrictions, "project '"+projectKey+"'")
}

func deleteProjectRestrictions(baseUrl string, projectKey string, token string, restrictions map[string]*Restrictions) error {
	url := fmt.Sprintf("%s/rest/branch-permissions/latest/projects/%s/restrictions/", baseUrl, projectKey)
	return deleteRestriction(url, token, restrictions, "project '"+projectKey+"'")
}

func updateRestrictions(url string, token string, restrictions map[string]*Restrictions, scope string) error {
	for matcherType, matcherRestriction := range restrictions {
		for branchType, branchRestriction := range matcherRestriction.Branches {
			for restrictionType, restriction := range branchRestriction.Restrictions {
				payload, err := json.Marshal(
					&types.CreateRestriction{
						Type: restrictionType,
						Matcher: &types.RestrictionMatcher{
							Id: branchType,
							Type: &types.RestrictionMatcherType{
								Id: matcherType,
							},
						},
						Users:      restriction.ExemptUsers,
						Groups:     restriction.ExemptGroups,
						AccessKeys: nil,
					})
				if err != nil {
					return err
				}
				if _, err := pkg.PostRequest(url, token, bytes.NewReader(payload), nil); err != nil {
					return err
				}
				pterm.Info.Println(pterm.Blue("⚙️ Updated/Created ") + "'" + restrictionType + "' restriction for '" + branchType + "' (" + matcherType + ") in " + scope)
			}
		}
	}
	return nil
}

func deleteRestriction(url string, token string, restrictions map[string]*Restrictions, scope string) error {
	for matcherType, matcherRestriction := range restrictions {
		for branchType, branchRestriction := range matcherRestriction.Branches {
			for restrictionType, restriction := range branchRestriction.Restrictions {
				if _, err := pkg.DeleteRequest(url+strconv.Itoa(restriction.id), token, nil); err != nil {
					return err
				}
				pterm.Info.Println(pterm.Red("🗑️ Deleted ") + "'" + restrictionType + "' restriction for '" + branchType + "' (" + matcherType + ") in " + scope)
			}
		}
	}
	return nil
}
