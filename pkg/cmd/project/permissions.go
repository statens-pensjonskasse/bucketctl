package project

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
	"os"
	"strings"
)

const (
	None = iota
	Read
	Write
	Admin
)

type Permission uint8

var (
	PermissionType = map[uint8]string{
		None:  "NONE",
		Read:  "PROJECT_READ",
		Write: "PROJECT_WRITE",
		Admin: "PROJECT_ADMIN",
	}
	PermissionValue = map[string]uint8{
		"NONE":          None,
		"PROJECT_READ":  Read,
		"PROJECT_WRITE": Write,
		"PROJECT_ADMIN": Admin,
	}
)

func (s *Permission) UnmarshalJSON(data []byte) (err error) {
	var permission string
	if err := json.Unmarshal(data, &permission); err != nil {
		return err
	}
	if *s, err = ParsePermission(permission); err != nil {
		return err
	}
	return nil
}

func ParsePermission(s string) (Permission, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	value, ok := PermissionValue[s]
	if !ok {
		return Permission(0), fmt.Errorf("%q is not a valid permission", s)
	}
	return Permission(value), nil

}

type PermissionSet map[Permission][]string

type ProjectPermissions struct {
	Groups []Group
	User   []User
}

type PSet map[Permission]*ProjectPermissions

type Group struct {
	Name string `json:"name"`
}

type GroupPermission struct {
	Group      Group      `json:"group"`
	Permission Permission `json:"permission"`
}

type groupPermissions struct {
	pkg.BitbucketResponse
	Values []GroupPermission `json:"values"`
}

type User struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	Active       bool   `json:"active"`
	DisplayName  string `json:"displayName"`
	Id           int    `json:"id"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
}

type UserPermission struct {
	User       User       `json:"user"`
	Permission Permission `json:"permission"`
}

type userPermissions struct {
	pkg.BitbucketResponse
	Values []UserPermission `json:"values"`
}

func getProjectGroupPermissions(baseUrl string, projectKey string, token string, limit int) (groupPermissions, error) {
	groupPermissionsUrl := fmt.Sprintf("%s/rest/api/1.0/projects/%s/permissions/groups?limit=%d", baseUrl, projectKey, limit)

	body, err := pkg.GetRequestBody(groupPermissionsUrl, token)
	if err != nil {
		return groupPermissions{}, err
	}

	var groups groupPermissions
	if err := json.Unmarshal(body, &groups); err != nil {
		return groupPermissions{}, err
	}

	return groups, nil
}

func getProjectUserPermissions(baseUrl string, projectKey string, token string, limit int) (userPermissions, error) {
	userPermissionsUrl := fmt.Sprintf("%s/rest/api/1.0/projects/%s/permissions/users?limit=%d", baseUrl, projectKey, limit)

	body, err := pkg.GetRequestBody(userPermissionsUrl, token)
	if err != nil {
		return userPermissions{}, err
	}

	var users userPermissions
	if err := json.Unmarshal(body, &users); err != nil {
		return userPermissions{}, err
	}

	return users, nil
}

func printGroupPermissions(permissions []GroupPermission) {
	var data [][]string

	data = append(data, []string{"Permission", "Groups"})

	permissionSet := make(PermissionSet)

	for _, p := range permissions {
		permissionSet[p.Permission] = append(permissionSet[p.Permission], p.Group.Name)
	}

	for permission, _ := range permissionSet {
		var groups string
		for _, g := range permissionSet[permission] {
			groups += g + "\n"
		}
		row := []string{PermissionType[uint8(permission)], strings.Trim(groups, "\n")}
		data = append(data, row)
	}

	pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}

func printUserPermissions(permissions []UserPermission) {
	var data [][]string

	data = append(data, []string{"Permission", "Users"})

	permissionSet := make(PermissionSet)

	for _, p := range permissions {
		permissionSet[p.Permission] = append(permissionSet[p.Permission], p.User.Name)
	}

	for permission, _ := range permissionSet {
		var users string
		for _, g := range permissionSet[permission] {
			users += g + "\n"
		}
		row := []string{PermissionType[uint8(permission)], strings.Trim(users, "\n")}
		data = append(data, row)
	}

	pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}

func listPermissions(cmd *cobra.Command, args []string) {
	var baseUrl = viper.GetString("baseUrl")
	var projectKey = viper.GetString("key")
	var token = viper.GetString("token")
	var limit = viper.GetInt("limit")

	gPerms, err := getProjectGroupPermissions(baseUrl, projectKey, token, limit)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}

	//userPermissions, err := getProjectUserPermissions(baseUrl, projectKey, token, limit)
	//if err != nil {
	//	pterm.Error.Println(err)
	//	os.Exit(1)
	//}

	pSet := make(PSet)
	pSet[None] = new(ProjectPermissions)
	pSet[Read] = new(ProjectPermissions)
	pSet[Write] = new(ProjectPermissions)
	pSet[Admin] = new(ProjectPermissions)
	for _, gp := range gPerms.Values {
		pSet[gp.Permission].Groups = append(pSet[gp.Permission].Groups, gp.Group)
	}

	pterm.Println(pSet)
	pterm.Println(pSet[3].Groups)

	//	printGroupPermissions(groupPermissions.Values)
	//	printUserPermissions(userPermissions.Values)
	//
	//	if !groupPermissions.IsLastPage {
	//		pterm.Warning.Println("Not all groupPermissions fetched, try with a higher limit")
	//	}
	//
	//	if !userPermissions.IsLastPage {
	//		pterm.Warning.Println("Not all userPermissions fetched, try with a higher limit")
	//	}
}
