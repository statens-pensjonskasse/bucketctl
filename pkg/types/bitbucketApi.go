package types

type response struct {
	Size          int    `json:"size"`
	Limit         int    `json:"limit"`
	IsLastPage    bool   `json:"isLastPage"`
	Start         int    `json:"start"`
	NextPageStart int    `json:"nextPageStart"`
	Values        []byte `json:"values"`
}

type Error struct {
	Errors []struct {
		Context       string `json:"context,omitempty"`
		Message       string `json:"message,omitempty"`
		ExceptionName string `json:"exceptionName,omitempty"`
	} `json:"errors"`
}

type Repository struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Slug          string  `json:"slug"`
	HierarchyId   string  `json:"hierarchyId"`
	ScmId         string  `json:"scmId"`
	State         string  `json:"state"`
	StatusMessage string  `json:"statusMessage"`
	Forkable      bool    `json:"forkable"`
	Public        bool    `json:"public"`
	Archived      bool    `json:"archived"`
	Project       Project `json:"project"`
}

type RepositoriesResponse struct {
	response
	Values []*Repository `json:"values"`
}

type Project struct {
	Id          int    `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type ProjectsResponse struct {
	response
	Values []*Project `json:"values"`
}

type Group struct {
	Name string `json:"name" yaml:"name"`
}

type GroupPermission struct {
	Group      Group  `json:"group"`
	Permission string `json:"permission"`
}

type GroupPermissionsResponse struct {
	response
	Values []*GroupPermission `json:"values"`
}

type User struct {
	Name         string `json:"name" yaml:"name"`
	EmailAddress string `json:"emailAddress" yaml:"emailAddress"`
	Active       bool   `json:"active" yaml:"active"`
	DisplayName  string `json:"displayName" yaml:"displayName"`
	Id           int    `json:"id" yaml:"id"`
	Slug         string `json:"slug" yaml:"slug"`
	Type         string `json:"type" yaml:"type"`
}

type UserPermission struct {
	User       User   `json:"user"`
	Permission string `json:"permission"`
}

type UserPermissionsResponse struct {
	response
	Values []*UserPermission `json:"values"`
}

type Webhook struct {
	Id                      int         `json:"id" yaml:"id"`
	Name                    string      `json:"name" yaml:"name"`
	CreatedDate             int         `json:"createdDate" yaml:"createdDate"`
	UpdatedDate             int         `json:"updatedDate" yaml:"updatedDate"`
	Events                  []string    `json:"events" yaml:"events"`
	Configuration           interface{} `json:"configuration" yaml:"configuration"`
	Url                     string      `json:"url" yaml:"url"`
	Active                  bool        `json:"active" yaml:"active"`
	ScopeType               string      `json:"scopeType" yaml:"scopeType"`
	SslVerificationRequired bool        `yaml:"sslVerificationRequired" yaml:"sslVerificationRequired" yaml:"sslVerificationRequired"`
}

type WebhooksResponse struct {
	response
	Values []*Webhook `json:"values"`
}
