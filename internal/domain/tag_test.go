package domain

import (
	"fmt"
	"net/http"
	"snykctl/internal/tools"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	tag   string
	key   string
	value string
	err   error
}

func Test_ParseTag(t *testing.T) {
	tests := []testData{
		testData{tag: "a=b", key: "a", value: "b", err: nil},
		testData{tag: "abc", key: "", value: "", err: fmt.Errorf("invalid tag. Not a key=value format")},
		testData{tag: "==", key: "", value: "", err: fmt.Errorf("invalid tag. Not a key=value format")},
		testData{tag: "=", key: "", value: "", err: fmt.Errorf("invalid tag. Not a key=value format")},
		testData{tag: "a=", key: "", value: "", err: fmt.Errorf("invalid tag. Not a key=value format")},
		testData{tag: "=b", key: "", value: "", err: fmt.Errorf("invalid tag. Not a key=value format")},
	}

	for _, test := range tests {
		k, v, err := ParseTag(test.tag)
		if k != test.key {
			t.Errorf("ParseTag(%#v) want key %s, got key %s", test.tag, test.key, k)
		}
		if v != test.value {
			t.Errorf("ParseTag(%#v) want value %s, got value %s", test.tag, test.value, v)
		}
		if err == nil && test.err != nil {
			t.Errorf("ParseTag(%#v) want err %s, got err %s", test.tag, test.err, err)
		}
	}
}

func Test_AddTag_OK(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ``
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	err := AddTag(client, "org1", "org2", "k=v")
	assert.Nil(t, err)
}

func Test_AddTag_parseFailed(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ``
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	err := AddTag(client, "org1", "org2", "vvv")
	expectedErrorMsg := "invalid tag. Not a key=value format"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

func Test_AddTag_KO(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ``
	client.StatusCode = http.StatusUnauthorized
	client.Status = "XXX"
	err := AddTag(client, "org1", "org2", "k=v")
	expectedErrorMsg := "Failed to add tag XXX"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

// k, v, err := ParseTag(tag)
// if err != nil {
// 	return err
// }

// path := fmt.Sprintf(tagPath, org_id, prj_id)

// tagBody := fmt.Sprintf(`{"key": "%s", "value": "%s"}`, k, v)
// var jsonStr = []byte(tagBody)

// resp := client.RequestPost(path, jsonStr)
// defer resp.Body.Close()
// if resp.StatusCode != http.StatusOK {
// 	return fmt.Errorf("Failed to add tag %s", resp.Status)
// }
// return nil
