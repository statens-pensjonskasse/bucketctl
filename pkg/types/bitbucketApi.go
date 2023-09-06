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

type Commit struct {
	Id                 string    `json:"id,omitempty"`
	DisplayId          string    `json:"displayId,omitempty"`
	Author             *User     `json:"author,omitempty"`
	AuthorTimestamp    int       `json:"authorTimestamp,omitempty"`
	Committer          *User     `json:"committer,omitempty"`
	CommitterTimestamp int       `json:"committerTimestamp,omitempty"`
	Message            string    `json:"message,omitempty"`
	Parents            []*Commit `json:"parents,omitempty"`
}

type CommitsResponse struct {
	response
	Values []*Commit `json:"values,omitempty"`
}

type Ref struct {
	Id         string      `json:"id,omitempty"`
	DisplayId  string      `json:"displayId,omitempty"`
	Type       string      `json:"type,omitempty"`
	Repository *Repository `json:"repository,omitempty"`
}

type Repository struct {
	Id            int      `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	Slug          string   `json:"slug,omitempty"`
	HierarchyId   string   `json:"hierarchyId,omitempty"`
	ScmId         string   `json:"scmId,omitempty"`
	State         string   `json:"state,omitempty"`
	StatusMessage string   `json:"statusMessage,omitempty"`
	Forkable      bool     `json:"forkable,omitempty"`
	Public        bool     `json:"public,omitempty"`
	Archived      bool     `json:"archived,omitempty"`
	Project       *Project `json:"project,omitempty"`
}

type RepositoriesResponse struct {
	response
	Values []*Repository `json:"values"`
}

type PullRequestParticipant struct {
	User            *User    `json:"user,omitempty"`
	Role            string   `json:"role,omitempty"`
	Status          string   `json:"status,omitempty"`
	HtmlDescription string   `json:"htmlDescription,omitempty"`
	Links           struct{} `json:"links"`
}

type DefaultReviewers struct {
	Id               int      `json:"id,omitempty"`
	Scope            *Scope   `json:"scope,omitempty"`
	SourceRefMatcher *Matcher `json:"sourceRefMatcher,omitempty"`
	TargetRefMatcher *Matcher `json:"targetRefMatcher,omitempty"`
	Reviewers        []*User  `json:"reviewers,omitempty"`
}

type Project struct {
	Id          int    `json:"id,omitempty"`
	Key         string `json:"key,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Public      bool   `json:"public,omitempty"`
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
	Name         string `json:"name,omitempty" yaml:"name"`
	EmailAddress string `json:"emailAddress,omitempty" yaml:"emailAddress"`
	Active       bool   `json:"active,omitempty" yaml:"active"`
	DisplayName  string `json:"displayName,omitempty" yaml:"displayName"`
	Id           int    `json:"id,omitempty" yaml:"id"`
	Slug         string `json:"slug,omitempty" yaml:"slug"`
	Type         string `json:"type,omitempty" yaml:"type"`
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

type Scope struct {
	ResourceId int    `json:"resourceId" yaml:"resourceId"`
	Type       string `json:"type" yaml:"type"`
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

type RestrictionResponse struct {
	response
	Values []*Restriction `json:"values"`
}

type BranchModel struct {
	Production  *Branch       `json:"production"`
	Development *Branch       `json:"development"`
	Types       []*BranchType `json:"types,omitempty"`
}

type Branch struct {
	Id              string `json:"id,omitempty"`
	DisplayId       string `json:"displayId,omitempty"`
	Type            string `json:"type,omitempty"`
	LatestCommit    string `json:"latestCommit,omitempty"`
	LatestChangeset string `json:"latestChangeset,omitempty"`
	IsDefault       bool   `json:"isDefault,omitempty"`
}

type BranchType struct {
	Id          string `json:"id,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Prefix      string `json:"prefix,omitempty"`
}

type PullRequest struct {
	Title       string                    `json:"title,omitempty"`
	Description string                    `json:"description,omitempty"`
	State       string                    `json:"state,omitempty"`
	Open        bool                      `json:"open,omitempty"`
	Closed      bool                      `json:"closed,omitempty"`
	Locked      bool                      `json:"locked,omitempty"`
	FromRef     *Ref                      `json:"fromRef,omitempty"`
	ToRef       *Ref                      `json:"toRef,omitempty"`
	Reviewers   []*PullRequestParticipant `json:"reviewers,omitempty"`
}

type PullRequestInfo struct {
	Id           int                       `json:"id,omitempty"`
	Version      int                       `json:"version,omitempty"`
	Title        string                    `json:"title,omitempty"`
	Description  string                    `json:"description,omitempty"`
	State        string                    `json:"state,omitempty"`
	Open         bool                      `json:"open,omitempty"`
	Closed       bool                      `json:"closed,omitempty"`
	Locked       bool                      `json:"locked,omitempty"`
	CreatedDate  int                       `json:"createdDate,omitempty"`
	UpdatedDate  int                       `json:"updatedDate,omitempty"`
	Author       *PullRequestParticipant   `json:"author,omitempty"`
	FromRef      *Ref                      `json:"fromRef,omitempty"`
	ToRef        *Ref                      `json:"toRef,omitempty"`
	Reviewers    []*PullRequestParticipant `json:"reviewers,omitempty"`
	Participants []*PullRequestParticipant `json:"participants,omitempty"`
	Links        *Links                    `json:"links,omitempty"`
}

type PullRequestProperties struct {
	MergeResult       *MergeResult `json:"mergeResult,omitempty"`
	QgStatus          string       `json:"qgStatus,omitempty"`
	ResolvedTaskCount int          `json:"resolvedTaskCount,omitempty"`
	CommentCount      int          `json:"commentCount,omitempty"`
	OpenTaskCount     int          `json:"openTaskCount,omitempty"`
}

type MergeResult struct {
	Outcome string `json:"outcome,omitempty"`
	Current bool   `json:"current,omitempty"`
}

type Links struct {
	Self  []*Href `json:"self,omitempty"`
	Clone []*Href `json:"clone,omitempty"`
}

type Href struct {
	Href string `json:"href,omitempty"`
	Name string `json:"name,omitempty"`
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
	if len(base.Events) != len(candidate.Events) {
		return false
	}
	elements := make(map[string]struct{}, len(base.Events))
	// Create a map with all the (unique) elements of list A as keys
	for _, v := range base.Events {
		elements[v] = struct{}{}
	}
	// Check that all the elements of list B has a key in the map
	for _, v := range candidate.Events {
		if _, exists := elements[v]; !exists {
			return false
		}
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
