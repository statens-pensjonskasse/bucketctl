package v1alpha1

import (
	"strings"

	"git.spk.no/infra/bucketctl/pkg/common"
)

type Entities struct {
	Groups []string `json:"groups,omitempty" yaml:"groups,omitempty"`
	Users  []string `json:"users,omitempty" yaml:"users,omitempty"`
}

type Permission struct {
	Name     string
	Entities *Entities `json:",omitempty,inline" yaml:",omitempty,inline"`
}

type Permissions []*Permission

func FindPermissionsToChange(desired *Permissions, actual *Permissions) (toCreate *Permissions, toDelete *Permissions) {
	if desired != nil {
		toCreate = desired.getPermissionsDifference(actual)
	} else {
		toCreate = new(Permissions)
	}
	if actual != nil {
		toDelete = actual.getPermissionsDifference(desired)
	} else {
		toDelete = new(Permissions)
	}
	return toCreate, toDelete
}

func UpdatePublicProperty(desired *ProjectConfigSpec, actual *ProjectConfigSpec) *bool {
	if desired.Public != nil && actual.Public != nil && *desired.Public != *actual.Public {
		return desired.Public
	}
	return nil
}

func UpdateDefaultProjectPermissionProperty(desired *ProjectConfigSpec, actual *ProjectConfigSpec) *string {
	if desired.DefaultPermission != nil && actual.DefaultPermission != nil && *desired.DefaultPermission != *actual.DefaultPermission {
		return desired.DefaultPermission
	}
	return nil
}

func (e *Entities) ContainsUser(user string) bool {
	for _, u := range e.Users {
		if strings.ToLower(user) == strings.ToLower(u) {
			return true
		}
	}
	return false
}

func (e *Entities) ContainsGroup(group string) bool {
	for _, g := range e.Groups {
		if strings.ToLower(group) == strings.ToLower(g) {
			return true
		}
	}
	return false
}

func (ps *Permissions) Copy() *Permissions {
	psCopy := new(Permissions)
	for _, p := range *ps {
		*psCopy = append(*psCopy, p.Copy())
	}
	return psCopy
}

func (p *Permission) Copy() *Permission {
	return &Permission{
		Name:     p.Name,
		Entities: p.Entities.Copy(),
	}
}

func (e *Entities) Copy() *Entities {
	eCopy := new(Entities)
	eCopy.Groups = append(eCopy.Groups, e.Groups...)
	eCopy.Users = append(eCopy.Users, e.Users...)
	return eCopy
}

// getPermissionsDifference Find the relative complement of B in ps, ps\B (What's in B that is not in ps)
func (ps *Permissions) getPermissionsDifference(B *Permissions) (difference *Permissions) {
	if B == nil || len(*B) == 0 {
		return ps.Copy()
	}

	relComplement := make(map[string]*Entities)

	setA := make(map[string]*Entities)
	for _, p := range *ps {
		setA[p.Name] = p.Entities
	}

	setB := make(map[string]*Entities)
	for _, p := range *B {
		setB[p.Name] = p.Entities
	}

	for permission := range setA {
		relComplement[permission] = new(Entities)
		entriesInA := setA[permission]
		entriesInB, existsInB := setB[permission]

		if !existsInB {
			// If the permission doesn't exist in B we add all entries of A in the relative complement
			relComplement[permission] = setA[permission]
		} else {
			// If the permission exist in B we have to check all the entries
			for _, user := range entriesInA.Users {
				if !entriesInB.ContainsUser(user) {
					relComplement[permission].Users = append(relComplement[permission].Users, user)
				}
			}
			for _, group := range entriesInA.Groups {
				if !entriesInB.ContainsGroup(group) {
					relComplement[permission].Groups = append(relComplement[permission].Groups, group)
				}
			}
		}
	}

	difference = new(Permissions)
	for name, e := range relComplement {
		// If there are changes in users or groups of a permission we add it
		if len(e.Users)+len(e.Groups) > 0 {
			*difference = append(*difference, &Permission{
				Name:     name,
				Entities: e,
			})
		}
	}

	return difference
}

func (ps *Permissions) toMap() map[string]*Permission {
	asMap := make(map[string]*Permission, len(*ps))
	for _, p := range *ps {
		asMap[p.Name] = p
	}
	return asMap
}

func (ps *Permissions) Equals(cmp *Permissions) bool {
	if ps == cmp {
		return true
	}
	if cmp == nil {
		return false
	}
	if len(*ps) != len(*cmp) {
		return false
	}
	psMap := ps.toMap()
	cmpMap := cmp.toMap()
	for name, p := range psMap {
		if !p.Equals(cmpMap[name]) {
			return false
		}
	}
	return true
}

func (p *Permission) Equals(cmp *Permission) bool {
	if p == cmp {
		return true
	}
	if cmp == nil {
		return false
	}
	if p.Name != cmp.Name {
		return false
	}
	if !p.Entities.Equals(cmp.Entities) {
		return false
	}
	return true
}

func (e *Entities) Equals(cmp *Entities) bool {
	if e == cmp {
		return true
	}
	if cmp == nil {
		return false
	}
	if !common.SlicesContainsSameElementsIgnoringCase(e.Groups, cmp.Groups) {
		return false
	}
	if !common.SlicesContainsSameElementsIgnoringCase(e.Users, cmp.Users) {
		return false
	}
	return true
}
