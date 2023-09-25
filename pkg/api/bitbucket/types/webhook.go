package types

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
