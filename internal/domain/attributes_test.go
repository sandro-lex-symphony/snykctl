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
	msg         string
}

func Test_ParseAttributes(t *testing.T) {
	tests := []testDataAttributes{
		testDataAttributes{env: "frontend", lifecycle: "", criticality: "", msg: ""},
		testDataAttributes{env: "xxx", lifecycle: "", criticality: "", msg: "invalid environment value: xxx\nValid values: [frontend backend internal external mobile saas on-prem hosted distributed]"},
		testDataAttributes{env: "", lifecycle: "production", criticality: "", msg: ""},
		testDataAttributes{env: "", lifecycle: "xxx", criticality: "", msg: "invalid lifecycle value: xxx\nValid values: [production development sandbox]"},
		testDataAttributes{env: "", lifecycle: "", criticality: "high", msg: ""},
		testDataAttributes{env: "", lifecycle: "", criticality: "xxx", msg: "invalid lifecycle value: xxx\nValid values: [critical high medium low]"},
		testDataAttributes{env: "frontend", lifecycle: "production", criticality: "medium", msg: ""},
		testDataAttributes{env: "xxx", lifecycle: "xxx", criticality: "xxx", msg: "invalid environment value: xxx\nValid values: [frontend backend internal external mobile saas on-prem hosted distributed]"},
	}
	// var err error
	for _, test := range tests {
		err := ParseAttributes(test.env, test.lifecycle, test.criticality)
		if err == nil {
			assert.Equal(t, test.msg, "")
		} else {
			assert.EqualErrorf(t, err, test.msg, "Error should be: %v, got: %v", test.msg, err)
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

func Test_BuildAttributeFilter(t *testing.T) {
	tests := []testDataAttributes{
		testDataAttributes{env: "", lifecycle: "", criticality: "", msg: ""},
		testDataAttributes{env: "frontend", lifecycle: "", criticality: "", msg: `"attributes": { "environment": ["frontend"] }`},
		testDataAttributes{env: "", lifecycle: "production", criticality: "", msg: `"attributes": { "lifecycle": ["production"] }`},
		testDataAttributes{env: "", lifecycle: "", criticality: "medium", msg: `"attributes": { "criticality": ["medium"] }`},
		testDataAttributes{env: "backend", lifecycle: "development", criticality: "medium", msg: `"attributes": { "environment": ["backend"],"lifecycle": ["development"],"criticality": ["medium"] }`},
	}

	for _, test := range tests {
		out := BuildAttributesFilter(test.env, test.lifecycle, test.criticality)
		assert.Equal(t, out, test.msg)
	}
}

func Test_BuildTagsFilter(t *testing.T) {
	var pTags map[string]string
	pTags = make(map[string]string)
	pTags["key"] = "value"

	out := BuildTagsFilter(pTags)
	expected := ` "tags": { "includes": [{ "key": "key", "value": "value" } ] }`
	assert.Equal(t, expected, out)

	pTags["k2"] = "v2"
	out = BuildTagsFilter(pTags)
	expected = ` "tags": { "includes": [{ "key": "key", "value": "value" } , { "key": "k2", "value": "v2" } ] }`
	assert.Equal(t, expected, out)
}

func Test_BuildFilterBody(t *testing.T) {
	var out, expected string
	var pTags map[string]string

	// attrs empty, tags empty
	out = BuildFilterBody("", "", "", pTags)
	expected = ""
	assert.Equal(t, expected, out)

	// attrs not empty
	out = BuildFilterBody("frontend", "", "", pTags)
	expected = `{ "filters": { "attributes": { "environment": ["frontend"] } } }`
	assert.Equal(t, expected, out)

	// attrs empty, tags not empty
	pTags = make(map[string]string)
	pTags["k"] = "v"
	out = BuildFilterBody("", "", "", pTags)
	expected = `{ "filters": {  "tags": { "includes": [{ "key": "k", "value": "v" } ] } } }`
	assert.Equal(t, expected, out)

	// attrs not empty / tags not empty
	out = BuildFilterBody("frontend", "", "", pTags)
	expected = `{ "filters": { "attributes": { "environment": ["frontend"] }, "tags": { "includes": [{ "key": "k", "value": "v" } ] } } }`
	assert.Equal(t, expected, out)

}
