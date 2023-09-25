package types

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

type UserPermission struct {
	User       User   `json:"user"`
	Permission string `json:"permission"`
}

type UserPermissionsResponse struct {
	response
	Values []*UserPermission `json:"values"`
}
