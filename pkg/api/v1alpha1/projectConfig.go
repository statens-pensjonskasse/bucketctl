package v1alpha1

import "bucketctl/pkg/common"

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
	Public             *bool                   `json:"public,omitempty" yaml:"public,omitempty"`
	DefaultPermission  *string                 `json:"defaultPermission,omitempty" yaml:"defaultPermission,omitempty"`
	Permissions        *Permissions            `json:"permissions,omitempty" yaml:"permissions,omitempty"`
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

func EmptyRepositoryProperties(repoSlug string) *RepositoryProperties {
	return &RepositoryProperties{
		RepoSlug:           repoSlug,
		Permissions:        &Permissions{},
		BranchRestrictions: &BranchRestrictions{},
		Webhooks:           &Webhooks{},
	}
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
	grouping := make(GroupedRepositories, len(*desired)+len(*actual))
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

	// Initialise other half of group if nil with empty properties
	for repoSlug, g := range grouping {
		if g.Actual == nil {
			g.Actual = EmptyRepositoryProperties(repoSlug)
		}
		if g.Desired == nil {
			g.Desired = EmptyRepositoryProperties(repoSlug)
		}
		grouping[repoSlug] = g
	}

	return grouping
}

func CombineProjectConfigSpecs(access *ProjectConfigSpec, branchRestrictions *ProjectConfigSpec, webhooks *ProjectConfigSpec) *ProjectConfigSpec {
	return &ProjectConfigSpec{
		ProjectKey:         access.ProjectKey,
		Public:             access.Public,
		DefaultPermission:  access.DefaultPermission,
		Permissions:        access.Permissions,
		BranchRestrictions: branchRestrictions.BranchRestrictions,
		Webhooks:           webhooks.Webhooks,
		Repositories:       CombineRepositoriesProperties(access.Repositories, branchRestrictions.Repositories, webhooks.Repositories),
	}
}

func CombineRepositoriesProperties(
	repoPermissions *RepositoriesProperties,
	repoBranchRestrictions *RepositoriesProperties,
	repoWebhooks *RepositoriesProperties) *RepositoriesProperties {
	repositoriesPropertiesMap := make(map[string]*RepositoryProperties, len(*repoPermissions))

	// All the different RepositoriesProperties should contain the same repositories
	// We use the first to initialise the map, but we can't be 100% sure that's all the repositories
	if repoPermissions != nil {
		for _, r := range *repoPermissions {
			repositoriesPropertiesMap[r.RepoSlug] = &RepositoryProperties{RepoSlug: r.RepoSlug}
			if len(*r.Permissions) > 0 {
				repositoriesPropertiesMap[r.RepoSlug].Permissions = r.Permissions
			}
		}
	}
	if repoBranchRestrictions != nil {
		for _, r := range *repoBranchRestrictions {
			if len(*r.BranchRestrictions) > 0 {
				if repositoriesPropertiesMap[r.RepoSlug] == nil {
					repositoriesPropertiesMap[r.RepoSlug] = &RepositoryProperties{RepoSlug: r.RepoSlug}
				}
				repositoriesPropertiesMap[r.RepoSlug].BranchRestrictions = r.BranchRestrictions
			}
		}
	}
	if repoWebhooks != nil {
		for _, r := range *repoWebhooks {
			if len(*r.Webhooks) > 0 {
				if repositoriesPropertiesMap[r.RepoSlug] == nil {
					repositoriesPropertiesMap[r.RepoSlug] = &RepositoryProperties{RepoSlug: r.RepoSlug}
				}
				repositoriesPropertiesMap[r.RepoSlug].Webhooks = r.Webhooks
			}
		}
	}

	repositoriesProperties := new(RepositoriesProperties)
	for _, repoSlug := range common.GetLexicallySortedKeys(repositoriesPropertiesMap) {
		*repositoriesProperties = append(*repositoriesProperties, repositoriesPropertiesMap[repoSlug])
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

	if pcs.Public != nil && cmp.Public == nil {
		return false
	}
	if pcs.Public == nil && cmp.Public != nil {
		return false
	}
	if pcs.Public != nil && cmp.Public != nil && *pcs.Public != *cmp.Public {
		return false
	}
	if pcs.DefaultPermission != nil && cmp.DefaultPermission == nil {
		return false
	}
	if pcs.DefaultPermission == nil && cmp.DefaultPermission != nil {
		return false
	}
	if pcs.DefaultPermission != nil && cmp.DefaultPermission != nil && *pcs.DefaultPermission != *cmp.DefaultPermission {
		return false
	}

	if pcs.Permissions == nil && cmp.Permissions != nil {
		return false
	}
	if pcs.Permissions != nil && !pcs.Permissions.Equals(cmp.Permissions) {
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
