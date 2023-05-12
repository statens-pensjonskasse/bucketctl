package project

import "testing"

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
	permissionSetA := &PermissionSet{Permissions: map[string]*Entities{
		"PROJECT_ADMIN": {
			Users:  []string{},
			Groups: []string{"A", "B"}},
	}}

	permissionSetA.Permissions["PROJECT_ADMIN"] = new(Entities)

}
