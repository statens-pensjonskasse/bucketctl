package v1alpha1

import (
	"bucketctl/pkg/api/bitbucket/types"
	"bucketctl/pkg/common"
)

type BranchRestrictions []*BranchRestriction

type BranchRestriction struct {
	Type           string          `json:"type" yaml:"type"`
	BranchMatchers *BranchMatchers `json:"matchers" yaml:"matchers"`
}

type BranchMatchers []*BranchMatcher

type BranchMatcher struct {
	Matching     string        `json:"matching" yaml:"matching"`
	Restrictions *Restrictions `json:"restrictions" yaml:"restrictions"`
}

type Restrictions []*Restriction

type Restriction struct {
	Id           int      `json:"-" yaml:"-"`
	Type         string   `json:"type" yaml:"type"`
	ExemptUsers  []string `json:"exempt-users,omitempty" yaml:"exempt-users,omitempty"`
	ExemptGroups []string `json:"exempt-groups,omitempty" yaml:"exempt-groups,omitempty"`
}

func FindBranchRestrictionsToChange(desired *BranchRestrictions, actual *BranchRestrictions) (toCreate *BranchRestrictions, toUpdate *BranchRestrictions, toDelete *BranchRestrictions) {
	creationCandidates := desired.findBranchRestrictionsDifference(actual).toMap()
	deletionCandidates := actual.findBranchRestrictionsDifference(desired).toMap()

	toUpdate = new(BranchRestrictions)
	for brType, br := range creationCandidates {
		// If the same BranchRestriction type appears in both we might want to update it,
		// but we have to check the restrictions first
		if _, exists := deletionCandidates[brType]; exists {
			toUpdateBr := &BranchRestriction{brType, new(BranchMatchers)}
			creationBmMap := br.BranchMatchers.toMap()
			deletionBmMap := deletionCandidates[brType].BranchMatchers.toMap()
			for bmMatching, bm := range creationBmMap {
				if deletionBm, exists := deletionBmMap[bmMatching]; exists {
					toUpdateBm := &BranchMatcher{bmMatching, new(Restrictions)}
					creationRMap := bm.Restrictions.toMap()
					deletionRMap := deletionBm.Restrictions.toMap()
					for rType, r := range creationRMap {
						if d, exists := deletionRMap[rType]; exists {
							toUpdateR := &Restriction{d.Id, r.Type, r.ExemptUsers, r.ExemptGroups}
							// Restriction was found in both differences,
							// this means it should neither be created nor deleted, but updated
							// We use the deletion candidate ID since we don't supply a desired ID
							delete(creationRMap, rType)
							delete(deletionRMap, rType)
							*toUpdateBm.Restrictions = append(*toUpdateBm.Restrictions, toUpdateR)
						}
					}

					// Propagate updated restrictions or remove branchMatcher if no restrictions
					if len(creationRMap) > 0 {
						*creationBmMap[bmMatching].Restrictions = toList(creationRMap)
					} else {
						delete(creationBmMap, bmMatching)
					}

					if len(deletionRMap) > 0 {
						*deletionBmMap[bmMatching].Restrictions = toList(deletionRMap)
					} else {
						delete(deletionBmMap, bmMatching)
					}

					// If we've added restrictions to update we append branchMatchers to update
					if len(*toUpdateBm.Restrictions) > 0 {
						*toUpdateBr.BranchMatchers = append(*toUpdateBr.BranchMatchers, toUpdateBm)
					}
				}
			}

			// Propagate updated branchMatchers or remove branchRestriction if no branchMatchers
			if len(creationBmMap) > 0 {
				*creationCandidates[brType].BranchMatchers = toList(creationBmMap)
			} else {
				delete(creationCandidates, brType)
			}

			if len(deletionBmMap) > 0 {
				*deletionCandidates[brType].BranchMatchers = toList(deletionBmMap)
			} else {
				delete(deletionCandidates, brType)
			}

			// If we've added a branchMatcher to update we propagate this to branchRestrictions to update
			if len(*toUpdateBr.BranchMatchers) > 0 {
				*toUpdate = append(*toUpdate, toUpdateBr)
			}
		}

	}

	toCreate = new(BranchRestrictions)
	*toCreate = toList(creationCandidates)
	toDelete = new(BranchRestrictions)
	*toDelete = toList(deletionCandidates)

	return toCreate, toUpdate, toDelete
}

func (brs *BranchRestrictions) AddRestriction(r *types.Restriction) {
	branchRestrictionMap := make(map[string]*BranchRestriction, len(*brs))
	for _, b := range *brs {
		branchRestrictionMap[b.Type] = b
	}

	matcherTypeId := r.Matcher.Type.Id
	if _, exists := branchRestrictionMap[matcherTypeId]; exists {
		branchRestrictionMap[matcherTypeId].BranchMatchers.addRestriction(r)
	} else {
		branchRestrictionMap[matcherTypeId] = &BranchRestriction{
			Type: matcherTypeId,
			BranchMatchers: &BranchMatchers{&BranchMatcher{
				Matching:     r.Matcher.Id,
				Restrictions: &Restrictions{createRestriction(r)},
			}},
		}
		*brs = append(*brs, branchRestrictionMap[matcherTypeId])
	}
}

func (bms *BranchMatchers) addRestriction(r *types.Restriction) {
	matchersMap := make(map[string]*BranchMatcher, len(*bms))
	for _, b := range *bms {
		matchersMap[b.Matching] = b
	}

	matcherId := r.Matcher.Id
	if _, exists := matchersMap[r.Matcher.Id]; exists {
		*matchersMap[matcherId].Restrictions = append(*matchersMap[matcherId].Restrictions, createRestriction(r))
	} else {
		matchersMap[matcherId] = &BranchMatcher{
			Matching:     matcherId,
			Restrictions: &Restrictions{createRestriction(r)},
		}
		*bms = append(*bms, matchersMap[matcherId])
	}
}

func createRestriction(r *types.Restriction) *Restriction {
	var users []string
	for _, u := range r.Users {
		users = append(users, u.Name)
	}

	return &Restriction{
		Id:           r.Id,
		Type:         r.Type,
		ExemptUsers:  users,
		ExemptGroups: r.Groups,
	}
}

func toList[T any](asMap map[string]*T) []*T {
	asList := new([]*T)
	for _, r := range asMap {
		*asList = append(*asList, r)
	}
	return *asList
}

func (brs *BranchRestrictions) toMap() map[string]*BranchRestriction {
	asMap := make(map[string]*BranchRestriction, len(*brs))
	for _, br := range *brs {
		asMap[br.Type] = br
	}
	return asMap
}

func (bms *BranchMatchers) toMap() map[string]*BranchMatcher {
	asMap := make(map[string]*BranchMatcher, len(*bms))
	for _, bm := range *bms {
		asMap[bm.Matching] = bm
	}
	return asMap
}

func (rs *Restrictions) toMap() map[string]*Restriction {
	asMap := make(map[string]*Restriction, len(*rs))
	for _, r := range *rs {
		asMap[r.Type] = r
	}
	return asMap
}

func (brs *BranchRestrictions) findBranchRestrictionsDifference(cmp *BranchRestrictions) (difference *BranchRestrictions) {
	baseMap := brs.toMap()
	comparisonMap := cmp.toMap()

	difference = new(BranchRestrictions)
	for brType, br := range baseMap {
		if _, exists := comparisonMap[brType]; !exists {
			*difference = append(*difference, br)
		} else {
			branchMatcherDifference := br.BranchMatchers.findBranchMatchersDifference(comparisonMap[brType].BranchMatchers)
			if len(*branchMatcherDifference) > 0 {
				*difference = append(*difference, &BranchRestriction{
					Type:           brType,
					BranchMatchers: branchMatcherDifference,
				})
			}
		}
	}
	return difference
}

func (bms *BranchMatchers) findBranchMatchersDifference(cmp *BranchMatchers) (difference *BranchMatchers) {
	baseMap := bms.toMap()
	comparisonMap := cmp.toMap()

	difference = new(BranchMatchers)
	for matching, bm := range baseMap {
		if _, exists := comparisonMap[matching]; !exists {
			*difference = append(*difference, bm)
		} else {
			restrictionsDifference := bm.Restrictions.findRestrictionsDifference(comparisonMap[matching].Restrictions)
			if len(*restrictionsDifference) > 0 {
				*difference = append(*difference, &BranchMatcher{
					Matching:     matching,
					Restrictions: restrictionsDifference,
				})
			}
		}
	}
	return difference
}

func (rs *Restrictions) findRestrictionsDifference(cmp *Restrictions) (difference *Restrictions) {
	baseMap := rs.toMap()
	comparisonMap := cmp.toMap()

	difference = new(Restrictions)
	for rType, r := range baseMap {
		if _, exists := comparisonMap[rType]; !exists {
			*difference = append(*difference, r)
		} else {
			// Assume there's no intersection between usernames and group names
			baseExemptions := append(r.ExemptUsers, r.ExemptGroups...)
			comparisonExemptions := append(comparisonMap[rType].ExemptUsers, comparisonMap[rType].ExemptGroups...)
			if !common.SlicesContainsSameElements(baseExemptions, comparisonExemptions) {
				*difference = append(*difference, r)
			}
		}
	}
	return difference
}

func (brs *BranchRestrictions) Equals(cmp *BranchRestrictions) bool {
	if brs == cmp {
		return true
	}
	if cmp == nil {
		return false
	}
	if len(*brs) != len(*cmp) {
		return false
	}
	baseMap := brs.toMap()
	cmpMap := cmp.toMap()
	if len(baseMap) != len(cmpMap) {
		return false
	}
	for name, r := range baseMap {
		if !r.Equals(cmpMap[name]) {
			return false
		}
	}
	return true
}

func (br *BranchRestriction) Equals(cmp *BranchRestriction) bool {
	if br == cmp {
		return true
	}
	if cmp == nil {
		return false
	}
	if br.Type != cmp.Type {
		return false
	}
	if !br.BranchMatchers.Equals(cmp.BranchMatchers) {
		return false
	}
	return true

}

func (bms *BranchMatchers) Equals(cmp *BranchMatchers) bool {
	if bms == cmp {
		return true
	}
	if cmp == nil {
		return false
	}
	if len(*bms) != len(*cmp) {
		return false
	}
	baseMap := bms.toMap()
	cmpMap := cmp.toMap()
	if len(baseMap) != len(cmpMap) {
		return false
	}
	for name, r := range baseMap {
		if !r.Equals(cmpMap[name]) {
			return false
		}
	}
	return true
}

func (bm *BranchMatcher) Equals(cmp *BranchMatcher) bool {
	if bm == cmp {
		return true
	}
	if cmp == nil {
		return false
	}
	if bm.Matching != cmp.Matching {
		return false
	}
	if !bm.Restrictions.Equals(cmp.Restrictions) {
		return false
	}
	return true
}

func (rs *Restrictions) Equals(cmp *Restrictions) bool {
	if rs == cmp {
		return true
	}
	if cmp == nil {
		return false
	}
	if len(*rs) != len(*cmp) {
		return false
	}
	baseMap := rs.toMap()
	cmpMap := cmp.toMap()
	if len(baseMap) != len(cmpMap) {
		return false
	}
	for name, r := range baseMap {
		if !r.Equals(cmpMap[name]) {
			return false
		}
	}
	return true
}

func (r *Restriction) Equals(cmp *Restriction) bool {
	if r == cmp {
		return true
	}
	if cmp == nil {
		return false
	}
	if r.Id != cmp.Id {
		return false
	}
	if r.Type != cmp.Type {
		return false
	}
	if !common.SlicesContainsSameElements(r.ExemptGroups, cmp.ExemptGroups) {
		return false
	}
	if !common.SlicesContainsSameElements(r.ExemptUsers, cmp.ExemptUsers) {
		return false
	}
	return true
}
