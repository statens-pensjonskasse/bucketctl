package types

type Project struct {
	Id          int    `json:"id,omitempty" yaml:"id,omitempty"`
	Key         string `json:"key,omitempty" yaml:"key,omitempty"`
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Public      bool   `json:"public" yaml:"public"`
	Type        string `json:"type,omitempty" yaml:"type,omitempty"`
	Links       *Links `json:"links,omitempty" yaml:"links,omitempty"`
}

type ProjectsResponse struct {
	response
	Values []*Project `json:"values"`
}
