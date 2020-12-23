package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

// TestConfigParses validates if yaml annotation do not overlap
func TestConfigParses(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, yaml.Unmarshal([]byte(""), &cfg))
}