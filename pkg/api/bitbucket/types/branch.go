package types

type Branch struct {
	Id              string `json:"id,omitempty"`
	DisplayId       string `json:"displayId,omitempty"`
	Type            string `json:"type,omitempty"`
	LatestCommit    string `json:"latestCommit,omitempty"`
	LatestChangeset string `json:"latestChangeset,omitempty"`
	IsDefault       bool   `json:"isDefault"`
}
