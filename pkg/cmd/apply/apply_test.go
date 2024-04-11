package apply

import (
	. "bucketctl/pkg/api/v1alpha1"
	"testing"
)

func Test_readProjectConfig(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "file not found",
			args:    args{file: "this-file-does-not-exist-🙈.yaml"},
			wantErr: true,
		},
		{
			name:    "read successfully",
			args:    args{file: "../../../testdata/cmd/apply/projectConfig/integration/desired.yaml"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readProjectConfig(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("readProjectConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got.Kind != ProjectConfigKind {
				t.Errorf("Got Kind %s, expected %s", got.Kind, ProjectConfigKind)
			}
			if got.APIVersion != ApiVersion {
				t.Errorf("Got ApiVersion %s, expected %s", got.APIVersion, ApiVersion)
			}
			if got.Metadata.Name != "desired_project_config" {
				t.Errorf("Got metadata.mame %s, expected %s", got.Metadata.Name, "desired_project_config")
			}
			if got.Spec.ProjectKey != "TEST" {
				t.Errorf("Got spec.projectKey %s, expected %s", got.Spec.ProjectKey, "TEST")
			}
		})
	}
}

func Test_findProjectConfigChanges(t *testing.T) {
	type args struct {
		desired *ProjectConfigSpec
		actual  *ProjectConfigSpec
	}
	type want struct {
		toCreate *ProjectConfigSpec
		toUpdate *ProjectConfigSpec
		toDelete *ProjectConfigSpec
	}

	desired, err := readProjectConfig("../../../testdata/cmd/apply/projectConfig/integration/desired.yaml")
	if err != nil {
		t.Errorf("Error reading testdata")
	}
	actual, err := readProjectConfig("../../../testdata/cmd/apply/projectConfig/integration/actual.yaml")
	if err != nil {
		t.Errorf("Error reading testdata")
	}
	toCreate, err := readProjectConfig("../../../testdata/cmd/apply/projectConfig/integration/toCreate.yaml")
	if err != nil {
		t.Errorf("Error reading testdata")
	}
	toUpdate, err := readProjectConfig("../../../testdata/cmd/apply/projectConfig/integration/toUpdate.yaml")
	if err != nil {
		t.Errorf("Error reading testdata")
	}
	toDelete, err := readProjectConfig("../../../testdata/cmd/apply/projectConfig/integration/toDelete.yaml")
	if err != nil {
		t.Errorf("Error reading testdata")
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "integration test",
			args: args{&desired.Spec, &actual.Spec},
			want: want{&toCreate.Spec, &toUpdate.Spec, &toDelete.Spec},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToCreate, gotToUpdate, gotToDelete := findProjectConfigChanges(tt.args.desired, tt.args.actual)
			if !gotToCreate.Equals(tt.want.toCreate) {
				t.Errorf("findProjectConfigChanges() gotToCreate = %v, want %v", gotToCreate, tt.want.toCreate)
			}
			if !gotToUpdate.Equals(tt.want.toUpdate) {
				t.Errorf("findProjectConfigChanges() gotToUpdate = %v, want %v", gotToUpdate, tt.want.toUpdate)
			}
			if !gotToDelete.Equals(tt.want.toDelete) {
				t.Errorf("findProjectConfigChanges() gotToDelete = %v, want %v", gotToDelete, tt.want.toDelete)
			}
		})
	}
}

func Test_removeBranchRestriction(t *testing.T) {
	type args struct {
		desired *ProjectConfigSpec
		actual  *ProjectConfigSpec
	}
	type want struct {
		toCreate *ProjectConfigSpec
		toUpdate *ProjectConfigSpec
		toDelete *ProjectConfigSpec
	}

	desired, err := readProjectConfig("../../../testdata/cmd/apply/projectConfig/removeBranchRestriction/desired.yaml")
	if err != nil {
		t.Errorf("Error reading testdata")
	}
	actual, err := readProjectConfig("../../../testdata/cmd/apply/projectConfig/removeBranchRestriction/actual.yaml")
	if err != nil {
		t.Errorf("Error reading testdata")
	}
	toDelete, err := readProjectConfig("../../../testdata/cmd/apply/projectConfig/removeBranchRestriction/toDelete.yaml")
	if err != nil {
		t.Errorf("Error reading testdata")
	}

	noChange := &ProjectConfigSpec{
		ProjectKey:         "INFRA",
		Permissions:        new(Permissions),
		BranchRestrictions: new(BranchRestrictions),
		Webhooks:           new(Webhooks),
		Repositories:       new(RepositoriesProperties),
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "integration test",
			args: args{&desired.Spec, &actual.Spec},
			want: want{noChange, noChange, &toDelete.Spec},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToCreate, gotToUpdate, gotToDelete := findProjectConfigChanges(tt.args.desired, tt.args.actual)
			if !gotToCreate.Equals(tt.want.toCreate) {
				t.Errorf("findProjectConfigChanges() gotToCreate = %v, want %v", gotToCreate, tt.want.toCreate)
			}
			if !gotToUpdate.Equals(tt.want.toUpdate) {
				t.Errorf("findProjectConfigChanges() gotToUpdate = %v, want %v", gotToUpdate, tt.want.toUpdate)
			}
			if !gotToDelete.Equals(tt.want.toDelete) {
				t.Errorf("findProjectConfigChanges() gotToDelete = %v, want %v", gotToDelete, tt.want.toDelete)
			}
		})
	}
}

func Test_caseInsensitivePermissions(t *testing.T) {
	type args struct {
		desired *ProjectConfigSpec
		actual  *ProjectConfigSpec
	}
	type want struct {
		toCreate *ProjectConfigSpec
		toUpdate *ProjectConfigSpec
		toDelete *ProjectConfigSpec
	}

	desired, err := readProjectConfig("../../../testdata/cmd/apply/projectConfig/caseInsensitiveNamesInPermissions/desired.yaml")
	if err != nil {
		t.Errorf("Error reading testdata")
	}
	actual, err := readProjectConfig("../../../testdata/cmd/apply/projectConfig/caseInsensitiveNamesInPermissions/actual.yaml")
	if err != nil {
		t.Errorf("Error reading testdata")
	}

	noChange := &ProjectConfigSpec{
		ProjectKey:         "TEST",
		Permissions:        new(Permissions),
		BranchRestrictions: new(BranchRestrictions),
		Webhooks:           new(Webhooks),
		Repositories:       new(RepositoriesProperties),
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "integration test",
			args: args{&desired.Spec, &actual.Spec},
			want: want{noChange, noChange, noChange},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToCreate, gotToUpdate, gotToDelete := findProjectConfigChanges(tt.args.desired, tt.args.actual)
			if !gotToCreate.Equals(tt.want.toCreate) {
				t.Errorf("findProjectConfigChanges() gotToCreate = %v, want %v", gotToCreate, tt.want.toCreate)
			}
			if !gotToUpdate.Equals(tt.want.toUpdate) {
				t.Errorf("findProjectConfigChanges() gotToUpdate = %v, want %v", gotToUpdate, tt.want.toUpdate)
			}
			if !gotToDelete.Equals(tt.want.toDelete) {
				t.Errorf("findProjectConfigChanges() gotToDelete = %v, want %v", gotToDelete, tt.want.toDelete)
			}
		})
	}
}
