package config

import (
	"testing"

	"gotest.tools/assert"
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
