package types

type response struct {
	Size          int    `json:"size"`
	Limit         int    `json:"limit"`
	IsLastPage    bool   `json:"isLastPage"`
	Start         int    `json:"start"`
	NextPageStart int    `json:"nextPageStart"`
	Values        []byte `json:"values"`
}

type Links struct {
	Self  []*Href `json:"self,omitempty" yaml:"self,omitempty"`
	Clone []*Href `json:"clone,omitempty" yaml:"clone,omitempty"`
}

type Href struct {
	Href string `json:"href,omitempty" yaml:"href,omitempty"`
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

type Scope struct {
	ResourceId int    `json:"resourceId" yaml:"resourceId"`
	Type       string `json:"type" yaml:"type"`
}
