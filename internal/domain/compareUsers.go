package domain

import (
	"fmt"
	"snykctl/internal/tools"
	"strings"
)

const missing = "--- MISSING ---"

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
