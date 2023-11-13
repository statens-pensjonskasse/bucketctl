package v1alpha1

import "bucketctl/pkg/api/bitbucket/types"

type DefaultBranch struct {
	Branch *types.Branch `json:",inline" yaml:",inline"`
}

func FindDefaultBranchesToChange(desired *RepositoriesProperties, actual *RepositoriesProperties) (toUpdate *RepositoriesProperties) {
	toUpdate = new(RepositoriesProperties)

	type candidates struct {
		desiredDefaultBranch *string
		actualDefaultBranch  *string
	}

	defaultBranchMap := make(map[string]*candidates)
	for _, repo := range *desired {
		defaultBranchMap[repo.RepoSlug] = &candidates{desiredDefaultBranch: repo.DefaultBranch}
	}
	for _, repo := range *actual {
		if defaultBranchMap[repo.RepoSlug] == nil {
			continue
		}
		defaultBranchMap[repo.RepoSlug].actualDefaultBranch = repo.DefaultBranch
	}

	for repoSlug, defaultBranches := range defaultBranchMap {
		if defaultBranches.desiredDefaultBranch != defaultBranches.actualDefaultBranch {
			*toUpdate = append(*toUpdate, &RepositoryProperties{RepoSlug: repoSlug, DefaultBranch: defaultBranches.desiredDefaultBranch})
		}
	}

	return toUpdate
}
