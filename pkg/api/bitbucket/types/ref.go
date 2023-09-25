package types

type Ref struct {
	Id         string      `json:"id,omitempty"`
	DisplayId  string      `json:"displayId,omitempty"`
	Type       string      `json:"type,omitempty"`
	Repository *Repository `json:"repository,omitempty"`
}
