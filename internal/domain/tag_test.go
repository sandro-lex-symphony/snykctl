package domain

import (
	"fmt"
	"testing"
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
