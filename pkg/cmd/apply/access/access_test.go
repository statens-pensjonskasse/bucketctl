package access

import (
	. "bucketctl/pkg/api/v1alpha1"
	"bucketctl/pkg/printer"
	"testing"
)

var (
	readPermission = &Permission{
		Name: "READ",
		Entities: &Entities{
			Users: []string{"reader"},
		},
	}
	adminPermission = &Permission{
		Name: "ADMIN",
		Entities: &Entities{
			Users: []string{"admin"},
		},
	}
)

func primitivePtr[T any](val T) *T {
	return &val
}

var (
	noChanges = &ProjectConfigSpec{
		ProjectKey: "A",
		Access: &ProjectAccess{
			Permissions: &Permissions{},
		},
		Repositories: &RepositoriesProperties{},
	}
	read = &ProjectConfigSpec{
		ProjectKey: "A",
		Access: &ProjectAccess{
			Permissions: &Permissions{readPermission},
		},
		Repositories: &RepositoriesProperties{},
	}
	admin = &ProjectConfigSpec{
		ProjectKey: "A",
		Access: &ProjectAccess{
			Permissions: &Permissions{adminPermission},
		},
		Repositories: &RepositoriesProperties{},
	}
	readAndAdmin = &ProjectConfigSpec{
		ProjectKey: "A",
		Access: &ProjectAccess{
			Permissions: &Permissions{readPermission, adminPermission},
		},
		Repositories: &RepositoriesProperties{},
	}
	publicDefaultRead = &ProjectConfigSpec{
		ProjectKey: "A",
		Access: &ProjectAccess{
			Public:            primitivePtr(true),
			DefaultPermission: primitivePtr("PROJECT_READ"),
			Permissions:       &Permissions{},
		},
		Repositories: &RepositoriesProperties{},
	}
	notPublicDefaultWrite = &ProjectConfigSpec{
		ProjectKey: "A",
		Access: &ProjectAccess{
			Public:            primitivePtr(false),
			DefaultPermission: primitivePtr("PROJECT_WRITE"),
			Permissions:       &Permissions{},
		},
		Repositories: &RepositoriesProperties{},
	}
)

func Test_FindAccessChanges(t *testing.T) {
	type args struct {
		desired *ProjectConfigSpec
		actual  *ProjectConfigSpec
	}
	type want struct {
		toCreate *ProjectConfigSpec
		toUpdate *ProjectConfigSpec
		toDelete *ProjectConfigSpec
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no difference",
			args: args{readAndAdmin, readAndAdmin},
			want: want{noChanges, noChanges, noChanges},
		},
		{
			name: "create only",
			args: args{readAndAdmin, read},
			want: want{admin, noChanges, noChanges},
		},
		{
			name: "update only",
			args: args{publicDefaultRead, notPublicDefaultWrite},
			want: want{noChanges, publicDefaultRead, noChanges},
		},
		{
			name: "delete only",
			args: args{noChanges, read},
			want: want{noChanges, noChanges, read},
		},
		{
			name: "create and delete",
			args: args{read, admin},
			want: want{read, noChanges, admin},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToCreate, gotToUpdate, gotToDelete := FindAccessChanges(tt.args.desired, tt.args.actual)
			if !gotToCreate.Equals(tt.want.toCreate) {
				t.Errorf("%s - create got %v, want %v", tt.name, gotToCreate, tt.want.toCreate)
			}
			if !gotToUpdate.Equals(tt.want.toUpdate) {
				t.Errorf("%s - update got %v, want %v", tt.name, gotToUpdate, tt.want.toUpdate)
				printer.PrintData(gotToUpdate, nil)
				printer.PrintData(tt.want.toUpdate, nil)
			}
			if !gotToDelete.Equals(tt.want.toDelete) {
				t.Errorf("%s - delete got %v, want %v", tt.name, gotToDelete, tt.want.toDelete)
			}
		})
	}
}
