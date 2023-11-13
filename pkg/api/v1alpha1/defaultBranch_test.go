package v1alpha1

import (
	"testing"
)

func Test_FindDefaultBranchesToChange(t *testing.T) {
	main := "main"
	master := "master"

	testRepoMain := &RepositoryProperties{RepoSlug: "test", DefaultBranch: &main}
	testRepoMaster := &RepositoryProperties{RepoSlug: "test", DefaultBranch: &master}

	type args struct {
		desired *RepositoriesProperties
		actual  *RepositoriesProperties
	}
	tests := []struct {
		name   string
		args   args
		wanted *RepositoriesProperties
	}{
		{
			name:   "no changes",
			args:   args{&RepositoriesProperties{testRepoMain}, &RepositoriesProperties{testRepoMain}},
			wanted: &RepositoriesProperties{},
		},
		{
			name:   "different",
			args:   args{&RepositoriesProperties{testRepoMain}, &RepositoriesProperties{testRepoMaster}},
			wanted: &RepositoriesProperties{testRepoMain},
		},
		{
			name:   "missing actual",
			args:   args{&RepositoriesProperties{testRepoMain}, &RepositoriesProperties{}},
			wanted: &RepositoriesProperties{testRepoMain},
		},
		{
			name:   "missing desired",
			args:   args{&RepositoriesProperties{}, &RepositoriesProperties{testRepoMaster}},
			wanted: &RepositoriesProperties{},
		},
		{
			name:   "nothing",
			args:   args{&RepositoriesProperties{}, &RepositoriesProperties{}},
			wanted: &RepositoriesProperties{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotToUpdate := FindDefaultBranchesToChange(tt.args.desired, tt.args.actual); !gotToUpdate.Equals(tt.wanted) {
				t.Errorf("FindDefaultBranchesToChange() = %v, want %v", gotToUpdate, tt.wanted)
			}
		})
	}
}
