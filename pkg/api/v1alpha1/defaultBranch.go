package v1alpha1

func FindDefaultBranchesToChange(desired *ProjectConfigSpec, actual *ProjectConfigSpec) (toUpdate *ProjectConfigSpec) {
	reposToUpdate := new(RepositoriesProperties)

	type candidates struct {
		desiredDefaultBranch string
		actualDefaultBranch  string
	}

	defaultBranchMap := make(map[string]*candidates)
	for _, repo := range *desired.Repositories {
		if repo.DefaultBranch != nil {
			defaultBranchMap[repo.RepoSlug] = &candidates{desiredDefaultBranch: *repo.DefaultBranch}
		} else if desired.DefaultBranch != nil {
			defaultBranchMap[repo.RepoSlug] = &candidates{desiredDefaultBranch: *desired.DefaultBranch}
		}
	}
	for _, repo := range *actual.Repositories {
		if defaultBranchMap[repo.RepoSlug] == nil {
			continue
		}
		if repo.DefaultBranch != nil {
			defaultBranchMap[repo.RepoSlug].actualDefaultBranch = *repo.DefaultBranch
		} else if actual.DefaultBranch != nil {
			defaultBranchMap[repo.RepoSlug].actualDefaultBranch = *actual.DefaultBranch
		}
	}

	for repoSlug, defaultBranches := range defaultBranchMap {
		if defaultBranches.desiredDefaultBranch != defaultBranches.actualDefaultBranch {
			*reposToUpdate = append(*reposToUpdate, &RepositoryProperties{RepoSlug: repoSlug, DefaultBranch: &defaultBranches.desiredDefaultBranch})
		}
	}

	toUpdate = &ProjectConfigSpec{ProjectKey: desired.ProjectKey, DefaultBranch: desired.DefaultBranch, Repositories: reposToUpdate}

	return toUpdate
}
