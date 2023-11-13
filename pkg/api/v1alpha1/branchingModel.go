package v1alpha1

import (
	"bucketctl/pkg/api/bitbucket/types"
	"strings"
)

type BranchingModel struct {
	Development *types.Branch        `json:"development,omitempty" yaml:"development,omitempty"`
	Production  *types.Branch        `json:"production,omitempty" yaml:"production,omitempty"`
	Types       *BranchingModelTypes `json:"types,omitempty" yaml:"types,omitempty"`
}

type BranchingModelTypes []*BranchingModelType

type BranchingModelType struct {
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	Prefix string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
}

func FromBitbucketBranchModelTypes(bitbucketBranchTypes []*types.BranchType) *BranchingModelTypes {
	activeBranchModelTypes := make(BranchingModelTypes, 0)
	for _, t := range bitbucketBranchTypes {
		if t.Enabled {
			activeBranchModelTypes = append(activeBranchModelTypes, &BranchingModelType{
				Name:   strings.ToLower(t.Id),
				Prefix: t.Prefix,
			})
		}
	}
	return &activeBranchModelTypes
}

func FindBranchingModelsToChange(desired *BranchingModel, actual *BranchingModel) (toCreate *BranchingModel, toUpdate *BranchingModel, toDelete *BranchingModel) {
	if desired == nil {
		desired = new(BranchingModel)
	}
	if actual == nil {
		actual = new(BranchingModel)
	}

	if actual.IsEmpty() && !desired.IsEmpty() {
		toCreate = desired
		return toCreate, toUpdate, toDelete
	} else if desired.IsEmpty() && !actual.IsEmpty() {
		toDelete = actual
		return toCreate, toUpdate, toDelete
	} else if !desired.Equals(actual) {
		toUpdate = desired
		return toCreate, toUpdate, toDelete
	}

	return toCreate, toUpdate, toDelete
}

func (bm *BranchingModel) IsEmpty() bool {
	if bm.Development == nil && bm.Production == nil && (bm.Types == nil || len(*bm.Types) <= 0) {
		return true
	}
	return false
}

func (bm *BranchingModel) Equals(cmp *BranchingModel) bool {
	if bm == cmp {
		return true
	}
	if cmp == nil {
		return false
	}

	if bm.Development == nil && cmp.Development != nil {
		return false
	}
	if bm.Development != nil && !bm.Development.Equals(cmp.Development) {
		return false
	}

	if bm.Production == nil && cmp.Production != nil {
		return false
	}
	if bm.Production != nil && !bm.Production.Equals(cmp.Production) {
		return false
	}

	if bm.Types == nil && cmp.Types != nil {
		return false
	}
	if bm.Types != nil && !bm.Types.Equals(cmp.Types) {
		return false
	}

	return true
}

func (bmTypes *BranchingModelTypes) Equals(cmp *BranchingModelTypes) bool {
	if bmTypes == cmp {
		return true
	}
	if cmp == nil {
		return false
	}
	if len(*bmTypes) != len(*cmp) {
		return false
	}
	bmTypesMap := make(map[string]*BranchingModelType, len(*bmTypes))
	for _, t := range *bmTypes {
		bmTypesMap[t.Name] = t
	}
	cmpTypesMap := make(map[string]*BranchingModelType, len(*cmp))
	for _, t := range *cmp {
		cmpTypesMap[t.Name] = t
	}
	for name, t := range bmTypesMap {
		if !t.Equals(cmpTypesMap[name]) {
			return false
		}
	}
	return true
}

func (bt *BranchingModelType) Equals(cmp *BranchingModelType) bool {
	if bt == cmp {
		return true
	}
	if cmp == nil {
		return false
	}
	if bt.Name != cmp.Name {
		return false
	}
	if bt.Prefix != cmp.Prefix {
		return false
	}
	return true
}
