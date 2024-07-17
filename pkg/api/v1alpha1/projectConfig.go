package v1alpha1

import (
	"git.spk.no/infra/bucketctl/pkg/common"
)

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
	DefaultBranch      *string                 `json:"defaultBranch,omitempty" yaml:"defaultBranch,omitempty"`
	DefaultPermission  *string                 `json:"defaultPermission,omitempty" yaml:"defaultPermission,omitempty"`
	BranchingModel     *BranchingModel         `json:"branchingModel,omitempty" yaml:"branchingModel,omitempty"`
	BranchRestrictions *BranchRestrictions     `json:"branchRestrictions,omitempty" yaml:"branchRestrictions,omitempty"`
	Permissions        *Permissions            `json:"permissions,omitempty" yaml:"permissions,omitempty"`
	Webhooks           *Webhooks               `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`
	Repositories       *RepositoriesProperties `json:"repositories,omitempty" yaml:"repositories,omitempty"`
}

type UncombinedProjectConfigSpecs struct {
	Access             *ProjectConfigSpec
	BranchingModels    *ProjectConfigSpec
	BranchRestrictions *ProjectConfigSpec
	DefaultBranches    *ProjectConfigSpec
	Webhooks           *ProjectConfigSpec
}

type RepositoriesProperties []*RepositoryProperties

type UncombinedRepositoriesProperties struct {
	BranchingModels    *RepositoriesProperties
	BranchRestrictions *RepositoriesProperties
	DefaultBranches    *RepositoriesProperties
	Permissions        *RepositoriesProperties
	Webhooks           *RepositoriesProperties
}

type RepositoryProperties struct {
	RepoSlug           string              `json:"name" yaml:"name"`
	DefaultBranch      *string             `json:"defaultBranch,omitempty" yaml:"defaultBranch,omitempty"`
	BranchingModel     *BranchingModel     `json:"branchingModel,omitempty" yaml:"branchingModel,omitempty"`
	BranchRestrictions *BranchRestrictions `json:"branchRestrictions,omitempty" yaml:"branchRestrictions,omitempty"`
	Permissions        *Permissions        `json:"permissions,omitempty" yaml:"permissions,omitempty"`
	Webhooks           *Webhooks           `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`
}

func EmptyRepositoryProperties(repoSlug string) *RepositoryProperties {
	return &RepositoryProperties{
		RepoSlug:           repoSlug,
		BranchingModel:     &BranchingModel{},
		BranchRestrictions: &BranchRestrictions{},
		Permissions:        &Permissions{},
		Webhooks:           &Webhooks{},
	}
}

func (pcs *ProjectConfigSpec) Copy() *ProjectConfigSpec {
	pcsCopy := &ProjectConfigSpec{
		ProjectKey:        pcs.ProjectKey,
		Public:            pcs.Public,
		DefaultBranch:     pcs.DefaultBranch,
		DefaultPermission: pcs.DefaultPermission,
	}
	if pcs.BranchingModel != nil {
		pcsCopy.BranchingModel = pcs.BranchingModel.Copy()
	}
	if pcs.BranchRestrictions != nil {
		pcsCopy.BranchRestrictions = pcs.BranchRestrictions.Copy()
	}
	if pcs.Permissions != nil {
		pcsCopy.Permissions = pcs.Permissions.Copy()
	}
	if pcs.Webhooks != nil {
		pcsCopy.Webhooks = pcs.Webhooks.Copy()
	}
	if pcs.Repositories != nil {
		pcsCopy.Repositories = pcs.Repositories.Copy()
	}

	return pcsCopy
}

func (rps *RepositoriesProperties) Copy() *RepositoriesProperties {
	rpsCopy := new(RepositoriesProperties)
	for _, rp := range *rps {
		*rpsCopy = append(*rpsCopy, rp.Copy())
	}
	return rpsCopy
}

func (rp *RepositoryProperties) Copy() *RepositoryProperties {
	rpCopy := &RepositoryProperties{
		RepoSlug: rp.RepoSlug,
	}
	if rp.DefaultBranch != nil {
		(*rpCopy).DefaultBranch = (*rp).DefaultBranch
	}
	if rp.BranchingModel != nil {
		rpCopy.BranchingModel = rp.BranchingModel.Copy()
	}
	if rp.BranchRestrictions != nil {
		rpCopy.BranchRestrictions = rp.BranchRestrictions.Copy()
	}
	if rp.Permissions != nil {
		rpCopy.Permissions = rp.Permissions.Copy()
	}
	if rp.Webhooks != nil {
		rpCopy.Webhooks = rp.Webhooks.Copy()
	}

	return rpCopy
}

// Validate TODO: Actually write function
func (pc *ProjectConfig) Validate() error {
	return nil
}

type GroupedRepositories map[string]struct {
	Desired *RepositoryProperties
	Actual  *RepositoryProperties
}

// GroupRepositories Group *RepositoriesProperties by repoSlug
func GroupRepositories(desired *RepositoriesProperties, actual *RepositoriesProperties) GroupedRepositories {
	if desired == nil {
		desired = new(RepositoriesProperties)
	}
	if actual == nil {
		actual = new(RepositoriesProperties)
	}
	grouping := make(GroupedRepositories, len(*desired)+len(*actual))
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

func CombineProjectConfigSpecs(specs *UncombinedProjectConfigSpecs) *ProjectConfigSpec {

	if specs.Access == nil {
		specs.Access = new(ProjectConfigSpec)
	}
	if specs.BranchingModels == nil {
		specs.BranchingModels = new(ProjectConfigSpec)
	}
	if specs.BranchRestrictions == nil {
		specs.BranchRestrictions = new(ProjectConfigSpec)
	}
	if specs.DefaultBranches == nil {
		specs.DefaultBranches = new(ProjectConfigSpec)
	}
	if specs.Webhooks == nil {
		specs.Webhooks = new(ProjectConfigSpec)
	}

	return &ProjectConfigSpec{
		ProjectKey:         specs.Access.ProjectKey,
		Public:             specs.Access.Public,
		DefaultBranch:      specs.DefaultBranches.DefaultBranch,
		DefaultPermission:  specs.Access.DefaultPermission,
		BranchingModel:     specs.BranchingModels.BranchingModel,
		BranchRestrictions: specs.BranchRestrictions.BranchRestrictions,
		Permissions:        specs.Access.Permissions,
		Webhooks:           specs.Webhooks.Webhooks,
		Repositories: CombineRepositoriesProperties(&UncombinedRepositoriesProperties{
			DefaultBranches:    specs.DefaultBranches.Repositories,
			BranchingModels:    specs.BranchingModels.Repositories,
			BranchRestrictions: specs.BranchRestrictions.Repositories,
			Permissions:        specs.Access.Repositories,
			Webhooks:           specs.Webhooks.Repositories,
		}),
	}
}

func CombineRepositoriesProperties(properties *UncombinedRepositoriesProperties) *RepositoriesProperties {
	repositoriesPropertiesMap := make(map[string]*RepositoryProperties)

	// All the different RepositoriesProperties should contain the same repositories
	// We use the first to initialise the map, but we can't be 100% sure that's all the repositories
	if properties.Permissions != nil {
		for _, r := range *properties.Permissions {
			repositoriesPropertiesMap[r.RepoSlug] = &RepositoryProperties{RepoSlug: r.RepoSlug}
			if len(*r.Permissions) > 0 {
				repositoriesPropertiesMap[r.RepoSlug].Permissions = r.Permissions
			}
		}
	}
	if properties.BranchingModels != nil {
		for _, r := range *properties.BranchingModels {
			if r.BranchingModel != nil {
				if repositoriesPropertiesMap[r.RepoSlug] == nil {
					repositoriesPropertiesMap[r.RepoSlug] = &RepositoryProperties{RepoSlug: r.RepoSlug}
				}
				repositoriesPropertiesMap[r.RepoSlug].BranchingModel = r.BranchingModel
			}
		}
	}
	if properties.BranchRestrictions != nil {
		for _, r := range *properties.BranchRestrictions {
			if len(*r.BranchRestrictions) > 0 {
				if repositoriesPropertiesMap[r.RepoSlug] == nil {
					repositoriesPropertiesMap[r.RepoSlug] = &RepositoryProperties{RepoSlug: r.RepoSlug}
				}
				repositoriesPropertiesMap[r.RepoSlug].BranchRestrictions = r.BranchRestrictions
			}
		}
	}
	if properties.DefaultBranches != nil {
		for _, r := range *properties.DefaultBranches {
			if r.DefaultBranch != nil {
				if repositoriesPropertiesMap[r.RepoSlug] == nil {
					repositoriesPropertiesMap[r.RepoSlug] = &RepositoryProperties{RepoSlug: r.RepoSlug}
				}
				repositoriesPropertiesMap[r.RepoSlug].DefaultBranch = r.DefaultBranch
			}
		}
	}
	if properties.Webhooks != nil {
		for _, r := range *properties.Webhooks {
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

	if pcs.DefaultBranch != nil && cmp.DefaultBranch == nil {
		return false
	}
	if pcs.DefaultBranch == nil && cmp.DefaultBranch != nil {
		return false
	}
	if pcs.Public != nil && cmp.DefaultBranch != nil && *pcs.DefaultBranch != *cmp.DefaultBranch {
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

	if pcs.BranchingModel == nil && cmp.BranchingModel != nil {
		return false
	}
	if pcs.BranchingModel != nil && !pcs.BranchingModel.Equals(cmp.BranchingModel) {
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

	if pcs.Repositories == nil && cmp.Repositories != nil {
		return false
	}
	if pcs.Repositories != nil && !pcs.Repositories.Equals(cmp.Repositories) {
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

	if rp.DefaultBranch == nil && cmp.DefaultBranch != nil {
		return false
	}
	if rp.DefaultBranch != nil && cmp.DefaultBranch == nil {
		return false
	}
	if rp.DefaultBranch != nil && cmp.DefaultBranch != nil && *rp.DefaultBranch != *cmp.DefaultBranch {
		return false
	}

	if rp.Permissions == nil && cmp.Permissions != nil {
		return false
	}
	if rp.Permissions != nil && !rp.Permissions.Equals(cmp.Permissions) {
		return false
	}

	if rp.BranchingModel == nil && cmp.BranchingModel != nil {
		return false
	}
	if rp.BranchingModel != nil && !rp.BranchingModel.Equals(cmp.BranchingModel) {
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
