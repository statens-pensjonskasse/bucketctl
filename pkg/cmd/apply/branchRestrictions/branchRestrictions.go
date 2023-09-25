package branchRestrictions

import (
	"bucketctl/pkg/api/bitbucket/types"
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/common"
	"bytes"
	"encoding/json"
	"github.com/pterm/pterm"
	"strconv"
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

func PrintBranchRestrictionChanges(toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) {
	printBranchRestriction(pterm.Green("🪵 create"), toCreate)
	printBranchRestriction(pterm.Blue("🌳 update"), toUpdate)
	printBranchRestriction(pterm.Red("🔥 delete"), toDelete)
}

func printBranchRestriction(action string, pcs *ProjectConfigSpec) {
	if pcs.BranchRestrictions != nil && len(*pcs.BranchRestrictions) > 0 {
		for _, br := range *pcs.BranchRestrictions {
			for _, bm := range *br.BranchMatchers {
				for _, r := range *bm.Restrictions {
					pterm.Printfln("%s %s (%s) %s restriction in project %s",
						action, pterm.Bold.Sprint(bm.Matching), pterm.Bold.Sprint(br.Type), pterm.Bold.Sprint(r.Type), pcs.ProjectKey)
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
							pterm.Printfln("%s %s (%s) %s restriction in repository %s/%s",
								action, pterm.Bold.Sprint(r.Type), pterm.Bold.Sprint(bm.Matching), pterm.Bold.Sprint(br.Type), pcs.ProjectKey, repo.RepoSlug)
						}
					}
				}
			}
		}
	}
}

func createBranchRestrictions(url string, token string, branchRestrictions *BranchRestrictions, action string, scope string) error {
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
				pterm.Info.Printfln("%s %s restriction for %s (%s) in %s", action, r.Type, bm.Matching, br.Type, scope)
			}
		}
	}
	return nil
}

func deleteRestrictions(url string, token string, restrictions *BranchRestrictions, scope string) error {
	for _, matcherRestriction := range *restrictions {
		for _, branchRestriction := range *matcherRestriction.BranchMatchers {
			for _, restriction := range *branchRestriction.Restrictions {
				if _, err := common.DeleteRequest(url+strconv.Itoa(restriction.Id), token, nil); err != nil {
					return err
				}
				pterm.Info.Printfln("%s %s restriction for %s (%s) in %s", pterm.Red("🗑️ Deleted"), restriction.Type, branchRestriction.Matching, matcherRestriction.Type, scope)
			}
		}
	}
	return nil
}
