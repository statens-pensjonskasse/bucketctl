package common

import (
	types2 "git.spk.no/infra/bucketctl/pkg/api/bitbucket/types"
	"testing"
)

const (
	mainBranch        = "refs/heads/main"
	developBranch     = "refs/heads/develop"
	featureBranch     = "refs/heads/feature/the-best-ever-feature"
	hotfixBranch      = "refs/heads/hotfix/not-so-good-feature-after-all"
	featureMainBranch = "refs/heads/feature/main"
)

var (
	branchModel = &types2.BranchingModel{
		Production: &types2.Branch{
			Id:         "refs/heads/main",
			DisplayId:  "main",
			Type:       "BRANCH",
			UseDefault: false,
		},
		Development: &types2.Branch{
			Id:         "refs/heads/develop",
			DisplayId:  "develop",
			Type:       "BRANCH",
			UseDefault: true,
		},
		Types: []*types2.BranchType{
			{
				Id:          "FEATURE",
				DisplayName: "Feature",
				Prefix:      "feature/",
			},
			{
				Id:          "HOTFIX",
				DisplayName: "Hotfix",
				Prefix:      "hotfix/",
			},
		},
	}
)

func Test_branchMatcher(t *testing.T) {
	mainBranchMatcher := &types2.Matcher{
		Id:        "refs/heads/main",
		DisplayID: "main",
		Active:    true,
		Type: &types2.MatcherType{
			Id:   "BRANCH",
			Name: "Branch",
		},
	}

	if !branchMatcher(mainBranchMatcher, mainBranch) {
		t.Errorf("Expected to match %s", mainBranch)
	}
	if branchMatcher(mainBranchMatcher, developBranch) {
		t.Errorf("Not expected to match %s", developBranch)
	}
	if branchMatcher(mainBranchMatcher, featureMainBranch) {
		t.Errorf("Not expected to match %s", developBranch)
	}
}

func Test_modelBranchMatcher(t *testing.T) {
	modelBranchProductionMatcher := &types2.Matcher{
		Id:        "production",
		DisplayID: "Production",
		Active:    true,
		Type: &types2.MatcherType{
			Id:   "MODEL_BRANCH",
			Name: "Branching model branch",
		},
	}

	if !modelBranchMatcher(modelBranchProductionMatcher, branchModel, mainBranch) {
		t.Errorf("Expected to match production branch")
	}
	if modelBranchMatcher(modelBranchProductionMatcher, branchModel, developBranch) {
		t.Errorf("Not expected to match development branch")
	}
	if modelBranchMatcher(modelBranchProductionMatcher, branchModel, featureMainBranch) {
		t.Errorf("Not expected to match feature branch")
	}
}

func Test_modelCategoryMatcher(t *testing.T) {
	modelCategoryFeatureMatcher := &types2.Matcher{
		Id:        "FEATURE",
		DisplayID: "Feature",
		Active:    true,
		Type: &types2.MatcherType{
			Id:   "MODEL_CATEGORY",
			Name: "Branching model category",
		},
	}

	if !modelCategoryMatcher(modelCategoryFeatureMatcher, branchModel, featureBranch) {
		t.Errorf("Expected to match feature branch %s", featureBranch)
	}
	if !modelCategoryMatcher(modelCategoryFeatureMatcher, branchModel, featureMainBranch) {
		t.Errorf("Expected to match feature branch %s", featureMainBranch)
	}
	if modelCategoryMatcher(modelCategoryFeatureMatcher, branchModel, hotfixBranch) {
		t.Errorf("Not expected to match hotfix branch %s", hotfixBranch)
	}
	if modelCategoryMatcher(modelCategoryFeatureMatcher, branchModel, mainBranch) {
		t.Errorf("Not expected ot match main branch %s", mainBranch)
	}
}

func Test_patternMatcher(t *testing.T) {
	patternType := &types2.MatcherType{
		Id:   "PATTERN",
		Name: "pattern",
	}

	patternFeatureMatcher := &types2.Matcher{
		Id:        "feature/",
		DisplayID: "feature/",
		Active:    true,
		Type:      patternType,
	}
	refsMatchingPattern(t, patternFeatureMatcher, featureBranch, featureMainBranch)
	refsNotMatchingPattern(t, patternFeatureMatcher, mainBranch, developBranch, hotfixBranch)

	// Tests from https://confluence.atlassian.com/bitbucketserver088/branch-permission-patterns-1216582116.html
	patternAnyMatcher := &types2.Matcher{
		Id:        "*",
		DisplayID: "*",
		Active:    true,
		Type:      patternType,
	}
	refsMatchingPattern(t, patternAnyMatcher, mainBranch, developBranch, featureBranch, featureMainBranch, "🍌")

	patternPROJECTMatcher := &types2.Matcher{
		Id:        "PROJECT-*",
		DisplayID: "PROJECT-*",
		Active:    true,
		Type:      patternType,
	}
	refsMatchingPattern(t, patternPROJECTMatcher, "refs/heads/PROJECT-1234", "refs/heads/stable/PROJECT-new", "refs/tags/PROJECT-1.1")
	refsNotMatchingPattern(t, patternPROJECTMatcher, "refs/heads/permission-1234", "refs/heads/stable/PROJECT_NEW", "refs/tags/PROJECT/1.1")

	patternQuestionMarkMatcher := &types2.Matcher{
		Id:        "?.?",
		DisplayID: "?.?",
		Active:    true,
		Type:      patternType,
	}
	refsMatchingPattern(t, patternQuestionMarkMatcher, "refs/heads/1.1", "refs/heads/stable/2.X", "refs/tags/3.1")
	refsNotMatchingPattern(t, patternQuestionMarkMatcher, "refs/heads/1-1")

	patternTagsMatcher := &types2.Matcher{
		Id:        "tags/",
		DisplayID: "tags/",
		Active:    true,
		Type:      patternType,
	}
	refsMatchingPattern(t, patternTagsMatcher, "refs/heads/stable/tags/some_branch", "refs/tags/permission-1.1.0")

	patternTagsStarMatcher := &types2.Matcher{
		Id:        "tags/**",
		DisplayID: "tags/**",
		Active:    true,
		Type:      patternType,
	}
	refsMatchingPattern(t, patternTagsStarMatcher, "refs/heads/stable/tags/some_branch", "refs/tags/permission-1.1.0")

	patternMasterBranchesMatcher := &types2.Matcher{
		Id:        "heads/**/master",
		DisplayID: "heads/**/master",
		Active:    true,
		Type:      patternType,
	}
	refsMatchingPattern(t, patternMasterBranchesMatcher, "refs/heads/master", "refs/heads/stable/master")
	refsNotMatchingPattern(t, patternMasterBranchesMatcher, "refs/tags/master", "refs/heads/stable/master/horse")
}

func refsMatchingPattern(t *testing.T, matcher *types2.Matcher, refs ...string) {
	for _, ref := range refs {
		if !patternMatcher(matcher, ref) {
			t.Errorf("Expected pattern %s to match %s", matcher.Id, ref)
		}
	}
}

func refsNotMatchingPattern(t *testing.T, matcher *types2.Matcher, refs ...string) {
	for _, ref := range refs {
		if patternMatcher(matcher, ref) {
			t.Errorf("Expected pattern %s to not match %s", matcher.Id, ref)
		}
	}
}
