package types

import (
	"testing"
)

func TestBranch_Copy(t *testing.T) {
	refId := "refId"

	branchA := &Branch{
		Id:              "id",
		RefId:           &refId,
		DisplayId:       "displayId",
		Type:            "type",
		LatestCommit:    "latestCommit",
		LatestChangeset: "latestChangeset",
		UseDefault:      true,
	}

	tests := []struct {
		name string
		recv Branch
		want *Branch
	}{
		{
			name: "Copy empty",
			recv: Branch{},
			want: &Branch{},
		},
		{
			name: "Copy",
			recv: *branchA,
			want: branchA,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.recv.Copy(); !got.Equals(tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}

	copyBranchA := branchA.Copy()

	refIdChanged := ""

	// Change every value
	copyBranchA.Id = ""
	copyBranchA.RefId = &refIdChanged
	copyBranchA.DisplayId = ""
	copyBranchA.Type = ""
	copyBranchA.LatestCommit = ""
	copyBranchA.LatestChangeset = ""
	copyBranchA.UseDefault = !copyBranchA.UseDefault

	if copyBranchA.Id == branchA.Id {
		t.Errorf("id not copied")
	}
	if copyBranchA.RefId == branchA.RefId {
		t.Errorf("refId not copied")
	}
	if *copyBranchA.RefId == *branchA.RefId {
		t.Errorf("refId value not copied")
	}
	if copyBranchA.DisplayId == branchA.DisplayId {
		t.Errorf("displayID not copied")
	}
	if copyBranchA.Type == branchA.Type {
		t.Errorf("type not copied")
	}
	if copyBranchA.LatestCommit == branchA.LatestCommit {
		t.Errorf("latestCommti not copied")
	}
	if copyBranchA.LatestChangeset == branchA.LatestChangeset {
		t.Errorf("latestChangeset not copied")
	}
	if copyBranchA.UseDefault == branchA.UseDefault {
		t.Errorf("useDefault not copied")
	}
}
