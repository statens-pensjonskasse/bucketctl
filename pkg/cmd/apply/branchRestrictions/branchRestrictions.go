package branchRestrictions

import (
	"bucketctl/pkg/api/bitbucket/types"
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/common"
	"bytes"
	"encoding/json"
	"github.com/pterm/pterm"
	"strconv"
	"strings"
)

func FindBranchRestrictionChanges(desired *ProjectConfigSpec, actual *ProjectConfigSpec) (toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) {
	projBrToCreate, projBrToUpdate, projBrToDelete := FindBranchRestrictionsToChange(desired.BranchRestrictions, actual.BranchRestrictions)
	repoBrToCreate, repoBrToUpdate, repoBrToDelete := findRepositoriesBranchRestrictionsChanges(desired.Repositories, actual.Repositories)

	toCreate = &ProjectConfigSpec{ProjectKey: desired.ProjectKey, BranchRestrictions: projBrToCreate, Repositories: repoBrToCreate}
	toUpdate = &ProjectConfigSpec{ProjectKey: desired.ProjectKey, BranchRestrictions: projBrToUpdate, Repositories: repoBrToUpdate}
	toDelete = &ProjectConfigSpec{ProjectKey: desired.ProjectKey, BranchRestrictions: projBrToDelete, Repositories: repoBrToDelete}

	return toCreate, toUpdate, toDelete
}

func SetBranchRestrictions(baseUrl string, projectKey string, token string, toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) error {
	if err := setProjectBranchRestrictions(baseUrl, projectKey, token, toCreate.BranchRestrictions, toUpdate.BranchRestrictions, toDelete.BranchRestrictions); err != nil {
		return err
	}

	if err := setRepositoriesBranchRestrictions(baseUrl, projectKey, token, toCreate.Repositories, toUpdate.Repositories, toDelete.Repositories); err != nil {
		return err
	}

	return nil
}

func GetChangesAsText(toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) (changes []string) {
	changes = append(changes, changesToText(pterm.Green("🪵 create"), toCreate)...)
	changes = append(changes, changesToText(pterm.Blue("🌳 update"), toUpdate)...)
	changes = append(changes, changesToText(pterm.Red("🔥 delete"), toDelete)...)
	return changes
}

func changesToText(action string, pcs *ProjectConfigSpec) (changes []string) {
	if pcs.BranchRestrictions != nil && len(*pcs.BranchRestrictions) > 0 {
		for _, br := range *pcs.BranchRestrictions {
			for _, bm := range *br.BranchMatchers {
				for _, r := range *bm.Restrictions {
					changes = append(changes,
						pterm.Sprintf("%s restriction %s (%s) %s (%s) in project %s",
							action,
							br.Type, bm.Matching,
							r.Type, formatRestrictionExemptions(r),
							pcs.ProjectKey))
				}
			}
		}
	}
	if pcs.Repositories != nil && len(*pcs.Repositories) > 0 {
		for _, repo := range *pcs.Repositories {
			if repo.BranchRestrictions != nil {
				for _, br := range *repo.BranchRestrictions {
					for _, bm := range *br.BranchMatchers {
						for _, r := range *bm.Restrictions {
							changes = append(changes,
								pterm.Sprintf("%s restriction %s (%s) %s (%s) in repository %s/%s",
									action,
									br.Type, bm.Matching,
									r.Type, formatRestrictionExemptions(r),
									pcs.ProjectKey, repo.RepoSlug))
						}
					}
				}
			}
		}
	}
	return changes
}

func formatRestrictionExemptions(r *Restriction) string {
	var exemptions = ""
	if r.ExemptUsers != nil && len(r.ExemptUsers) > 0 {
		exemptions += "users: "
		for _, u := range r.ExemptUsers {
			exemptions += u + ", "
		}
	}
	if r.ExemptGroups != nil && len(r.ExemptGroups) > 0 {
		exemptions += "groups: "
		for _, g := range r.ExemptGroups {
			exemptions += g + ", "
		}
	}
	exemptions = strings.TrimSuffix(exemptions, ", ")
	if len(exemptions) == 0 {
		exemptions += "NONE"
	}

	return "Exemptions: " + exemptions
}

func createBranchRestrictions(url string, token string, branchRestrictions *BranchRestrictions, action string, scope string) error {
	if branchRestrictions != nil && len(*branchRestrictions) > 0 {
		for _, br := range *branchRestrictions {
			for _, bm := range *br.BranchMatchers {
				for _, r := range *bm.Restrictions {
					payload, err := json.Marshal(
						&types.CreateRestriction{
							Type: r.Type,
							Matcher: &types.Matcher{
								Id: bm.Matching,
								Type: &types.MatcherType{
									Id: br.Type,
								},
							},
							Users:      r.ExemptUsers,
							Groups:     r.ExemptGroups,
							AccessKeys: nil,
						})
					if err != nil {
						return err
					}
					if _, err := common.PostRequest(url, token, bytes.NewReader(payload), nil); err != nil {
						return err
					}
					pterm.Printfln("%s %s restriction for %s (%s) in %s", action, r.Type, bm.Matching, br.Type, scope)
				}
			}
		}
	}
	return nil
}

func deleteRestrictions(url string, token string, branchRestrictions *BranchRestrictions, scope string) error {
	if branchRestrictions != nil && len(*branchRestrictions) > 0 {
		for _, matcherRestriction := range *branchRestrictions {
			for _, branchRestriction := range *matcherRestriction.BranchMatchers {
				for _, restriction := range *branchRestriction.Restrictions {
					if _, err := common.DeleteRequest(url+"/"+strconv.Itoa(restriction.Id), token, nil); err != nil {
						return err
					}
					pterm.Printfln("%s %s restriction for %s (%s) in %s", pterm.Red("🗑️ Deleted"), restriction.Type, branchRestriction.Matching, matcherRestriction.Type, scope)
				}
			}
		}
	}
	return nil
}
