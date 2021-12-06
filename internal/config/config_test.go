package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_obfuscated(t *testing.T) {
	var config ConfigProperties
	config.SetSync(true) // avoid reading from filesystem

	token := "abcdefgh"
	want := "cdefgh"

	config.SetToken(token)
	assert.Equal(t, want, config.ObfuscatedToken())

	config.SetId(token)
	assert.Equal(t, want, config.ObfuscatedId())
}

func Test_Init_defaultValues(t *testing.T) {
	var config ConfigProperties
	config.Init(false)
	// avoid reading conf again
	config.SetSync(true)

	assert.Equal(t, 10, config.Timeout())
	assert.Equal(t, 10, config.WorkerSize())
	assert.Equal(t, "https://snyk.io/api/v1", config.Url())
	assert.Equal(t, ".snykctl.yaml", config.Filename())
}

func Test_ReadConf(t *testing.T) {
	var config ConfigProperties
	config.SetFilename(("tmp.yaml"))

	err := config.ReadConf()
	home, _ := os.UserHomeDir()

	expectedErrorMsg := fmt.Sprintf("readConf failed: open %s/tmp.yaml: no such file or directory", home)
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}
