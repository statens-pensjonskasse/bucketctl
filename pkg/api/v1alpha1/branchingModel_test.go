package v1alpha1

import (
	"bucketctl/pkg/api/bitbucket/types"
	"testing"
)

func Test_BranchingModel_isEmpty(t *testing.T) {
	type fields struct {
		Development *types.Branch
		Production  *types.Branch
		Types       *BranchingModelTypes
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "all nil fields",
			fields: fields{nil, nil, nil},
			want:   true,
		},
		{
			name:   "nil branches and empty branchModelTypes",
			fields: fields{nil, nil, &BranchingModelTypes{}},
			want:   true,
		},
		{
			name:   "only development branch",
			fields: fields{&types.Branch{}, nil, nil},
			want:   false,
		},
		{
			name:   "only production branch",
			fields: fields{nil, &types.Branch{}, nil},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := &BranchingModel{
				Development: tt.fields.Development,
				Production:  tt.fields.Production,
				Types:       tt.fields.Types,
			}
			if got := bm.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

var (
	mainBranch                  = "/refs/heads/master"
	useDefaultDevelopmentBranch = &BranchingModel{
		Development: &types.Branch{UseDefault: true},
		Production:  nil,
		Types:       nil,
	}
	namedDevelopmentBranch = &BranchingModel{
		Development: &types.Branch{RefId: &mainBranch},
		Production:  nil,
		Types:       nil,
	}
)

func Test_FindBranchingModelsToChange(t *testing.T) {
	type args struct {
		desired *BranchingModel
		actual  *BranchingModel
	}
	type wants struct {
		toCreate *BranchingModel
		toUpdate *BranchingModel
		toDelete *BranchingModel
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "nil in, no change out",
			args:  args{nil, nil},
			wants: wants{nil, nil, nil},
		},
		{
			name:  "empty in, no change out",
			args:  args{&BranchingModel{}, &BranchingModel{}},
			wants: wants{nil, nil, nil},
		},
		{
			name:  "use defaultBranch as developmentBranch",
			args:  args{useDefaultDevelopmentBranch, &BranchingModel{}},
			wants: wants{useDefaultDevelopmentBranch, nil, nil},
		},
		{
			name:  "change from named branch to defaultBranch",
			args:  args{useDefaultDevelopmentBranch, namedDevelopmentBranch},
			wants: wants{nil, useDefaultDevelopmentBranch, nil},
		},
		{
			name:  "remove named branch as default development branch",
			args:  args{nil, useDefaultDevelopmentBranch},
			wants: wants{nil, nil, useDefaultDevelopmentBranch},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToCreate, gotToUpdate, gotToDelete := FindBranchingModelsToChange(tt.args.desired, tt.args.actual)
			if !gotToCreate.Equals(tt.wants.toCreate) {
				t.Errorf("FindBranchingModelsToChange() gotToCreate = %v, want %v", gotToCreate, tt.wants.toCreate)
			}
			if !gotToUpdate.Equals(tt.wants.toUpdate) {
				t.Errorf("FindBranchingModelsToChange() gotToUpdate = %v, want %v", gotToUpdate, tt.wants.toUpdate)
			}
			if !gotToDelete.Equals(tt.wants.toDelete) {
				t.Errorf("FindBranchingModelsToChange() gotToDelete = %v, want %v", gotToDelete, tt.wants.toDelete)
			}
		})
	}
}

var (
	bmt = &BranchingModelType{
		Name:   "name",
		Prefix: "prefix",
	}
	bmts = &BranchingModelTypes{bmt}
	bm   = &BranchingModel{
		Development: &types.Branch{UseDefault: true},
		Production:  &types.Branch{Id: "main"},
		Types:       bmts,
	}
)

func Test_BranchingModelType_Copy(t *testing.T) {
	tests := []struct {
		name string
		recv BranchingModelType
		want *BranchingModelType
	}{
		{
			name: "Copy empty",
			recv: BranchingModelType{},
			want: &BranchingModelType{},
		},
		{
			name: "Copy",
			recv: *bmt,
			want: bmt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.recv.Copy(); !got.Equals(tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}

	copyBranchingModelType := bmt.Copy()

	// Change every value
	copyBranchingModelType.Name = ""
	copyBranchingModelType.Prefix = ""

	if copyBranchingModelType.Name == bmt.Name {
		t.Errorf("Name not copied")
	}
	if copyBranchingModelType.Prefix == bmt.Prefix {
		t.Errorf("Prefix not copied")
	}
}

func Test_BranchingModelTypes_Copy(t *testing.T) {
	tests := []struct {
		name string
		recv BranchingModelTypes
		want *BranchingModelTypes
	}{
		{
			name: "Copy empty",
			recv: BranchingModelTypes{},
			want: &BranchingModelTypes{},
		},
		{
			name: "Copy",
			recv: *bmts,
			want: bmts,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.recv.Copy(); !got.Equals(tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}

	bmtsCopy := bmts.Copy()
	(*bmtsCopy)[0].Name = ""

	if (*bmtsCopy)[0].Equals((*bmts)[0]) {
		t.Errorf("BranchModelType not copied")
	}
}

func Test_BranchingModel_Copy(t *testing.T) {
	tests := []struct {
		name string
		recv BranchingModel
		want *BranchingModel
	}{
		{
			name: "Copy empty",
			recv: BranchingModel{},
			want: &BranchingModel{},
		},
		{
			name: "Copy",
			recv: *bm,
			want: bm,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.recv.Copy(); !got.Equals(tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}

	bmCopy := bm.Copy()

	bmCopy.Development.Type = "type"
	bmCopy.Production = &types.Branch{}
	bmCopy.Types = &BranchingModelTypes{}

	if bmCopy.Development.Equals(bm.Development) {
		t.Errorf("development not copied")
	}
	if bmCopy.Production.Equals(bm.Production) {
		t.Errorf("production not copied")
	}
	if bmCopy.Types.Equals(bm.Types) {
		t.Errorf("types not copied")
	}

}
