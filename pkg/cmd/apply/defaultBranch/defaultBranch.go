package defaultBranch

import . "bucketctl/pkg/api/v1alpha1"

func FindDefaultBranchChanges(desired *ProjectConfigSpec, actual *ProjectConfigSpec) (toUpdate *ProjectConfigSpec) {
	toUpdate = new(ProjectConfigSpec)
	return toUpdate
}
