package types

type Repository struct {
	Id            int      `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	Slug          string   `json:"slug,omitempty"`
	HierarchyId   string   `json:"hierarchyId,omitempty"`
	ScmId         string   `json:"scmId,omitempty"`
	State         string   `json:"state,omitempty"`
	StatusMessage string   `json:"statusMessage,omitempty"`
	Forkable      bool     `json:"forkable"`
	Public        bool     `json:"public"`
	Archived      bool     `json:"archived"`
	Project       *Project `json:"project,omitempty"`
}

type RepositoriesResponse struct {
	response
	Values []*Repository `json:"values"`
}
