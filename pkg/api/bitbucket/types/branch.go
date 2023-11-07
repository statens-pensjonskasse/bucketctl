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
