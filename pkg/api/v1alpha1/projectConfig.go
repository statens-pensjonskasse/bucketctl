package v1alpha1

const (
	ApiVersion        string = "bucketctl.spk.no/v1alpha1"
	ProjectConfigKind string = "ProjectConfig"
)

func ProjectConfigV1alpha1() *ProjectConfig {
	return &ProjectConfig{
		TypeMeta: TypeMeta{
			Kind:       ProjectConfigKind,
			APIVersion: ApiVersion,
		},
		Metadata: ObjectMeta{},
		Spec:     ProjectConfigSpec{},
	}
}

type ProjectConfig struct {
	TypeMeta `json:",inline" yaml:",inline"`
	Metadata ObjectMeta        `json:"metadata" yaml:"metadata"`
	Spec     ProjectConfigSpec `json:"spec" yaml:"spec"`
}

type ProjectConfigSpec struct {
	ProjectKey         string                  `json:"projectKey" yaml:"projectKey"`
	Access             *ProjectAccess          `json:"access,omitempty" yaml:"access,omitempty"`
	BranchRestrictions *BranchRestrictions     `json:"branchRestrictions,omitempty" yaml:"branchRestrictions,omitempty"`
	Webhooks           *Webhooks               `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`
	Repositories       *RepositoriesProperties `json:"repositories,omitempty" yaml:"repositories,omitempty"`
}

type RepositoriesProperties []*RepositoryProperties

type RepositoryProperties struct {
	RepoSlug           string              `json:"name" yaml:"name"`
	Permissions        *Permissions        `json:"permissions,omitempty" yaml:"permissions,omitempty"`
	BranchRestrictions *BranchRestrictions `json:"branchRestrictions,omitempty" yaml:"branchRestrictions,omitempty"`
	Webhooks           *Webhooks           `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`
}

func (a *ProjectConfig) Validate() error {
	return nil
}

type GroupedRepositories map[string]struct {
	Desired *RepositoryProperties
	Actual  *RepositoryProperties
}

// GroupRepositories Group *RepositoriesProperties by repoSlug
func GroupRepositories(desired *RepositoriesProperties, actual *RepositoriesProperties) GroupedRepositories {
	grouping := make(GroupedRepositories)
	if desired == nil {
		desired = new(RepositoriesProperties)
	}
	if actual == nil {
		actual = new(RepositoriesProperties)
	}
	for _, d := range *desired {
		g := grouping[d.RepoSlug]
		g.Desired = d
		grouping[d.RepoSlug] = g
	}
	for _, a := range *actual {
		g := grouping[a.RepoSlug]
		g.Actual = a
		grouping[a.RepoSlug] = g
	}
	return grouping
}

func CombineProjectConfigSpecs(access *ProjectConfigSpec, branchRestrictions *ProjectConfigSpec, webhooks *ProjectConfigSpec) *ProjectConfigSpec {
	return &ProjectConfigSpec{
		ProjectKey:         access.ProjectKey,
		Access:             access.Access,
		BranchRestrictions: branchRestrictions.BranchRestrictions,
		Webhooks:           webhooks.Webhooks,
		Repositories:       CombineRepositoriesProperties(access.Repositories, branchRestrictions.Repositories, webhooks.Repositories),
	}
}

func CombineRepositoriesProperties(permissions *RepositoriesProperties, branchRestrictions *RepositoriesProperties, webhooks *RepositoriesProperties) *RepositoriesProperties {
	repositoriesPropertiesMap := make(map[string]*RepositoryProperties)
	if permissions != nil {
		for _, r := range *permissions {
			g := repositoriesPropertiesMap[r.RepoSlug]
			if g == nil {
				g = new(RepositoryProperties)
			}
			if len(*r.Permissions) > 0 {
				g.RepoSlug = r.RepoSlug
				g.Permissions = r.Permissions
				repositoriesPropertiesMap[r.RepoSlug] = g
			}
		}
	}
	if branchRestrictions != nil {
		for _, r := range *branchRestrictions {
			g := repositoriesPropertiesMap[r.RepoSlug]
			if g == nil {
				g = new(RepositoryProperties)
			}
			if len(*r.BranchRestrictions) > 0 {
				g.RepoSlug = r.RepoSlug
				g.BranchRestrictions = r.BranchRestrictions
				repositoriesPropertiesMap[r.RepoSlug] = g
			}
		}
	}
	if webhooks != nil {
		for _, r := range *webhooks {
			g := repositoriesPropertiesMap[r.RepoSlug]
			if g == nil {
				g = new(RepositoryProperties)
			}
			if len(*r.Webhooks) > 0 {
				g.RepoSlug = r.RepoSlug
				g.Webhooks = r.Webhooks
				repositoriesPropertiesMap[r.RepoSlug] = g
			}
		}
	}

	repositoriesProperties := new(RepositoriesProperties)
	for _, v := range repositoriesPropertiesMap {
		*repositoriesProperties = append(*repositoriesProperties, v)
	}

	return repositoriesProperties
}

func (pcs *ProjectConfigSpec) Equals(cmp *ProjectConfigSpec) bool {
	if pcs == cmp {
		return true
	}
	if cmp == nil {
		return false
	}

	if pcs.ProjectKey != cmp.ProjectKey {
		return false
	}

	if pcs.Access == nil && cmp.Access != nil {
		return false
	}
	if pcs.Access != nil && !pcs.Access.Equals(cmp.Access) {
		return false
	}

	if pcs.BranchRestrictions == nil && cmp.BranchRestrictions != nil {
		return false
	}
	if pcs.BranchRestrictions != nil && !pcs.BranchRestrictions.Equals(cmp.BranchRestrictions) {
		return false
	}

	if pcs.Webhooks == nil && cmp.Webhooks != nil {
		return false
	}
	if pcs.Webhooks != nil && !pcs.Webhooks.Equals(cmp.Webhooks) {
		return false
	}

	return true
}

func (rps *RepositoriesProperties) toMap() map[string]*RepositoryProperties {
	asMap := make(map[string]*RepositoryProperties, len(*rps))
	for _, rp := range *rps {
		asMap[rp.RepoSlug] = rp
	}
	return asMap
}

func (rps *RepositoriesProperties) Equals(cmp *RepositoriesProperties) bool {
	if rps == cmp {
		return true
	}
	if cmp == nil {
		return false
	}
	if len(*rps) != len(*cmp) {
		return false
	}
	rpsMap := rps.toMap()
	cmpMap := cmp.toMap()
	for slug, rp := range rpsMap {
		if !rp.Equals(cmpMap[slug]) {
			return false
		}
	}
	return true
}

func (rp *RepositoryProperties) Equals(cmp *RepositoryProperties) bool {
	if rp == cmp {
		return true
	}
	if cmp == nil {
		return false
	}

	if rp.RepoSlug != cmp.RepoSlug {
		return false
	}

	if rp.Permissions == nil && cmp.Permissions != nil {
		return false
	}
	if rp.Permissions != nil && !rp.Permissions.Equals(cmp.Permissions) {
		return false
	}

	if rp.BranchRestrictions == nil && cmp.BranchRestrictions != nil {
		return false
	}
	if rp.BranchRestrictions != nil && !rp.BranchRestrictions.Equals(cmp.BranchRestrictions) {
		return false
	}

	if rp.Webhooks == nil && cmp.Webhooks != nil {
		return false
	}
	if rp.Webhooks != nil && !rp.Webhooks.Equals(cmp.Webhooks) {
		return false
	}

	return true
}
