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

	desiredPermissions := &Permissions{
		permission: &Entities{
			Users:  []string{"User1", "User2"},
			Groups: []string{"Group1", "Group2", "Group3"},
		},
	}

	actualPermissions := &Permissions{
		permission: &Entities{
			Users:  []string{"User3", "User4"},
			Groups: []string{"Group3", "Group4"},
		},
	}

	expectedToBeGranted := &Permissions{
		permission: &Entities{
			Users:  []string{"User1", "User2"},
			Groups: []string{"Group1", "Group2"},
		},
	}

	actualToBeGranted := desiredPermissions.getPermissionsDifference(actualPermissions)
	if !reflect.DeepEqual(expectedToBeGranted, actualToBeGranted) {
		t.Fatal("Forventer å gi tilgang til 'User1', 'User2', 'Group1' og 'Group2'")
	}

	expectedToBeRemoved := &Permissions{
		permission: &Entities{
			Users:  []string{"User3", "User4"},
			Groups: []string{"Group4"},
		},
	}
	actualToBeRemoved := actualPermissions.getPermissionsDifference(desiredPermissions)
	if !reflect.DeepEqual(expectedToBeRemoved, actualToBeRemoved) {
		t.Fatal("Forventer å fjerne tilgang for 'User3', 'User4' og 'Group4'")
	}
}
