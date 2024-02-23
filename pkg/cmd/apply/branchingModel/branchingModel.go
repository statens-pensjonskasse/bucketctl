package branchingModel

import (
	"bucketctl/pkg/api/bitbucket/types"
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/common"
	"bucketctl/pkg/logger"
	"bytes"
	"encoding/json"
	"github.com/pterm/pterm"
	"strings"
)

func FindBranchingModelChanges(desired *ProjectConfigSpec, actual *ProjectConfigSpec) (toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) {
	// All projects need a branching model, can't create a new or delete, only update
	_, projBmToUpdate, _ := FindBranchingModelsToChange(desired.BranchingModel, actual.BranchingModel)
	repoBmToCreate, repoBmToUpdate, repoBmToDelete := findRepositoryBranchingModelChanges(desired.Repositories, actual.Repositories)

	toCreate = &ProjectConfigSpec{ProjectKey: desired.ProjectKey, BranchingModel: nil, Repositories: repoBmToCreate}
	toUpdate = &ProjectConfigSpec{ProjectKey: desired.ProjectKey, BranchingModel: projBmToUpdate, Repositories: repoBmToUpdate}
	toDelete = &ProjectConfigSpec{ProjectKey: desired.ProjectKey, BranchingModel: nil, Repositories: repoBmToDelete}

	return toCreate, toUpdate, toDelete
}

func SetBranchingModels(baseUrl string, projectKey string, token string, toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) error {
	if err := setProjectBranchingModel(baseUrl, projectKey, token, toUpdate.BranchingModel); err != nil {
		return nil
	}

	if err := setRepositoriesBranchingModels(baseUrl, projectKey, token, toCreate.Repositories, toUpdate.Repositories, toDelete.Repositories); err != nil {
		return err
	}

	return nil
}

func GetChangesAsText(toCreate *ProjectConfigSpec, toUpdate *ProjectConfigSpec, toDelete *ProjectConfigSpec) (changes []string) {
	changes = append(changes, changesToText(pterm.Green("🌱 create"), toCreate)...)
	changes = append(changes, changesToText(pterm.Blue("🪴 update"), toUpdate)...)
	changes = append(changes, changesToText(pterm.Red("🍂 delete"), toDelete)...)
	return changes
}

func changesToText(action string, pcs *ProjectConfigSpec) (changes []string) {
	if pcs.BranchingModel != nil && !pcs.BranchingModel.IsEmpty() {
		changes = append(changes,
			pterm.Sprintf("%s branching model (dev: %s, prod: %s) in project %s",
				action,
				branchingModelBranchAsString(pcs.BranchingModel.Development),
				branchingModelBranchAsString(pcs.BranchingModel.Production), pcs.ProjectKey))
	}
	if pcs.Repositories != nil && len(*pcs.Repositories) > 0 {
		for _, repo := range *pcs.Repositories {
			if repo.BranchingModel != nil && !repo.BranchingModel.IsEmpty() {
				changes = append(changes,
					pterm.Sprintf("%s branching model (dev: %s, prod: %s) in repository %s/%s",
						action,
						branchingModelBranchAsString(repo.BranchingModel.Development),
						branchingModelBranchAsString(repo.BranchingModel.Production),
						pcs.ProjectKey, repo.RepoSlug))
			}
		}
	}
	return changes
}

func createBranchingModel(url string, token string, branchingModel *BranchingModel, action string, scope string) error {
	if branchingModel != nil && !branchingModel.IsEmpty() {
		developmentBranch := &types.Branch{
			RefId:      branchingModel.Development.RefId,
			UseDefault: branchingModel.Development.UseDefault,
		}
		productionBranch := &types.Branch{}
		if branchingModel.Production != nil {
			productionBranch.RefId = branchingModel.Production.RefId
			productionBranch.UseDefault = branchingModel.Production.UseDefault
		}
		var branchTypes []*types.BranchType
		for _, t := range *branchingModel.Types {
			branchTypes = append(branchTypes, &types.BranchType{
				Id:          t.Name,
				DisplayName: t.Name,
				Enabled:     true,
				Prefix:      t.Prefix,
			})
		}
		scopeType := &types.Scope{}
		if strings.Contains(scope, "/") {
			scopeType.Type = "REPOSITORY"
		} else {
			scopeType.Type = "PROJECT"
		}

		payload, err := json.Marshal(
			&types.BranchingModel{
				Development: developmentBranch,
				Production:  productionBranch,
				Types:       branchTypes,
				Scope:       scopeType,
			},
		)
		if err != nil {
			return err
		}
		if _, err := common.PutRequest(url, token, bytes.NewReader(payload), nil); err != nil {
			return err
		}

		logger.Log("%s branching model (dev: %s, prod: %s) for %s", action, branchingModelBranchAsString(developmentBranch), branchingModelBranchAsString(productionBranch), scope)
	}
	return nil
}

func branchingModelBranchAsString(branch *types.Branch) string {
	branchingModelBranch := "N/A"
	if branch == nil {
		return branchingModelBranch
	}
	if branch.RefId != nil {
		branchingModelBranch = *branch.RefId
	}
	// If useDefault is true it overrides refId
	if branch.UseDefault {
		branchingModelBranch = "useDefault"
	}
	return branchingModelBranch
}

func deleteBranchingModel(url string, token string, branchingModel *BranchingModel, scope string) error {
	if branchingModel != nil && !branchingModel.IsEmpty() {
		if _, err := common.DeleteRequest(url, token, nil); err != nil {
			return err
		}
		logger.Log("%s branching model for %s", pterm.Red("🍂 Deleted"), scope)
	}
	return nil
}
