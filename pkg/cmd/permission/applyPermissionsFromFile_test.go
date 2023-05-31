package permission

import (
	"reflect"
	"testing"
)

func TestEntitiesContains(t *testing.T) {
	entities := &Entities{
		Users:  []string{"😐"},
		Groups: []string{"👯"},
	}

	if !entities.containsUser("😐") {
		t.Fatal("Skal inneholde '😐'")
	}
	if entities.containsUser("🫥") {
		t.Fatal("Skal ikke inneholde '🫥'")
	}
	if !entities.containsGroup("👯") {
		t.Fatal("Skal inneholde '👯'")
	}
	if entities.containsGroup("🤼") {
		t.Fatal("Skal ikke inneholde '🤼'")
	}
}

func TestGetPermissionSetDifference(t *testing.T) {
	var permission = "DUMMY_PERMISSION"

	desiredPermissions := &ProjectPermissions{
		PermissionSet: Permissions{
			map[string]*Entities{
				permission: {
					Users:  []string{"User1", "User2"},
					Groups: []string{"Group1", "Group2", "Group3"},
				},
			},
		},
	}

	actualPermissions := &ProjectPermissions{
		PermissionSet: Permissions{
			map[string]*Entities{
				permission: {
					Users:  []string{"User3", "User4"},
					Groups: []string{"Group3", "Group4"},
				},
			},
		},
	}

	expectedToBeGranted := &Entities{
		Users:  []string{"User1", "User2"},
		Groups: []string{"Group1", "Group2"},
	}
	actualToBeGranted := desiredPermissions.getPermissionSetDifference(actualPermissions).Permissions[permission]
	if !reflect.DeepEqual(expectedToBeGranted, actualToBeGranted) {
		t.Fatal("Forventer å gi tilgang til 'User1', 'User2', 'Group1' og 'Group2'")
	}

	expectedToBeRemoved := &Entities{
		Users:  []string{"User3", "User4"},
		Groups: []string{"Group4"},
	}
	actualToBeRemoved := actualPermissions.getPermissionSetDifference(desiredPermissions).Permissions[permission]
	if !reflect.DeepEqual(expectedToBeRemoved, actualToBeRemoved) {
		t.Fatal("Forventer å fjerne tilgang for 'User3', 'User4' og 'Group4'")
	}
}
