package v1alpha1

func FindDefaultBranchesToChange(desired *RepositoriesProperties, actual *RepositoriesProperties) (toUpdate *RepositoriesProperties) {
	toUpdate = new(RepositoriesProperties)

	type candidates struct {
		desiredDefaultBranch string
		actualDefaultBranch  string
	}

	defaultBranchMap := make(map[string]*candidates)
	for _, repo := range *desired {
		if repo.DefaultBranch != nil {
			defaultBranchMap[repo.RepoSlug] = &candidates{desiredDefaultBranch: *repo.DefaultBranch}
		}
	}
	for _, repo := range *actual {
		if defaultBranchMap[repo.RepoSlug] == nil {
			continue
		}
		if repo.DefaultBranch != nil {
			defaultBranchMap[repo.RepoSlug].actualDefaultBranch = *repo.DefaultBranch
		}
	}

	for repoSlug, defaultBranches := range defaultBranchMap {
		if defaultBranches.desiredDefaultBranch != defaultBranches.actualDefaultBranch {
			*toUpdate = append(*toUpdate, &RepositoryProperties{RepoSlug: repoSlug, DefaultBranch: &defaultBranches.desiredDefaultBranch})
		}
	}

	return toUpdate
}
