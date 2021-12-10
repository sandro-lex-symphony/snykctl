package domain

import (
	"net/http"
	"snykctl/internal/tools"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FormatIgnore(t *testing.T) {
	u1 := User{Id: "u1", Name: "u1", Role: "collaboratr", Email: "u1@example.com"}
	content := IgnoreContent{Reason: "x1", Created: "x2", Expires: "x3", ReasonType: "x4", IgnoredBy: u1}
	ignore := Ignore{Id: "id", Content: content}

	out := FormatIgnore(ignore, "")
	expected := "id                            x2                            u1@example.com                x1\n"
	assert.Equal(t, expected, out)

	out = FormatIgnore(ignore, "p1")
	expected = "p1                                    id                            x2                            u1@example.com                x1\n"
	assert.Equal(t, expected, out)
}

func Test_FormatIgnoreResult(t *testing.T) {
	u1 := User{Id: "u1", Name: "u1", Role: "collaboratr", Email: "u1@example.com"}
	content := IgnoreContent{Reason: "x1", Created: "x2", Expires: "x3", ReasonType: "x4", IgnoredBy: u1}
	star := IgnoreStar{Star: content}
	ignoreResult := make(map[string][]IgnoreStar)

	var stars []IgnoreStar
	stars = append(stars, star)
	ignoreResult["v1"] = stars

	out := FormatIgnoreResult(ignoreResult, "")
	expected := "v1                            x2                            u1@example.com                x1\n"
	assert.Equal(t, expected, out)

	out = FormatIgnoreResult(ignoreResult, "p1")
	expected = "p1                                    v1                            x2                            u1@example.com                x1\n"
	assert.Equal(t, expected, out)
}

func Test_GetProjectIgnores_OK(t *testing.T) {
	raw := `{
		"npm:qs:20140806-1": [
		  {
			"*": {
			  "reason": "No fix available",
			  "created": "2017-10-31T11:24:00.932Z",
			  "expires": "2017-12-10T15:39:38.099Z",
			  "ignoredBy": {
				"id": "a3952187-0d8e-45d8-9aa2-036642857b4f",
				"name": "Joe Bloggs",
				"email": "jbloggs@gmail.com"
			  },
			  "reasonType": "temporary-ignore",
			  "disregardIfFixable": true
			}
		  }
		],
		"npm:negotiator:20160616": [
		  {
			"*": {
			  "reason": "Not vulnerable via this path",
			  "created": "2017-10-31T11:24:45.365Z",
			  "ignoredBy": {
				"id": "a3952187-0d8e-45d8-9aa2-036642857b4f",
				"name": "Joe Bloggs",
				"email": "jbloggs@gmail.com"
			  },
			  "reasonType": "not-vulnerable",
			  "disregardIfFixable": false
			}
		  }
		],
		"npm:electron:20170426": [
		  {
			"*": {
			  "reason": "Low impact",
			  "created": "2017-10-31T11:25:17.138Z",
			  "ignoredBy": {
				"id": "a3952187-0d8e-45d8-9aa2-036642857b4f",
				"name": "Joe Bloggs",
				"email": "jbloggs@gmail.com"
			  },
			  "reasonType": "wont-fix",
			  "disregardIfFixable": false
			}
		  }
		]
	  }`
	client := tools.NewMockClient()
	client.ResponseBody = raw
	client.StatusCode = http.StatusOK
	client.Status = "XXX"

	res, err := GetProjectIgnores(client, "org", "prj")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(res))
}

func Test_FormatGetIgnoreResult_KO(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ""
	client.StatusCode = http.StatusUnauthorized
	client.Status = "XXX"

	_, err := GetProjectIgnores(client, "org", "prj")
	expectedErrorMsg := "getProjectsIgnores failed: XXX"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

func Test_FormatGetIgnoreResult_BodyKO(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `{"key":"value"}`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"

	_, err := GetProjectIgnores(client, "org", "prj")
	expectedErrorMsg := "json: cannot unmarshal string into Go value of type []domain.IgnoreStar"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}
