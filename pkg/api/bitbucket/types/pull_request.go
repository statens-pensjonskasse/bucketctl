package types

type PullRequestParticipant struct {
	User   *User    `json:"user,omitempty"`
	Role   string   `json:"role,omitempty"`
	Status string   `json:"status,omitempty"`
	Links  struct{} `json:"links"`
}

type DefaultReviewers struct {
	Id               int      `json:"id,omitempty"`
	Scope            *Scope   `json:"scope,omitempty"`
	SourceRefMatcher *Matcher `json:"sourceRefMatcher,omitempty"`
	TargetRefMatcher *Matcher `json:"targetRefMatcher,omitempty"`
	Reviewers        []*User  `json:"reviewers,omitempty"`
}

type PullRequest struct {
	Title           string                    `json:"title,omitempty"`
	Description     string                    `json:"description,omitempty"`
	State           string                    `json:"state,omitempty"`
	Open            bool                      `json:"open"`
	Closed          bool                      `json:"closed"`
	Locked          bool                      `json:"locked"`
	FromRef         *Ref                      `json:"fromRef,omitempty"`
	ToRef           *Ref                      `json:"toRef,omitempty"`
	Reviewers       []*PullRequestParticipant `json:"reviewers,omitempty"`
	HtmlDescription string                    `json:"htmlDescription,omitempty"`
	Links           *Links
}

type PullRequestInfo struct {
	Id           int                       `json:"id,omitempty"`
	Version      int                       `json:"version,omitempty"`
	Title        string                    `json:"title,omitempty"`
	Description  string                    `json:"description,omitempty"`
	State        string                    `json:"state,omitempty"`
	Open         bool                      `json:"open"`
	Closed       bool                      `json:"closed"`
	Locked       bool                      `json:"locked"`
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
