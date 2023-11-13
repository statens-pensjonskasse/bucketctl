package types

type BranchingModel struct {
	Development *Branch       `json:"development" yaml:"development"`
	Production  *Branch       `json:"production,omitempty" yaml:"production,omitempty"`
	Types       []*BranchType `json:"types,omitempty" yaml:"types,omitempty"`
	Scope       *Scope        `json:"scope,omitempty" yaml:"scope,omitempty"`
}

type BranchType struct {
	Id          string `json:"id,omitempty" yaml:"id,omitempty"`
	DisplayName string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Enabled     bool   `json:"enabled" yaml:"enabled,omitempty"`
	Prefix      string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
}
