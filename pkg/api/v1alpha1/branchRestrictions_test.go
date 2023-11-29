package v1alpha1

import (
	"reflect"
	"testing"
)

var (
	// restrictionA (rA)
	rA = &Restriction{
		Id:           0,
		Type:         "restriction_a",
		ExemptUsers:  []string{"userA"},
		ExemptGroups: []string{"groupA"},
	}
	rB = &Restriction{
		Id:           1,
		Type:         "restriction_b",
		ExemptUsers:  []string{"userA", "userB"},
		ExemptGroups: []string{"groupA"},
	}
	rC = &Restriction{
		Id:           3,
		Type:         "restriction_c",
		ExemptUsers:  []string{"userA", "userB"},
		ExemptGroups: []string{"groupA", "groupC"},
	}
	rADiffExemptions = &Restriction{
		Id:           0,
		Type:         rA.Type,
		ExemptUsers:  []string{"userA", "userB"},
		ExemptGroups: []string{"groupA"},
	}
)

func Test_Restrictions_findRestrictionsDifference(t *testing.T) {
	type args struct {
		comparison *Restrictions
	}
	tests := []struct {
		name string
		base Restrictions
		args args
		want *Restrictions
	}{
		{
			name: "base restrictions equal comparison restrictions",
			base: Restrictions{rA},
			args: args{&Restrictions{rA}},
			want: new(Restrictions),
		},
		{
			name: "no restrictions in comparison",
			base: Restrictions{rA},
			args: args{new(Restrictions)},
			want: &Restrictions{rA},
		},
		{
			name: "different restrictions in comparison",
			base: Restrictions{rA},
			args: args{&Restrictions{rB}},
			want: &Restrictions{rA},
		},
		{
			name: "comparison includes restrictions",
			base: Restrictions{rA},
			args: args{&Restrictions{rA, rB}},
			want: new(Restrictions),
		},
		{
			name: "comparison includes restrictions with different exemptions",
			base: Restrictions{rA},
			args: args{&Restrictions{rADiffExemptions}},
			want: &Restrictions{rA},
		},
		{
			name: "multiple missing restrictions",
			base: Restrictions{rA, rB, rC},
			args: args{&Restrictions{rADiffExemptions}},
			want: &Restrictions{rA, rB, rC},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.base.findRestrictionsDifference(tt.args.comparison); !got.Equals(tt.want) {
				t.Errorf("%s - got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

var (
	// branchMatcherA (bmA)
	bmA = &BranchMatcher{
		Matching:     "matcher a",
		Restrictions: new(Restrictions),
	}
	// branchMatcherA (bmA) with restrictionA (rA)
	bmArA = &BranchMatcher{
		Matching:     bmA.Matching,
		Restrictions: &Restrictions{rA},
	}
	bmArB = &BranchMatcher{
		Matching:     bmA.Matching,
		Restrictions: &Restrictions{rB},
	}
	bmB = &BranchMatcher{
		Matching:     "matcher b",
		Restrictions: new(Restrictions),
	}
	bmBrB = &BranchMatcher{
		Matching:     bmB.Matching,
		Restrictions: &Restrictions{rB},
	}
	bmArArB = &BranchMatcher{
		Matching:     bmA.Matching,
		Restrictions: &Restrictions{rA, rB},
	}
	bmArBrA = &BranchMatcher{
		Matching:     bmA.Matching,
		Restrictions: &Restrictions{rB, rA},
	}
)

func Test_BranchMatchers_findBranchMatchersNotIn(t *testing.T) {
	type args struct {
		comparison *BranchMatchers
	}
	tests := []struct {
		name string
		base BranchMatchers
		args args
		want *BranchMatchers
	}{
		{
			name: "comparing itself",
			base: BranchMatchers{bmArA},
			args: args{&BranchMatchers{bmArA}},
			want: new(BranchMatchers),
		},
		{
			name: "comparing with empty",
			base: BranchMatchers{bmArA, bmBrB},
			args: args{new(BranchMatchers)},
			want: &BranchMatchers{bmArA, bmBrB},
		},
		{
			name: "comparing empty with empty",
			base: BranchMatchers{},
			args: args{new(BranchMatchers)},
			want: new(BranchMatchers),
		},
		{
			name: "missing restriction in comparison",
			base: BranchMatchers{bmArArB},
			args: args{&BranchMatchers{bmArA}},
			want: &BranchMatchers{bmArB},
		},
		{
			name: "extra restriction in comparison",
			base: BranchMatchers{bmArA},
			args: args{&BranchMatchers{bmArArB}},
			want: new(BranchMatchers),
		},
		{
			name: "comparing same set, but different ordering",
			base: BranchMatchers{bmArArB},
			args: args{&BranchMatchers{bmArBrA}},
			want: new(BranchMatchers),
		},
		{
			name: "comparing different branchMatchers with same restrictions",
			base: BranchMatchers{bmBrB},
			args: args{&BranchMatchers{bmArB}},
			want: &BranchMatchers{bmBrB},
		},
		{
			name: "empty branchMatcher",
			base: BranchMatchers{bmA},
			args: args{&BranchMatchers{bmArArB}},
			want: new(BranchMatchers),
		},
		{
			name: "empty comparison",
			base: BranchMatchers{bmArArB, bmBrB},
			args: args{&BranchMatchers{bmA, bmBrB}},
			want: &BranchMatchers{bmArArB},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.base.findBranchMatchersDifference(tt.args.comparison); !got.Equals(tt.want) {
				t.Errorf("%s - got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

var (
	// branchRestrictionA (brA)
	brA = &BranchRestriction{
		Type:           "branch restriction a",
		BranchMatchers: new(BranchMatchers),
	}
	// branchRestrictionA (brA) with branchMatcherA (bmA) and restrictionA (rA)
	brAbmArA = &BranchRestriction{
		Type:           brA.Type,
		BranchMatchers: &BranchMatchers{bmArA},
	}
	brAbmArArB = &BranchRestriction{
		Type:           brA.Type,
		BranchMatchers: &BranchMatchers{bmArArB},
	}
	brAbmArB = &BranchRestriction{
		Type:           brA.Type,
		BranchMatchers: &BranchMatchers{bmArB},
	}
	brB = &BranchRestriction{
		Type:           "branch restriction b",
		BranchMatchers: new(BranchMatchers),
	}
	brBbmArA = &BranchRestriction{
		Type:           brB.Type,
		BranchMatchers: &BranchMatchers{bmArA},
	}
	brBbmArB = &BranchRestriction{
		Type:           brB.Type,
		BranchMatchers: &BranchMatchers{bmArB},
	}
	brBbmBrB = &BranchRestriction{
		Type:           brB.Type,
		BranchMatchers: &BranchMatchers{bmBrB},
	}
	brBbmArArBbmBrB = &BranchRestriction{
		Type:           brB.Type,
		BranchMatchers: &BranchMatchers{bmArArB, bmBrB},
	}
	brBbmArBbmBrB = &BranchRestriction{
		Type:           brB.Type,
		BranchMatchers: &BranchMatchers{bmArB, bmBrB},
	}
)

func Test_BranchRestrictions_findBranchRestrictionsNotIn(t *testing.T) {
	type args struct {
		comparison *BranchRestrictions
	}
	tests := []struct {
		name string
		base BranchRestrictions
		args args
		want *BranchRestrictions
	}{
		{
			name: "comparing itself",
			base: BranchRestrictions{brAbmArA},
			args: args{&BranchRestrictions{brAbmArA}},
			want: new(BranchRestrictions),
		},
		{
			name: "comparing with empty",
			base: BranchRestrictions{brAbmArArB},
			args: args{new(BranchRestrictions)},
			want: &BranchRestrictions{brAbmArArB},
		},
		{
			name: "missing restriction in comparison",
			base: BranchRestrictions{brAbmArArB},
			args: args{&BranchRestrictions{brAbmArA}},
			want: &BranchRestrictions{brAbmArB},
		},
		{
			name: "extra restriction in comparison",
			base: BranchRestrictions{brAbmArA},
			args: args{&BranchRestrictions{brAbmArArB}},
			want: new(BranchRestrictions),
		},
		{
			name: "comparison with different branch matcher, but same restriction",
			base: BranchRestrictions{brBbmArB},
			args: args{&BranchRestrictions{brBbmBrB}},
			want: &BranchRestrictions{brBbmArB},
		},
		{
			name: "complex comparison",
			base: BranchRestrictions{brBbmArArBbmBrB},
			args: args{&BranchRestrictions{brAbmArA, brBbmArA}},
			want: &BranchRestrictions{brBbmArBbmBrB},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.base.findBranchRestrictionsDifference(tt.args.comparison); !got.Equals(tt.want) {
				t.Errorf("%s - got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

var (
	bmArADiff = &BranchMatcher{
		Matching:     bmA.Matching,
		Restrictions: &Restrictions{rADiffExemptions},
	}
	brAbmArADiff = &BranchRestriction{
		Type:           brA.Type,
		BranchMatchers: &BranchMatchers{bmArADiff},
	}
	brAbmBrB = &BranchRestriction{
		Type:           brA.Type,
		BranchMatchers: &BranchMatchers{bmBrB},
	}
	brBbmArADiff = &BranchRestriction{
		Type:           brB.Type,
		BranchMatchers: &BranchMatchers{bmArADiff},
	}
)

func TestFindBranchRestrictionsToChange(t *testing.T) {
	type args struct {
		desired *BranchRestrictions
		actual  *BranchRestrictions
	}
	type wants struct {
		toCreate *BranchRestrictions
		toUpdate *BranchRestrictions
		toDelete *BranchRestrictions
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "only create branch restriction",
			args:  args{&BranchRestrictions{brAbmArA}, new(BranchRestrictions)},
			wants: wants{&BranchRestrictions{brAbmArA}, new(BranchRestrictions), new(BranchRestrictions)},
		},
		{
			name:  "only delete branch restriction",
			args:  args{&BranchRestrictions{brAbmArA}, &BranchRestrictions{brAbmArArB}},
			wants: wants{new(BranchRestrictions), new(BranchRestrictions), &BranchRestrictions{brAbmArB}},
		},
		{ // Here we can just update exemptions on existing Bitbucket Restriction
			name:  "change restriction exemptions",
			args:  args{&BranchRestrictions{brAbmArADiff}, &BranchRestrictions{brAbmArA}},
			wants: wants{new(BranchRestrictions), &BranchRestrictions{brAbmArADiff}, new(BranchRestrictions)},
		},
		{ // Here we have to delete and recreate due to Bitbucket API
			name:  "change restriction for branch restriction",
			args:  args{&BranchRestrictions{brAbmArB}, &BranchRestrictions{brAbmArA}},
			wants: wants{&BranchRestrictions{brAbmArB}, new(BranchRestrictions), &BranchRestrictions{brAbmArA}},
		},
		{
			name:  "change branch matcher for branch restriction",
			args:  args{&BranchRestrictions{brAbmArB}, &BranchRestrictions{brAbmBrB}},
			wants: wants{&BranchRestrictions{brAbmArB}, new(BranchRestrictions), &BranchRestrictions{brAbmBrB}},
		},
		{
			name:  "change branch matcher for branch restriction and add restriction",
			args:  args{&BranchRestrictions{brAbmArArB}, &BranchRestrictions{brAbmBrB}},
			wants: wants{&BranchRestrictions{brAbmArArB}, new(BranchRestrictions), &BranchRestrictions{brAbmBrB}},
		},
		{
			name:  "change branch matcher for branch restriction and remove restriction",
			args:  args{&BranchRestrictions{brAbmBrB}, &BranchRestrictions{brAbmArArB}},
			wants: wants{&BranchRestrictions{brAbmBrB}, new(BranchRestrictions), &BranchRestrictions{brAbmArArB}},
		},
		{
			name:  "change branch restriction",
			args:  args{&BranchRestrictions{brBbmArA}, &BranchRestrictions{brAbmArA}},
			wants: wants{&BranchRestrictions{brBbmArA}, new(BranchRestrictions), &BranchRestrictions{brAbmArA}},
		},
		{
			name:  "add restriction to branchRestriction",
			args:  args{&BranchRestrictions{brAbmArArB}, &BranchRestrictions{brAbmArA}},
			wants: wants{&BranchRestrictions{brAbmArB}, new(BranchRestrictions), new(BranchRestrictions)},
		},
		{
			name:  "update and delete",
			args:  args{&BranchRestrictions{brBbmArADiff}, &BranchRestrictions{brBbmArArBbmBrB}},
			wants: wants{new(BranchRestrictions), &BranchRestrictions{brBbmArADiff}, &BranchRestrictions{brBbmArBbmBrB}},
		},
		{
			name:  "update and create",
			args:  args{&BranchRestrictions{brAbmArArB}, &BranchRestrictions{brAbmArADiff}},
			wants: wants{&BranchRestrictions{brAbmArB}, &BranchRestrictions{brAbmArA}, new(BranchRestrictions)},
		},
		{
			name:  "create and delete",
			args:  args{&BranchRestrictions{brAbmArA}, &BranchRestrictions{brAbmArB}},
			wants: wants{&BranchRestrictions{brAbmArA}, new(BranchRestrictions), &BranchRestrictions{brAbmArB}},
		},
		{
			name:  "create, update and delete",
			args:  args{&BranchRestrictions{brAbmArA, brBbmArADiff}, &BranchRestrictions{brBbmArArBbmBrB}},
			wants: wants{&BranchRestrictions{brAbmArA}, &BranchRestrictions{brBbmArADiff}, &BranchRestrictions{brBbmArBbmBrB}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToCreate, gotToUpdate, gotToDelete := FindBranchRestrictionsToChange(tt.args.desired, tt.args.actual)
			if !gotToCreate.Equals(tt.wants.toCreate) {
				t.Errorf("%s - create got %v, want %v", tt.name, gotToCreate, tt.wants.toCreate)
			}
			if !gotToUpdate.Equals(tt.wants.toUpdate) {
				t.Errorf("%s - update got %v, want %v", tt.name, gotToUpdate, tt.wants.toUpdate)
			}
			if !gotToDelete.Equals(tt.wants.toDelete) {
				t.Errorf("%s - delete got %v, want %v", tt.name, gotToDelete, tt.wants.toDelete)
			}
		})
	}
}

func TestBranchRestrictions_Copy(t *testing.T) {
	tests := []struct {
		name string
		brs  BranchRestrictions
		want *BranchRestrictions
	}{
		{
			name: "Copy empty",
			brs:  BranchRestrictions{},
			want: &BranchRestrictions{},
		},
		{
			name: "Copy with single entry",
			brs: BranchRestrictions{&BranchRestriction{
				Type:           "test",
				BranchMatchers: &BranchMatchers{},
			}},
			want: &BranchRestrictions{&BranchRestriction{
				Type:           "test",
				BranchMatchers: &BranchMatchers{},
			}},
		},
		{
			name: "Deep copy",
			brs: BranchRestrictions{&BranchRestriction{
				Type: "test",
				BranchMatchers: &BranchMatchers{
					&BranchMatcher{
						Matching: "deep",
						Restrictions: &Restrictions{
							&Restriction{
								Id:           0,
								Type:         "type",
								ExemptUsers:  nil,
								ExemptGroups: nil,
							},
						},
					},
				},
			}},
			want: &BranchRestrictions{&BranchRestriction{
				Type: "test",
				BranchMatchers: &BranchMatchers{
					&BranchMatcher{
						Matching: "deep",
						Restrictions: &Restrictions{
							&Restriction{
								Id:           0,
								Type:         "type",
								ExemptUsers:  nil,
								ExemptGroups: nil,
							},
						},
					},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.brs.Copy(); !got.Equals(tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_BranchMatchers_Copy(t *testing.T) {
	tests := []struct {
		name string
		bms  BranchMatchers
		want *BranchMatchers
	}{
		{
			name: "Copy empty",
			bms:  BranchMatchers{},
			want: &BranchMatchers{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bms.Copy(); !got.Equals(tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestriction_Copy(t *testing.T) {
	type fields struct {
		Id           int
		Type         string
		ExemptUsers  []string
		ExemptGroups []string
	}
	tests := []struct {
		name   string
		fields fields
		want   *Restriction
	}{
		{
			name: "empty copy",
			fields: fields{
				Id:           0,
				Type:         "",
				ExemptUsers:  nil,
				ExemptGroups: nil,
			},
			want: &Restriction{
				Id:           0,
				Type:         "",
				ExemptUsers: nil,
				ExemptGroups: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Restriction{
				Id:           tt.fields.Id,
				Type:         tt.fields.Type,
				ExemptUsers:  tt.fields.ExemptUsers,
				ExemptGroups: tt.fields.ExemptGroups,
			}
			if got := r.Copy(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestrictions_Copy(t *testing.T) {
	tests := []struct {
		name string
		rs   Restrictions
		want *Restrictions
	}{
		{
			name: "empty copy",
			rs: Restrictions{
				&Restriction{
					Id:           0,
					Type:         "",
					ExemptUsers:  []string{""},
					ExemptGroups: []string{""},
				},
			},
			want:  &Restrictions{
				&Restriction{
					Id:           0,
					Type:         "",
					ExemptUsers:  []string{""},
					ExemptGroups: []string{""},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rs.Copy(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %+v, want %+v", got, tt.want)
			}
		})
	}
}