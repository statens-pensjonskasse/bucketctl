package v1alpha1

import "git.spk.no/infra/bucketctl/pkg/api/bitbucket/types"

type Branch struct {
	Id              string  `json:"id,omitempty" yaml:"id,omitempty"`
	RefId           *string `json:"refId,omitempty" yaml:"refId,omitempty"`
	DisplayId       string  `json:"displayId,omitempty" yaml:"displayId,omitempty"`
	Type            string  `json:"type,omitempty" yaml:"type,omitempty"`
	LatestCommit    string  `json:"latestCommit,omitempty" yaml:"latestCommit,omitempty"`
	LatestChangeset string  `json:"latestChangeset,omitempty" yaml:"latestChangeset,omitempty"`
	UseDefault      bool    `json:"useDefault,omitempty" yaml:"useDefault,omitempty"`
}

type Branches []*Branch

func FromBitbucketBranch(bitbucketBranch *types.Branch) *Branch {
	return &Branch{
		Id:              bitbucketBranch.Id,
		RefId:           bitbucketBranch.RefId,
		DisplayId:       bitbucketBranch.DisplayId,
		Type:            bitbucketBranch.Type,
		LatestCommit:    bitbucketBranch.LatestCommit,
		LatestChangeset: bitbucketBranch.LatestChangeset,
		UseDefault:      bitbucketBranch.UseDefault,
	}
}

func (b *Branches) ContainsBranchId(branchId *string) bool {
	for _, branch := range *b {
		if branch.Id == *branchId {
			return true
		}
	}
	return false
}

func (b *Branches) ContainsBranchDisplayId(branchDisplayId *string) bool {
	for _, branch := range *b {
		if branch.DisplayId == *branchDisplayId {
			return true
		}
	}
	return false
}
