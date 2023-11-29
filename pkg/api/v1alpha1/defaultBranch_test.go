package v1alpha1

import (
	"testing"
)

func Test_FindDefaultBranchesToChange(t *testing.T) {
	main := "main"
	master := "master"
	develop := "develop"

	testRepoMain := &RepositoryProperties{RepoSlug: "test", DefaultBranch: &main}
	testRepoMaster := &RepositoryProperties{RepoSlug: "test", DefaultBranch: &master}

	type args struct {
		desired *ProjectConfigSpec
		actual  *ProjectConfigSpec
	}
	tests := []struct {
		name   string
		args   args
		wanted *ProjectConfigSpec
	}{
		{
			name: "no changes",
			args: args{
				&ProjectConfigSpec{Repositories: &RepositoriesProperties{testRepoMain}},
				&ProjectConfigSpec{Repositories: &RepositoriesProperties{testRepoMain}}},
			wanted: &ProjectConfigSpec{Repositories: &RepositoriesProperties{}},
		},
		{
			name: "different",
			args: args{
				&ProjectConfigSpec{Repositories: &RepositoriesProperties{testRepoMain}},
				&ProjectConfigSpec{Repositories: &RepositoriesProperties{testRepoMaster}}},
			wanted: &ProjectConfigSpec{Repositories: &RepositoriesProperties{testRepoMain}},
		},
		{
			name: "missing actual",
			args: args{
				&ProjectConfigSpec{Repositories: &RepositoriesProperties{testRepoMain}},
				&ProjectConfigSpec{Repositories: &RepositoriesProperties{}}},
			wanted: &ProjectConfigSpec{Repositories: &RepositoriesProperties{testRepoMain}},
		},
		{
			name: "missing desired",
			args: args{
				&ProjectConfigSpec{Repositories: &RepositoriesProperties{}},
				&ProjectConfigSpec{Repositories: &RepositoriesProperties{testRepoMaster}}},
			wanted: &ProjectConfigSpec{Repositories: &RepositoriesProperties{}},
		},
		{
			name: "nothing",
			args: args{
				&ProjectConfigSpec{Repositories: &RepositoriesProperties{}},
				&ProjectConfigSpec{Repositories: &RepositoriesProperties{}}},
			wanted: &ProjectConfigSpec{Repositories: &RepositoriesProperties{}},
		},
		{
			name: "project defaultBranch",
			args: args{
				&ProjectConfigSpec{DefaultBranch: &main, Repositories: &RepositoriesProperties{{RepoSlug: "repo1"}}},
				&ProjectConfigSpec{DefaultBranch: &master, Repositories: &RepositoriesProperties{{RepoSlug: "repo1"}}}},
			wanted: &ProjectConfigSpec{DefaultBranch: &main, Repositories: &RepositoriesProperties{{RepoSlug: "repo1", DefaultBranch: &main}}},
		},
		{
			name: "repo defaultBranch override project defaultBranch",
			args: args{
				&ProjectConfigSpec{DefaultBranch: &main, Repositories: &RepositoriesProperties{{RepoSlug: "repo1", DefaultBranch: &master}}},
				&ProjectConfigSpec{Repositories: &RepositoriesProperties{{RepoSlug: "repo1", DefaultBranch: &main}}}},
			wanted: &ProjectConfigSpec{DefaultBranch: &main, Repositories: &RepositoriesProperties{{RepoSlug: "repo1", DefaultBranch: &master}}},
		},
		{
			name: "repo defaultBranch override project defaultBranch",
			args: args{
				&ProjectConfigSpec{
					DefaultBranch: &main,
					Repositories: &RepositoriesProperties{
						{RepoSlug: "repo0", DefaultBranch: &develop},
						{RepoSlug: "repo1", DefaultBranch: &master},
						{RepoSlug: "repo2", DefaultBranch: nil},
						{RepoSlug: "repo3", DefaultBranch: &main},
					}},
				&ProjectConfigSpec{
					DefaultBranch: &master,
					Repositories: &RepositoriesProperties{
						{RepoSlug: "repo0", DefaultBranch: &develop},
						{RepoSlug: "repo1", DefaultBranch: &main},
						{RepoSlug: "repo2", DefaultBranch: nil},
						{RepoSlug: "repo3", DefaultBranch: &main},
					}}},
			wanted: &ProjectConfigSpec{
				DefaultBranch: &main,
				Repositories: &RepositoriesProperties{
					{RepoSlug: "repo1", DefaultBranch: &master},
					{RepoSlug: "repo2", DefaultBranch: &main},
				}},
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
