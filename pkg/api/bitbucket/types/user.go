package types

type User struct {
	Name         string `json:"name,omitempty" yaml:"name"`
	EmailAddress string `json:"emailAddress,omitempty" yaml:"emailAddress"`
	Active       bool   `json:"active" yaml:"active"`
	DisplayName  string `json:"displayName,omitempty" yaml:"displayName"`
	Id           int    `json:"id,omitempty" yaml:"id"`
	Slug         string `json:"slug,omitempty" yaml:"slug"`
	Type         string `json:"type,omitempty" yaml:"type"`
}
