package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testAttributesBody struct {
	env         string
	lifecycle   string
	criticality string
	body        string
}

func Test_Build_AttributesBody(t *testing.T) {
	tests := []testAttributesBody{
		testAttributesBody{env: "", lifecycle: "", criticality: "", body: ""},
		testAttributesBody{env: "production", lifecycle: "", criticality: "", body: `{ "environment": ["production"] }`},
		testAttributesBody{env: "", lifecycle: "frontend", criticality: "", body: `{ "lifecycle": ["frontend"] }`},
		testAttributesBody{env: "", lifecycle: "", criticality: "high", body: `{ "criticality": ["high"] }`},
		testAttributesBody{env: "production", lifecycle: "frontend", criticality: "", body: `{ "environment": ["production"],"lifecycle": ["frontend"] }`},
		testAttributesBody{env: "production", lifecycle: "frontend", criticality: "medium", body: `{ "environment": ["production"],"lifecycle": ["frontend"],"criticality": ["medium"] }`},
	}

	for _, test := range tests {
		body := BuildAttributesBody(test.env, test.lifecycle, test.criticality)
		assert.Equal(t, test.body, body)
	}
}

type testDataAttributes struct {
	env         string
	lifecycle   string
	criticality string
	errorMsg    string
}

func Test_ParseAttributes(t *testing.T) {
	tests := []testDataAttributes{
		testDataAttributes{env: "frontend", lifecycle: "", criticality: "", errorMsg: ""},
		testDataAttributes{env: "xxx", lifecycle: "", criticality: "", errorMsg: "invalid environment value: xxx\nValid values: [frontend backend internal external mobile saas on-prem hosted distributed]"},
		testDataAttributes{env: "", lifecycle: "production", criticality: "", errorMsg: ""},
		testDataAttributes{env: "", lifecycle: "xxx", criticality: "", errorMsg: "invalid lifecycle value: xxx\nValid values: [production development sandbox]"},
		testDataAttributes{env: "", lifecycle: "", criticality: "high", errorMsg: ""},
		testDataAttributes{env: "", lifecycle: "", criticality: "xxx", errorMsg: "invalid lifecycle value: xxx\nValid values: [critical high medium low]"},
		testDataAttributes{env: "frontend", lifecycle: "production", criticality: "medium", errorMsg: ""},
		testDataAttributes{env: "xxx", lifecycle: "xxx", criticality: "xxx", errorMsg: "invalid environment value: xxx\nValid values: [frontend backend internal external mobile saas on-prem hosted distributed]"},
	}
	// var err error
	for _, test := range tests {
		err := ParseAttributes(test.env, test.lifecycle, test.criticality)
		if err == nil {
			assert.Equal(t, test.errorMsg, "")
		} else {
			assert.EqualErrorf(t, err, test.errorMsg, "Error should be: %v, got: %v", test.errorMsg, err)
		}
	}
}

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

func Test_ParseTags(t *testing.T) {
	// one ok
	tag := []string{"key=value"}
	pTags, err := ParseTags(tag)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(pTags))
	assert.Equal(t, pTags["key"], "value")

	// two ok
	tag = append(tag, "key2=value2")
	pTags, err = ParseTags(tag)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(pTags))
	assert.Equal(t, pTags["key"], "value")
	assert.Equal(t, pTags["key2"], "value2")

	// ko
	tag = []string{"qweasd"}
	_, err = ParseTags(tag)
	expected := "invalid tag. Not a key=value format"
	assert.EqualErrorf(t, err, expected, "Error should be: %v, got: %v", expected, err)
}
