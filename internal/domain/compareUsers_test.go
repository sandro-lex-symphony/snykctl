package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
