package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"snykctl/internal/tools"
)

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

func (u *Users) Get() error {
	return u.baseGet(false)
}

func (u *Users) GetRaw() (string, error) {
	err := u.baseGet(true)
	if err != nil {
		return "", err
	}
	return u.rawResponse, nil
}

func (u *Users) baseGet(raw bool) error {
	path := fmt.Sprintf("/org/%s/members", u.Org.Id)
	resp := u.client.RequestGet(path)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GetUsers failed: %s", resp.Status)
	}

	if raw {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("GetProjects failed: %s", err)
		}
		u.rawResponse = string(bodyBytes)
	} else {
		var result []*User
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
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