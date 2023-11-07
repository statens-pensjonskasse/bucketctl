package v1alpha1

import (
	"bucketctl/pkg/api/bitbucket/types"
	"strings"
)

type BranchModel struct {
	DefaultBranch *string            `json:"defaultBranch,omitempty" yaml:"defaultBranch,omitempty"`
	Development   *types.Branch      `json:"development,omitempty" yaml:"development,omitempty"`
	Production    *types.Branch      `json:"production,omitempty" yaml:"production,omitempty"`
	Types         []*BranchModelType `json:"types,omitempty" yaml:"types,omitempty"`
}

type BranchModelType struct {
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	Prefix string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
}

func FromBitbucketBranchModelTypes(bitbucketBranchTypes []*types.BranchType) []*BranchModelType {
	activeBranchModelTypes := make([]*BranchModelType, 0)
	for _, t := range bitbucketBranchTypes {
		if t.Enabled {
			activeBranchModelTypes = append(activeBranchModelTypes, &BranchModelType{
				Name:   strings.ToLower(t.Id),
				Prefix: t.Prefix,
			})
		}
	}
	return activeBranchModelTypes
}
