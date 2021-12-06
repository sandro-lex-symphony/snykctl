package domain

import (
	"fmt"
	"snykctl/internal/tools"
)

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

	colsize := 40
	mid := ""
	left := fillSpaces(orgName1, colsize, " ")
	fmt.Printf("%s%s%s\n", left, mid, orgName2)

	leftBar := fillSpaces("", len(orgName1), "=")
	leftBar = fillSpaces(leftBar, colsize, " ")
	right := fillSpaces("", len(orgName2), "=")
	fmt.Printf("%s%s%s\n", leftBar, mid, right)

	r3 := mergeUsers(users1.Users, users2.Users)
	for i := 0; i < len(r3); i++ {
		if containUser(users1.Users, r3[i]) && containUser(users2.Users, r3[i]) {
			fmt.Printf("%s%s%s\n", fillSpaces(r3[i].Name, colsize, " "), mid, r3[i].Name)
		} else if containUser(users1.Users, r3[i]) {
			fmt.Printf("%s%s--- MISSING ---\n", fillSpaces(r3[i].Name, colsize, " "), mid)
		} else {
			fmt.Printf("%s%s%s\n", fillSpaces("--- MISSING ---", colsize, " "), mid, r3[i].Name)
		}
	}

	return nil
}

func addSpaces(size int, filler string) string {
	ret := ""
	for i := 0; i < size; i++ {
		ret += filler
	}
	return ret
}

func fillSpaces(s string, size int, fillerChar string) string {
	if len(s) >= size {
		return s
	}
	filler := addSpaces(size-len(s), fillerChar)
	return s + filler
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
