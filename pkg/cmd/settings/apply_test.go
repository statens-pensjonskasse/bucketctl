package settings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	desired = map[string]*Restrictions{
		"WANTED": {
			Branches: map[string]*BranchRestrictions{
				"wanted": {
					Restrictions: map[string]*Restriction{
						"wanted": {
							ExemptGroups: []string{"wanted"},
						},
						"toAdd": {},
						"toUpdate": {
							ExemptUsers: []string{"wanted", "toAdd"},
						},
					},
				},
				"toUpdate": {
					Restrictions: map[string]*Restriction{
						"toUpdate": {},
					},
				},
			},
		},
		"TO_ADD": {
			Branches: map[string]*BranchRestrictions{
				"toAdd": {
					Restrictions: map[string]*Restriction{
						"toAdd": {
							ExemptGroups: []string{"toAdd"},
						},
					},
				},
			},
		},
	}

	actual = map[string]*Restrictions{
		"NOT_WANTED": {
			Branches: map[string]*BranchRestrictions{
				"notWanted": {
					Restrictions: map[string]*Restriction{
						"notWanted": {},
					},
				},
			},
		},
		"WANTED": {
			Branches: map[string]*BranchRestrictions{
				"wanted": {
					Restrictions: map[string]*Restriction{
						"wanted": {
							ExemptGroups: []string{"wanted"},
						},
						"notWanted": {},
						"toUpdate": {
							ExemptUsers: []string{"wanted", "notWanted"},
						},
					},
				},
				"notWanted": {
					Restrictions: map[string]*Restriction{
						"notWanted": {},
					},
				},
			},
		},
	}
)

func Test_findRestrictionsToUpdate(t *testing.T) {
	toUpdate := findRestrictionsToUpdate(desired, actual)

	// Expect to add wanted restrictions
	if toUpdate["TO_ADD"] == nil {
		t.Errorf("Expected to find restriction to add")
	} else if toUpdate["TO_ADD"].Branches["toAdd"] == nil {
		t.Errorf("Expected to find restriction on branch to add")
	} else if toUpdate["TO_ADD"].Branches["toAdd"].Restrictions["toAdd"] == nil {
		t.Errorf("Expected to find restriction type of branch to add")
	} else if toUpdate["TO_ADD"].Branches["toAdd"].Restrictions["toAdd"].ExemptGroups == nil {
		t.Errorf("Expected to find restriction exemption of group to add")
	}

	// Expect to only add new restrictions
	if toUpdate["WANTED"] == nil {
		t.Errorf("Expected to find matcher type where subtype is to be updated")
	} else if toUpdate["WANTED"].Branches["wanted"] == nil {
		t.Errorf("Expected to find branch where subtype is to be updated")
	} else if toUpdate["WANTED"].Branches["wanted"].Restrictions["toAdd"] == nil {
		t.Errorf("Expected to find restriction to be added")
	} else if toUpdate["WANTED"].Branches["wanted"].Restrictions["toUpdate"] == nil {
		t.Errorf("Expected to find restriction to update")
	} else if toUpdate["WANTED"].Branches["wanted"].Restrictions["toUpdate"].ExemptUsers == nil {
		t.Errorf("Expected to find restriction exemptions")
	} else if toUpdate["WANTED"].Branches["wanted"].Restrictions["wanted"] != nil {
		t.Errorf("Expected to not update restriction that's already in place")
	} else if toUpdate["WANTED"].Branches["wanted"].Restrictions["toUpdate"] == nil {
		t.Errorf("Expected to update restriction with different exemptions")
	}

	assert.ElementsMatch(t,
		toUpdate["WANTED"].Branches["wanted"].Restrictions["toUpdate"].ExemptUsers,
		desired["WANTED"].Branches["wanted"].Restrictions["toUpdate"].ExemptUsers,
		"Expected to update restrictions with desired restricitons")
}

func Test_findRestrictionsToDelete(t *testing.T) {
	toDelete := findRestrictionsToDelete(desired, actual)

	// Expect to delete not wanted restrictions
	if toDelete["NOT_WANTED"] == nil {
		t.Errorf("Expected to find unwanted restrictions")
	} else if toDelete["NOT_WANTED"].Branches["notWanted"] == nil {
		t.Errorf("Expected to find unwanted restriction of branch")
	} else if toDelete["NOT_WANTED"].Branches["notWanted"].Restrictions["notWanted"] == nil {
		t.Errorf("Expected to find unwanted restriction type of branch")
	}

	// Expect to only delete not wanted sub-restrictions
	if toDelete["WANTED"] == nil {
		t.Errorf("Expected to find matcher type where subtype is unwanted")
	} else if toDelete["WANTED"].Branches["wanted"] == nil {
		t.Errorf("Expected to find branch where subtype is unwanted")
	} else if toDelete["WANTED"].Branches["notWanted"] == nil {
		t.Errorf("Expected to find branch where restriction is unwanted")
	} else if toDelete["WANTED"].Branches["wanted"].Restrictions["notWanted"] == nil {
		t.Errorf("Expected to find unwanted branch restrictions")
	} else if toDelete["WANTED"].Branches["wanted"].Restrictions["desired"] != nil {
		t.Errorf("Expected to not find wanted branch restriction")
	}
}
