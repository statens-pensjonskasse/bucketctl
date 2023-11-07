package types

type Restriction struct {
	Id         int            `json:"id" yaml:"id"`
	Type       string         `json:"type" yaml:"type"`
	Scope      *Scope         `json:"scope" yaml:"scope"`
	Matcher    *Matcher       `json:"matcher" yaml:"matcher"`
	Users      []*User        `json:"users" yaml:"users"`
	Groups     []string       `json:"groups" yaml:"groups"`
	AccessKeys []*interface{} `json:"accessKeys" yaml:"accessKeys"`
}

type CreateRestriction struct {
	Type       string   `json:"type,omitempty" yaml:"type"`
	Matcher    *Matcher `json:"matcher,omitempty" yaml:"matcher"`
	Users      []string `json:"users,omitempty" yaml:"users"`
	Groups     []string `json:"groups,omitempty" yaml:"groups"`
	AccessKeys []string `json:"accessKeys" yaml:"accessKeys"`
}

type RestrictionResponse struct {
	response
	Values []*Restriction `json:"values"`
}

type Matcher struct {
	Id        string       `json:"id" yaml:"id"`
	DisplayID string       `json:"displayID" yaml:"displayID"`
	Active    bool         `json:"active" yaml:"active"`
	Type      *MatcherType `json:"type" yaml:"type"`
}

type MatcherType struct {
	Id   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}
