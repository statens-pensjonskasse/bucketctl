package common

import (
	"bucketctl/pkg/types"
	"github.com/vibrantbyte/go-antpath/antpath"
	"strings"
)

const (
	AnyRef        = "ANY_REF"
	Branch        = "BRANCH"
	Pattern       = "PATTERN"
	ModelBranch   = "MODEL_BRANCH"
	ModelCategory = "MODEL_CATEGORY"
)

func RefMatcher(matcher *types.Matcher, branchModel *types.BranchModel, ref string) bool {
	matcherTypeId := matcher.Type.Id

	switch matcherTypeId {
	case AnyRef:
		return true
	case Branch:
		return branchMatcher(matcher, ref)
	case Pattern:
		return patternMatcher(matcher, ref)
	case ModelBranch:
		return modelBranchMatcher(matcher, branchModel, ref)
	case ModelCategory:
		return modelCategoryMatcher(matcher, branchModel, ref)
	default:
		return false
	}
}

func branchMatcher(matcher *types.Matcher, ref string) bool {
	return matcher.Id == ref
}

func patternMatcher(matcher *types.Matcher, ref string) bool {
	// https://confluence.atlassian.com/bitbucketserver088/branch-permission-patterns-1216582116.html
	pattern := matcher.Id

	if strings.HasSuffix(pattern, "/") {
		pattern += "**"
	}

	antMatcher := antpath.New()

	// Check each part of the ref separately
	for _, part := range strings.Split(ref, "/") {
		if antMatcher.Match(pattern, part) {
			return true
		}
	}

	// Check the remainder of each ref after removing the preceding part and separator
	remainder := ref
	for _, part := range strings.Split(ref, "/") {
		remainder = strings.TrimPrefix(remainder, part+"/")
		if antMatcher.Match(pattern, remainder) {
			return true
		}
	}

	return false
}

func modelBranchMatcher(matcher *types.Matcher, branchModel *types.BranchModel, ref string) bool {
	switch matcher.Id {
	case "production":
		if ref == branchModel.Production.Id {
			return true
		}
	case "development":
		if ref == branchModel.Development.Id {
			return true
		}
	}
	return false
}

func modelCategoryMatcher(matcher *types.Matcher, branchModel *types.BranchModel, ref string) bool {
	// Branch is fetched in the form "refs/heads/<branch>", e.g. "refs/heads/feature/myAwesomeFeature"
	branch := strings.TrimPrefix(ref, "refs/heads/")

	// Find the prefix of the matcher branch category
	var matcherBranchTypePrefix string
	for _, category := range branchModel.Types {
		if category.Id == matcher.Id {
			matcherBranchTypePrefix = category.Prefix
			break
		}
	}

	if strings.HasPrefix(branch, matcherBranchTypePrefix) {
		return true
	}

	return false
}
