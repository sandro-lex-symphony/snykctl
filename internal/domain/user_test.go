package domain

import (
	"fmt"
	"net/http"
	"snykctl/internal/tools"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_User_Get_httpError(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ""
	client.StatusCode = http.StatusUnauthorized
	client.Status = "XXX"
	u := NewUsers(client, "xxx")

	err := u.Get()
	expectedErrorMsg := "GetUsers failed: XXX"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
	assert.Equal(t, false, u.Sync())
}

func Test_Users_Get_badBody(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = "filler"
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	u := NewUsers(client, "xxx")

	err := u.Get()
	expectedErrorMsg := "GetUsers failed:"
	assert.Containsf(t, err.Error(), expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
	assert.Equal(t, false, u.Sync())
}

func Test_User_Get_Ok(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `[{"id":"233588e1-fc72-42c1-a42b-01ea0763d954","username":"power.user@example.com","name":"Power User","email":"power.user@example.com","role":"admin"},{"id":"22cfa37d-8b84-49bf-bb75-43ec2c0dba6b","username":"normal.user@example.com","name":"Normal User","email":"normal.user@example.com","role":"collaborator"}]`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	u := NewUsers(client, "xxx")

	err := u.Get()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(u.Users))

	for _, user := range u.Users {
		if user.Id == "233588e1-fc72-42c1-a42b-01ea0763d954" {
			assert.Equal(t, "power.user@example.com", user.Email)
			assert.Equal(t, "Power User", user.Name)
		}
	}
}

func Test_User_GetRaw(t *testing.T) {
	client := tools.NewMockClient()

	raw := `[{"id":"233588e1-fc72-42c1-a42b-01ea0763d954","username":"power.user@example.com","name":"Power User","email":"power.user@example.com","role":"admin"},{"id":"22cfa37d-8b84-49bf-bb75-43ec2c0dba6b","username":"normal.user@example.com","name":"Normal User","email":"normal.user@example.com","role":"collaborator"}]`
	client.ResponseBody = raw
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	u := NewUsers(client, "xxx")

	out, err := u.GetRaw()
	assert.Nil(t, err)
	assert.Equal(t, raw, out)

}

func Test_User_GetGroup_Ok(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `[{"id":"233588e1-fc72-42c1-a42b-01ea0763d954","username":"power.user@example.com","name":"Power User","email":"power.user@example.com","role":"admin"},{"id":"22cfa37d-8b84-49bf-bb75-43ec2c0dba6b","username":"normal.user@example.com","name":"Normal User","email":"normal.user@example.com","role":"collaborator"}]`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	u := NewUsers(client, "xxx")

	err := u.GetGroup()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(u.Users))

	for _, user := range u.Users {
		if user.Id == "233588e1-fc72-42c1-a42b-01ea0763d954" {
			assert.Equal(t, "power.user@example.com", user.Email)
			assert.Equal(t, "Power User", user.Name)
		}
	}
}

func Test_User_GetGroupRaw(t *testing.T) {
	client := tools.NewMockClient()
	raw := `[{"id":"233588e1-fc72-42c1-a42b-01ea0763d954","username":"power.user@example.com","name":"Power User","email":"power.user@example.com","role":"admin"},{"id":"22cfa37d-8b84-49bf-bb75-43ec2c0dba6b","username":"normal.user@example.com","name":"Normal User","email":"normal.user@example.com","role":"collaborator"}]`
	client.ResponseBody = raw
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	u := NewUsers(client, "xxx")

	out, err := u.GetGroupRaw()
	assert.Nil(t, err)
	assert.Equal(t, raw, out)
}

func Test_User_Get_String(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `[{"id":"233588e1-fc72-42c1-a42b-01ea0763d954","username":"power.user@example.com","name":"Power User","email":"power.user@example.com","role":"admin"},{"id":"22cfa37d-8b84-49bf-bb75-43ec2c0dba6b","username":"normal.user@example.com","name":"Normal User","email":"normal.user@example.com","role":"collaborator"}]`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	u := NewUsers(client, "xxx")

	err := u.Get()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(u.Users))

	expected := "233588e1-fc72-42c1-a42b-01ea0763d954   admin         Power User\n22cfa37d-8b84-49bf-bb75-43ec2c0dba6b   collaborator  Normal User\n"
	assert.Equal(t, expected, u.String())
}

func Test_User_Get_Quiet(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `[{"id":"233588e1-fc72-42c1-a42b-01ea0763d954","username":"power.user@example.com","name":"Power User","email":"power.user@example.com","role":"admin"},{"id":"22cfa37d-8b84-49bf-bb75-43ec2c0dba6b","username":"normal.user@example.com","name":"Normal User","email":"normal.user@example.com","role":"collaborator"}]`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	u := NewUsers(client, "xxx")

	err := u.Get()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(u.Users))

	expected := "233588e1-fc72-42c1-a42b-01ea0763d954\n22cfa37d-8b84-49bf-bb75-43ec2c0dba6b\n"
	assert.Equal(t, expected, u.Quiet())
}

func Test_User_Get_Name(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `[{"id":"233588e1-fc72-42c1-a42b-01ea0763d954","username":"power.user@example.com","name":"Power User","email":"power.user@example.com","role":"admin"},{"id":"22cfa37d-8b84-49bf-bb75-43ec2c0dba6b","username":"normal.user@example.com","name":"Normal User","email":"normal.user@example.com","role":"collaborator"}]`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	u := NewUsers(client, "xxx")

	err := u.Get()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(u.Users))

	expected := "Power User\nNormal User\n"
	assert.Equal(t, expected, u.Name())
}

func Test_mergeUsers(t *testing.T) {

	u1 := User{Id: "123", Name: "user1", Role: "admin", Email: "user1@example.com"}
	u2 := User{Id: "456", Name: "user2", Role: "collaborator", Email: "user2@example.com"}
	u3 := User{Id: "789", Name: "user3", Role: "collaborator", Email: "user3@example.com"}

	var users1, users2 []*User

	users1 = append(users1, &u1, &u2, &u3)
	users2 = append(users2, &u1)

	users3 := mergeUsers(users1, users2)

	assert.Equal(t, users3, users1)

}

func Test_containUser(t *testing.T) {
	u1 := User{Id: "123", Name: "user1", Role: "admin", Email: "user1@example.com"}
	u2 := User{Id: "456", Name: "user2", Role: "collaborator", Email: "user2@example.com"}
	u3 := User{Id: "789", Name: "user3", Role: "collaborator", Email: "user3@example.com"}

	var users1 []*User

	users1 = append(users1, &u1, &u2)

	assert.True(t, containUser(users1, &u1))
	assert.True(t, containUser(users1, &u2))
	assert.False(t, containUser(users1, &u3))
}

func Test_compare(t *testing.T) {

	u1 := User{Id: "123", Name: "user1", Role: "admin", Email: "user1@example.com"}
	u2 := User{Id: "456", Name: "user2", Role: "collaborator", Email: "user2@example.com"}
	u3 := User{Id: "789", Name: "user3", Role: "collaborator", Email: "user3@example.com"}

	var users1, users2 []*User

	users1 = append(users1, &u1, &u2)
	users2 = append(users2, &u1, &u3)

	o1 := "o1"
	o2 := "o2"
	bar := "=="
	expected := fmt.Sprintf("%-40s%s\n", o1, o2)
	expected += fmt.Sprintf("%-40s%s\n", bar, bar)
	expected += fmt.Sprintf("%-40s%s\n", u1.Name, u1.Name)
	expected += fmt.Sprintf("%-40s%s\n", u2.Name, missing)
	expected += fmt.Sprintf("%-40s%s\n", missing, u3.Name)
	out := compare("o1", "o2", users1, users2)

	assert.Equal(t, expected, out)

}
