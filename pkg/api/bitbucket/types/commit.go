package types

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
