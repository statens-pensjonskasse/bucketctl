package types

import (
	"reflect"
)

type response struct {
	Size          int    `json:"size"`
	Limit         int    `json:types.LimitFlag`
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
	Key         string `json:types.ProjectKeyFlag`
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

type DefaultProjectPermission struct {
	Permitted bool `json:"permitted" yaml:"permitted"`
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

type Restriction struct {
	Id         int                 `json:"id" yaml:"id"`
	Type       string              `json:"type" yaml:"type"`
	Scope      *RestrictionScope   `json:"scope" yaml:"scope"`
	Matcher    *RestrictionMatcher `json:"matcher" yaml:"matcher"`
	Users      []*User             `json:"users" yaml:"users"`
	Groups     []string            `json:"groups" yaml:"groups"`
	AccessKeys []*interface{}      `json:"accessKeys" yaml:"accessKeys"`
}

type CreateRestriction struct {
	Type       string              `json:"type,omitempty" yaml:"type"`
	Matcher    *RestrictionMatcher `json:"matcher,omitempty" yaml:"matcher"`
	Users      []string            `json:"users,omitempty" yaml:"users"`
	Groups     []string            `json:"groups,omitempty" yaml:"groups"`
	AccessKeys []string            `json:"accessKeys" yaml:"accessKeys"`
}

type RestrictionScope struct {
	ResourceId int    `json:"resourceId" yaml:"resourceId"`
	Type       string `json:"type" yaml:"type"`
}

type RestrictionMatcher struct {
	Id        string                  `json:"id" yaml:"id"`
	DisplayID string                  `json:"displayID" yaml:"displayID"`
	Active    bool                    `json:"active" yaml:"active"`
	Type      *RestrictionMatcherType `json:"type" yaml:"type"`
}

type RestrictionMatcherType struct {
	Id   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

type RestrictionResponse struct {
	response
	Values []*Restriction `json:"values"`
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

func (base *Webhook) Copy() *Webhook {
	copied := &Webhook{
		Id:                      base.Id,
		Name:                    base.Name,
		CreatedDate:             base.CreatedDate,
		UpdatedDate:             base.UpdatedDate,
		Configuration:           base.Configuration,
		Url:                     base.Url,
		Active:                  base.Active,
		ScopeType:               base.ScopeType,
		SslVerificationRequired: base.SslVerificationRequired,
	}
	copied.Events = append(copied.Events, base.Events...)

	return copied
}

func (base *Webhook) Equivalent(candidate *Webhook) bool {
	if candidate == nil {
		return false
	}
	if base.Name != candidate.Name {
		return false
	}
	if base.Url != candidate.Url {
		return false
	}
	if base.Active != candidate.Active {
		return false
	}
	if base.SslVerificationRequired != candidate.SslVerificationRequired {
		return false
	}
	if !reflect.DeepEqual(base.Configuration, candidate.Configuration) {
		return false
	}
	if !reflect.DeepEqual(base.Events, candidate.Events) {
		return false
	}
	return true
}

// Similarity Finner gir en score på hvor like to webhooks er mellom 0.0 og 1.0
// Dersom ID er lik antas webhookene å være de samme
func (base *Webhook) Similarity(candidate *Webhook) float64 {
	if candidate == nil {
		return 0.0
	}
	if base.Id == candidate.Id {
		return 1.0
	}
	similarityScore := 0.0
	if base.Name == candidate.Name {
		similarityScore += 0.3
	}
	if base.Url == candidate.Url {
		similarityScore += 0.1
	}
	if base.Active == candidate.Active {
		similarityScore += 0.1
	}
	if base.ScopeType == candidate.ScopeType {
		similarityScore += 0.1
	}
	if base.SslVerificationRequired == candidate.SslVerificationRequired {
		similarityScore += 0.1
	}
	if reflect.DeepEqual(base.Configuration, candidate.Configuration) {
		similarityScore += 0.1
	}
	if len(base.Events) == len(candidate.Events) {
		similarityScore += 0.1
	}
	return similarityScore
}
