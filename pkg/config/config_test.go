package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

// TestConfigParses validates if yaml annotation do not overlap
func TestConfigParses(t *testing.T) {
	tests := []struct {
		name   string
		cfg    string
		expCfg Config
		expErr bool
	}{
		{
			name:   "should render empty config",
			cfg:    "",
			expCfg: Config{},
		},
		{
			name: "should render multiple filters",
			cfg: `
source:
  filters:
    - jmespath: "obj.metadata.namespace == 'kube-system'"
    - jmespath: "obj.metadata.namespace != 'playground'"
`,
			expCfg: Config{
				Source: Source{
					Filters: []JMESPathFilter{
						{JMESPath: "obj.metadata.namespace == 'kube-system'"},
						{JMESPath: "obj.metadata.namespace != 'playground'"},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			gotCfg := Config{}
			err := yaml.Unmarshal([]byte(test.cfg), &gotCfg)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expCfg, gotCfg)
			}
		})
	}
}
