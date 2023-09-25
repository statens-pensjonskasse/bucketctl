package types

type BranchModel struct {
	Production  *Branch       `json:"production"`
	Development *Branch       `json:"development"`
	Types       []*BranchType `json:"types,omitempty"`
}

type BranchType struct {
	Id          string `json:"id,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Prefix      string `json:"prefix,omitempty"`
}
