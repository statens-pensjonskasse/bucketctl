package types

type Branch struct {
	Id              string  `json:"id,omitempty" yaml:"id,omitempty"`
	RefId           *string `json:"refId,omitempty" yaml:"refId,omitempty"`
	DisplayId       string  `json:"displayId,omitempty" yaml:"displayId,omitempty"`
	Type            string  `json:"type,omitempty" yaml:"type,omitempty"`
	LatestCommit    string  `json:"latestCommit,omitempty" yaml:"latestCommit,omitempty"`
	LatestChangeset string  `json:"latestChangeset,omitempty" yaml:"latestChangeset,omitempty"`
	UseDefault      bool    `json:"useDefault,omitempty" yaml:"useDefault,omitempty"`
}

func (branch *Branch) Copy() *Branch {
	return &Branch{
		Id:              branch.Id,
		RefId:           branch.RefId,
		DisplayId:       branch.DisplayId,
		Type:            branch.Type,
		LatestCommit:    branch.LatestCommit,
		LatestChangeset: branch.LatestChangeset,
		UseDefault:      branch.UseDefault,
	}
}

func (branch *Branch) Equals(cmp *Branch) bool {
	if branch == cmp {
		return true
	}
	if cmp == nil {
		return false
	}
	if branch.Id != cmp.Id {
		return false
	}
	if branch.RefId == nil && cmp.RefId != nil {
		return false
	}
	if branch.RefId != nil && cmp.RefId == nil {
		return false
	}
	if branch.RefId != nil && cmp.RefId != nil && (*branch.RefId != *cmp.RefId) {
		return false
	}
	if branch.DisplayId != cmp.DisplayId {
		return false
	}
	if branch.Type != cmp.Type {
		return false
	}
	if branch.LatestCommit != cmp.LatestCommit {
		return false
	}
	if branch.LatestChangeset != cmp.LatestChangeset {
		return false
	}
	if branch.UseDefault != cmp.UseDefault {
		return false
	}

	return true
}
