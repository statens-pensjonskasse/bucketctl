package v1alpha1

import (
	"reflect"
	"testing"
)

func Test_EntitiesContains(t *testing.T) {
	entities := &Entities{
		Users:  []string{"😐"},
		Groups: []string{"👯"},
	}

	if !entities.ContainsUser("😐") {
		t.Fatal("Should contain '😐'")
	}
	if entities.ContainsUser("🫥") {
		t.Fatal("Should not contain '🫥'")
	}
	if !entities.ContainsGroup("👯") {
		t.Fatal("Should contain '👯'")
	}
	if entities.ContainsGroup("🤼") {
		t.Fatal("Should not contain '🤼'")
	}
}

func Test_Copy(t *testing.T) {
	orig := &Permissions{
		{
			Name: "copy_me",
			Entities: &Entities{
				Users:  []string{"User1", "User2"},
				Groups: []string{"Group1", "Group2", "Group3"},
			},
		},
	}
	cpy := orig.Copy()

	if cpy == orig {
		t.Errorf("Expected a pointer to a new object")
	}
	if !reflect.DeepEqual(cpy, orig) {
		t.Errorf("Expected a copy of all values")
	}

	(*cpy)[0].Name = "copy"
	(*cpy)[0].Entities.Users[0] = "Δ"

	if (*cpy)[0].Name == (*orig)[0].Name {
		t.Errorf("Expected original to retain name")
	}

	if (*cpy)[0].Entities.Users[0] == (*orig)[0].Entities.Users[0] {
		t.Errorf("Expected original to retain same user name")
	}

}

var (
	read = &Permission{
		Name: "READ",
		Entities: &Entities{
			Users: []string{"reader"},
		},
	}
	admin = &Permission{
		Name: "ADMIN",
		Entities: &Entities{
			Users: []string{"admin"},
		},
	}
)

func Test_Permissions_getPermissionsDifference(t *testing.T) {
	tests := []struct {
		name string
		ps   Permissions
		args *Permissions
		want *Permissions
	}{
		{
			name: "no difference",
			ps:   Permissions{read},
			args: &Permissions{read},
			want: new(Permissions),
		},
		{
			name: "something comparing empty",
			ps:   Permissions{admin},
			args: new(Permissions),
			want: &Permissions{admin},
		},
		{
			name: "empty comparing something",
			ps:   Permissions{},
			args: &Permissions{read},
			want: &Permissions{},
		},
		{
			name: "missing one",
			ps:   Permissions{read, admin},
			args: &Permissions{admin},
			want: &Permissions{read},
		},
		{
			name: "including",
			ps:   Permissions{read},
			args: &Permissions{read, admin},
			want: new(Permissions),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ps.getPermissionsDifference(tt.args); !got.Equals(tt.want) {
				t.Errorf("%s, getPermissionsDifference() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_getPermissionsDifference(t *testing.T) {
	var permission = "DUMMY_PERMISSION"

	desiredPermissions := &Permissions{
		{
			Name: permission,
			Entities: &Entities{
				Users:  []string{"User1", "User2"},
				Groups: []string{"Group1", "Group2", "Group3"},
			},
		},
	}

	allDesired := desiredPermissions.getPermissionsDifference(new(Permissions))
	if allDesired == desiredPermissions {
		t.Errorf("Expected a pointer to a new object")
	}
	if !reflect.DeepEqual(desiredPermissions, allDesired) {
		t.Errorf("Expected all values when paired with an empty slice")
	}

	actualPermissions := &Permissions{
		{
			Name: permission,
			Entities: &Entities{
				Users:  []string{"User3", "User4"},
				Groups: []string{"Group3", "Group4"},
			},
		},
	}

	expectedToBeGranted := &Permissions{
		{
			Name: permission,
			Entities: &Entities{
				Users:  []string{"User1", "User2"},
				Groups: []string{"Group1", "Group2"},
			},
		},
	}

	actualToBeGranted := desiredPermissions.getPermissionsDifference(actualPermissions)
	if !reflect.DeepEqual(expectedToBeGranted, actualToBeGranted) {
		t.Errorf("Forventer å gi tilgang til 'User1', 'User2', 'Group1' og 'Group2'")
	}

	expectedToBeRemoved := &Permissions{
		{
			Name: permission,
			Entities: &Entities{
				Users:  []string{"User3", "User4"},
				Groups: []string{"Group4"},
			},
		},
	}

	actualToBeRemoved := actualPermissions.getPermissionsDifference(desiredPermissions)
	if !reflect.DeepEqual(expectedToBeRemoved, actualToBeRemoved) {
		t.Errorf("Forventer å fjerne tilgang for 'User3', 'User4' og 'Group4'")
	}
}

func Test_group(t *testing.T) {
	desired := &RepositoriesProperties{
		{
			RepoSlug: "AAA",
			Permissions: &Permissions{
				{
					Name:     "desired",
					Entities: &Entities{Users: []string{"user"}},
				},
			},
		},
		{
			RepoSlug: "CCC",
			Permissions: &Permissions{
				{
					Name:     "both",
					Entities: &Entities{Users: []string{"user"}},
				},
			},
		},
	}

	actual := &RepositoriesProperties{
		{
			RepoSlug: "BBB",
			Permissions: &Permissions{
				{
					Name:     "actual",
					Entities: &Entities{Users: []string{"user"}},
				},
			},
		},
		{
			RepoSlug: "CCC",
			Permissions: &Permissions{
				{
					Name:     "both",
					Entities: &Entities{Users: []string{"user"}},
				},
			},
		},
	}

	grouped := GroupRepositories(desired, actual)

	if grouped["AAA"].Desired == nil {
		t.Errorf("Expected grouped permissions to contain desired permissions for repo %s", "AAA")
	}
	if len(*grouped["AAA"].Actual.Permissions) != 0 {
		t.Errorf("Expected grouped permissions to not contain any permissions for repo %s", "AAA")
	}

	if len(*grouped["BBB"].Desired.Permissions) != 0 {
		t.Errorf("Expected grouped permissions to not contain any permissions for repo %s", "BBB")
	}
	if grouped["BBB"].Actual == nil {
		t.Errorf("Expected grouped permissions to contain actual permissions for repo %s", "BBB")
	}

	if grouped["CCC"].Desired == nil {
		t.Errorf("Expected grouped permissions to contain desired permissions for repo %s", "CCC")
	}
	if grouped["CCC"].Actual == nil {
		t.Errorf("Expected grouped permissions to contain actual permissions for repo %s", "CCC")
	}
}

func Test_FindPermissionsToChange(t *testing.T) {
	type args struct {
		desired *Permissions
		actual  *Permissions
	}
	type want struct {
		toCreate *Permissions
		toDelete *Permissions
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no difference",
			args: args{&Permissions{read}, &Permissions{read}},
			want: want{new(Permissions), new(Permissions)},
		},
		{
			name: "something comparing empty",
			args: args{&Permissions{read, admin}, new(Permissions)},
			want: want{&Permissions{read, admin}, new(Permissions)},
		},
		{
			name: "empty comparing something",
			args: args{new(Permissions), &Permissions{read, admin}},
			want: want{new(Permissions), &Permissions{read, admin}},
		},
		{
			name: "add permission",
			args: args{&Permissions{read, admin}, &Permissions{read}},
			want: want{&Permissions{admin}, new(Permissions)},
		},
		{
			name: "remove permission",
			args: args{&Permissions{admin}, &Permissions{admin, read}},
			want: want{new(Permissions), &Permissions{read}},
		},
		{
			name: "add and remove permission",
			args: args{&Permissions{admin}, &Permissions{read}},
			want: want{&Permissions{admin}, &Permissions{read}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toCreate, toDelete := FindPermissionsToChange(tt.args.desired, tt.args.actual)
			if !toCreate.Equals(tt.want.toCreate) {
				t.Errorf("%s got create %v, want %v", tt.name, toCreate, tt.want.toCreate)
			}
			if !toDelete.Equals(tt.want.toDelete) {
				t.Errorf("%s got delete %v, want %v", tt.name, toDelete, tt.want.toDelete)
			}
		})
	}
}
