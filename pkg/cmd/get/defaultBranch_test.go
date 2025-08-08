package get

import (
	"testing"

	. "git.spk.no/infra/bucketctl/pkg/api/v1alpha1"
)

var (
	mainBranch   = "main"
	masterBranch = "master"
)

func Test_getMostCommonDefaultBranch(t *testing.T) {
	tests := []struct {
		name string
		args *RepositoriesProperties
		want string
	}{
		{
			name: "No repositories",
			args: &RepositoriesProperties{},
			want: "N/A",
		},
		{
			name: "One repository",
			args: &RepositoriesProperties{{RepoSlug: "repo1", DefaultBranch: &mainBranch}},
			want: mainBranch,
		},
		{
			name: "Two main, one master",
			args: &RepositoriesProperties{
				{RepoSlug: "repo1", DefaultBranch: &mainBranch},
				{RepoSlug: "repo2", DefaultBranch: &mainBranch},
				{RepoSlug: "repo3", DefaultBranch: &masterBranch}},
			want: mainBranch,
		},
		{
			name: "Two main, two master. Lexical preference.",
			args: &RepositoriesProperties{
				{RepoSlug: "repo1", DefaultBranch: &masterBranch},
				{RepoSlug: "repo2", DefaultBranch: &mainBranch},
				{RepoSlug: "repo3", DefaultBranch: &masterBranch},
				{RepoSlug: "repo4", DefaultBranch: &mainBranch}},
			want: mainBranch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMostCommonDefaultBranch(tt.args); got != tt.want {
				t.Errorf("getMostCommonDefaultBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FilterRepositoryDefaultBranchMatchingProjectDefaultBranch(t *testing.T) {
	tests := []struct {
		name string
		args *ProjectConfigSpec
		want *ProjectConfigSpec
	}{
		{
			name: "No project defaultBranch",
			args: &ProjectConfigSpec{
				DefaultBranch: nil,
				Repositories: &RepositoriesProperties{
					{RepoSlug: "repo1", DefaultBranch: &mainBranch},
					{RepoSlug: "repo2", DefaultBranch: &masterBranch},
				},
			},
			want: &ProjectConfigSpec{
				DefaultBranch: nil,
				Repositories: &RepositoriesProperties{
					{RepoSlug: "repo1", DefaultBranch: &mainBranch},
					{RepoSlug: "repo2", DefaultBranch: &masterBranch},
				},
			},
		},
		{
			name: "main branch project defaultBranch",
			args: &ProjectConfigSpec{
				DefaultBranch: &mainBranch,
				Repositories: &RepositoriesProperties{
					{RepoSlug: "repo1", DefaultBranch: &mainBranch},
					{RepoSlug: "repo2", DefaultBranch: &masterBranch},
				},
			},
			want: &ProjectConfigSpec{
				DefaultBranch: &mainBranch,
				Repositories: &RepositoriesProperties{
					{RepoSlug: "repo1", DefaultBranch: nil},
					{RepoSlug: "repo2", DefaultBranch: &masterBranch},
				},
			},
		},
		{
			name: "project defaultBranch not present in repos",
			args: &ProjectConfigSpec{
				DefaultBranch: &mainBranch,
				Repositories: &RepositoriesProperties{
					{RepoSlug: "repo1", DefaultBranch: &masterBranch},
					{RepoSlug: "repo2", DefaultBranch: &masterBranch},
				},
			},
			want: &ProjectConfigSpec{
				DefaultBranch: &mainBranch,
				Repositories: &RepositoriesProperties{
					{RepoSlug: "repo1", DefaultBranch: &masterBranch},
					{RepoSlug: "repo2", DefaultBranch: &masterBranch},
				},
			},
		},
		{
			name: "project defaultBranch present in all repos",
			args: &ProjectConfigSpec{
				DefaultBranch: &mainBranch,
				Repositories: &RepositoriesProperties{
					{RepoSlug: "repo1", DefaultBranch: &mainBranch},
					{RepoSlug: "repo2", DefaultBranch: &mainBranch},
				},
			},
			want: &ProjectConfigSpec{
				DefaultBranch: &mainBranch,
				Repositories: &RepositoriesProperties{
					{RepoSlug: "repo1", DefaultBranch: nil},
					{RepoSlug: "repo2", DefaultBranch: nil},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterRepositoryDefaultBranchMatchingProjectDefaultBranch(tt.args); !got.Equals(tt.want) {
				t.Errorf("FilterRepositoryDefaultBranchMatchingProjectDefaultBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}
