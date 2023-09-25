package v1alpha1

type TypeMeta struct {
	Kind       string `json:"kind,omitempty" yaml:"kind"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion"`
}

type ObjectMeta struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}
