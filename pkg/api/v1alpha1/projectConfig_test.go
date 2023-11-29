package v1alpha1

import (
	"bucketctl/pkg/common"
	"gopkg.in/yaml.v3"
	"reflect"
	"testing"
)

var (
	defaultBranch = "defaultBranch"
	rp            = &RepositoryProperties{
		RepoSlug:           "repoSlug",
		DefaultBranch:      &defaultBranch,
		Permissions:        &Permissions{},
		BranchingModel:     &BranchingModel{},
		BranchRestrictions: &BranchRestrictions{},
		Webhooks:           &Webhooks{},
	}
	rps = &RepositoriesProperties{rp}

	public            = true
	defaultPermission = "PROJECT_READ"
	pcs               = &ProjectConfigSpec{
		ProjectKey:         "key",
		Public:             &public,
		DefaultBranch:      &defaultBranch,
		DefaultPermission:  &defaultPermission,
		Permissions:        &Permissions{},
		BranchingModel:     &BranchingModel{},
		BranchRestrictions: &BranchRestrictions{},
		Webhooks:           &Webhooks{},
		Repositories:       &RepositoriesProperties{},
	}
)

func Test_RepositoryProperties_Copy(t *testing.T) {
	tests := []struct {
		name string
		recv RepositoryProperties
		want *RepositoryProperties
	}{
		{
			name: "Copy empty",
			recv: RepositoryProperties{},
			want: &RepositoryProperties{},
		},
		{
			name: "Copy",
			recv: *rp,
			want: rp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.recv.Copy(); !got.Equals(tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}

	rpCopy := rp.Copy()

	rpCopy.RepoSlug = ""
	rpCopy.DefaultBranch = nil
	rpCopy.Permissions = nil
	rpCopy.BranchingModel = nil
	rpCopy.BranchRestrictions = nil
	rpCopy.Webhooks = nil

	if rpCopy.RepoSlug == rp.RepoSlug {
		t.Error("repoSlug not copied")
	}
	if rpCopy.DefaultBranch == rp.DefaultBranch {
		t.Error("defaultBranch not copied")
	}
	if rpCopy.Permissions == rp.Permissions {
		t.Error("permissions not copied")
	}
	if rpCopy.BranchingModel == rp.BranchingModel {
		t.Error("branchingModel not copied")
	}
	if rpCopy.BranchRestrictions == rp.BranchRestrictions {
		t.Error("branchRestrictions not copied")
	}
	if rpCopy.Webhooks == rp.Webhooks {
		t.Error("webhook not copied")
	}
}

func TestRepositoriesProperties_Copy(t *testing.T) {
	tests := []struct {
		name string
		recv RepositoriesProperties
		want *RepositoriesProperties
	}{
		{
			name: "Copy empty",
			recv: RepositoriesProperties{},
			want: &RepositoriesProperties{},
		},
		{
			name: "Copy",
			recv: *rps,
			want: rps,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.recv.Copy(); !got.Equals(tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}

	rpsCopy := rps.Copy()

	(*rpsCopy)[0].RepoSlug = ""
	if rpsCopy.Equals(rps) {
		t.Error("repositories not copied")
	}

}

func TestProjectConfigSpec_Copy(t *testing.T) {
	tests := []struct {
		name string
		recv ProjectConfigSpec
		want *ProjectConfigSpec
	}{
		{
			name: "Copy empty",
			recv: ProjectConfigSpec{},
			want: &ProjectConfigSpec{},
		},
		{
			name: "Copy",
			recv: *pcs,
			want: pcs,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.recv.Copy(); !got.Equals(tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}

	pcsCopy := pcs.Copy()

	pcs.ProjectKey = ""
	pcs.Public = nil
	pcs.DefaultBranch = nil
	pcs.DefaultPermission = nil
	pcs.Permissions = nil
	pcs.BranchingModel = nil
	pcs.BranchRestrictions = nil
	pcs.Webhooks = nil
	pcs.Repositories = nil

	if pcsCopy.ProjectKey == pcs.ProjectKey {
		t.Error("projectKey not copied")
	}
	if pcsCopy.Public == pcs.Public {
		t.Error("public not copied")
	}
	if pcsCopy.DefaultBranch == pcs.DefaultBranch {
		t.Error("defaultBranch not copied")
	}
	if pcsCopy.DefaultPermission == pcs.DefaultPermission {
		t.Error("defaultPermission not copied")
	}
	if pcsCopy.Permissions == pcs.Permissions {
		t.Error("permission not copied")
	}
	if pcsCopy.BranchingModel == pcs.BranchingModel {
		t.Error("branchingModel not copied")
	}
	if pcsCopy.BranchRestrictions == pcs.BranchRestrictions {
		t.Error("branchRestricitons not copied")
	}
	if pcsCopy.Webhooks == pcs.Webhooks {
		t.Error("webhooks not copied")
	}
	if pcsCopy.Repositories == pcs.Repositories {
		t.Error("repositories not copied")
	}
}

func TestProjectConfigSpec_Copy_Equals_Deep(t *testing.T) {

	file := "../../../testdata/cmd/apply/projectConfig/integration/desired.yaml"

	var projectConfig ProjectConfig
	if err := common.ReadConfigFile(file, &projectConfig); err != nil {
		t.Errorf("error reading file %s", file)
	}

	pcsOrig := &projectConfig.Spec
	pcsCopy := pcsOrig.Copy()

	// Compare with hard-coded Equals. If this fails there's an error in Copy() and/or Equals().
	if !pcsCopy.Equals(pcsOrig) {
		t.Error("not equal")
	}

	// Compare with DeepEqual. If this fails there is most probably something wrong with Equals(). Maybe a missing field.
	if !reflect.DeepEqual(pcsCopy, pcsOrig) {
		t.Error("not deeply")
	}

	// Marshal data to YAML and back again
	var pcsDeepCopy *ProjectConfigSpec
	yamlData, err := yaml.Marshal(pcsOrig)
	if err != nil {
		t.Error(err)
	}
	err = yaml.Unmarshal(yamlData, &pcsDeepCopy)
	if err != nil {
		t.Error(err)
	}

	// Compare with DeepCopy. If this fails there is most probably something wrong with Copy(). Possibly a missing field.
	if !reflect.DeepEqual(pcsCopy, pcsDeepCopy) {
		t.Error("not equal")
	}
}
