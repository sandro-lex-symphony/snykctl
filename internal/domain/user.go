package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"snykctl/internal/config"
	"snykctl/internal/tools"
	"strings"
)

const membersPath = "/org/%s/members"
const groupMembersPath = "/group/%s/members"
const addUserPath = "/group/%s/org/%s/members"
const deleteUserPath = "/org/%s/members/%s"

const missing = "--- MISSING ---"

type Users struct {
	Users       []*User
	Org         Org
	client      tools.HttpClient
	sync        bool
	rawResponse string
}

type User struct {
	Id    string
	Name  string
	Role  string
	Email string
}

func NewUsers(c tools.HttpClient, org_id string) *Users {
	u := new(Users)
	u.Org.Id = org_id
	u.SetClient(c)
	return u
}

func (u *Users) SetClient(c tools.HttpClient) {
	u.client = c
}

func (u *Users) GetGroup() error {
	return u.baseGet(false, groupMembersPath)
}

func (u *Users) GetGroupRaw() (string, error) {
	if err := u.baseGet(true, groupMembersPath); err != nil {
		return "", err
	}
	return u.rawResponse, nil
}

func (u *Users) Get() error {
	return u.baseGet(false, membersPath)
}

func (u *Users) GetRaw() (string, error) {
	if err := u.baseGet(true, membersPath); err != nil {
		return "", err
	}
	return u.rawResponse, nil
}

func (u *Users) baseGet(raw bool, endpoint string) error {
	path := fmt.Sprintf(endpoint, u.Org.Id)
	resp := u.client.RequestGet(path)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GetUsers failed: %s", resp.Status)
	}

	if raw {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("GetUsers failed: %s", err)
		}
		u.rawResponse = string(bodyBytes)
	} else {
		var result []*User
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return fmt.Errorf("GetUsers failed: %s", err)
		}
		u.Users = result
	}

	u.sync = true
	return nil
}

func (u Users) String() string {
	return u.toString("")
}

func (u Users) Quiet() string {
	return u.toString("id")
}

func (u Users) Name() string {
	return u.toString("name")
}

func (u Users) toString(filter string) string {
	var out string
	for _, user := range u.Users {
		if filter == "id" {
			out += fmt.Sprintf("%s\n", user.Id)
		} else if filter == "name" {
			out += fmt.Sprintf("%s\n", user.Name)
		} else {
			out += fmt.Sprintf("%-38s %-14s%s\n", user.Id, user.Role, user.Name)
		}
	}
	return out
}

func (u Users) Sync() bool {
	return u.sync
}

func AddUser(client tools.HttpClient, org_id, user_id, role string) error {
	path := fmt.Sprintf(addUserPath, config.Instance.Id(), org_id)
	jsonValue, _ := json.Marshal(map[string]string{
		"userId": user_id,
		"role":   role,
	})
	resp := client.RequestPost(path, jsonValue)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("add User failed: %s", resp.Status)
	}
	return nil
}

func DeleteUser(client tools.HttpClient, org_id, user_id string) error {
	path := fmt.Sprintf(deleteUserPath, org_id, user_id)
	resp := client.RequestDelete(path)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("deleteUsers failed: %s", resp.Status)
	}

	return nil
}

func CopyUsers(client tools.HttpClient, org1, org2 string) error {
	users1 := NewUsers(client, org1)
	if err := users1.Get(); err != nil {
		return err
	}

	for _, user := range users1.Users {
		if err := AddUser(client, org2, user.Id, "collaborator"); err != nil {
			return err
		}
	}

	return nil
}

func CompareUsers(client tools.HttpClient, org1, org2 string) error {
	orgs := NewOrgs(client)

	var err error
	orgName1, err := orgs.GetOrgName(org1)
	if err != nil {
		return err
	}
	users1 := NewUsers(client, org1)
	err = users1.Get()
	if err != nil {
		return err
	}

	orgName2, err := orgs.GetOrgName(org2)
	if err != nil {
		return err
	}
	users2 := NewUsers(client, org2)
	err = users2.Get()
	if err != nil {
		return err
	}

	out := compare(orgName1, orgName2, users1.Users, users2.Users)
	fmt.Print(out)

	return nil
}

func compare(orgName1 string, orgName2 string, users1 []*User, users2 []*User) string {
	var out string
	out += fmt.Sprintf("%-40s%s\n", orgName1, orgName2)

	leftBar := strings.Repeat("=", len(orgName1))
	rightBar := strings.Repeat("=", len(orgName2))
	out += fmt.Sprintf("%-40s%s\n", leftBar, rightBar)

	r3 := mergeUsers(users1, users2)
	for i := 0; i < len(r3); i++ {
		if containUser(users1, r3[i]) && containUser(users2, r3[i]) {
			out += fmt.Sprintf("%-40s%s\n", r3[i].Name, r3[i].Name)
		} else if containUser(users1, r3[i]) {
			out += fmt.Sprintf("%-40s%s\n", r3[i].Name, missing)
		} else {
			out += fmt.Sprintf("%-40s%s\n", missing, r3[i].Name)
		}
	}

	return out
}

func mergeUsers(u1 []*User, u2 []*User) []*User {
	var u3 []*User
	u3 = u1
	for i := 0; i < len(u2); i++ {
		if !containUser(u1, u2[i]) {
			u3 = append(u3, u2[i])
		}

	}
	return u3
}

func containUser(u1 []*User, x *User) bool {
	for _, v := range u1 {
		if v.Id == x.Id {
			return true
		}
	}
	return false
}
