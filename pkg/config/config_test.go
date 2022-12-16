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
		{
			name: "should render tags config",
			cfg: `
target:
  type: aws
  aws:
    accountId: 123456789
    region: ap-southeast-2
    role: arn:aws:iam::123456789012:role/roleName
    ecrOptions:
      tags:
        - key: CreatedBy
          value: k8s-image-swapper
        - key: A
          value: B
`,
			expCfg: Config{
				Target: Target{
					AWS: AWS{
						AccountID: "123456789",
						Region:    "ap-southeast-2",
						Role:      "arn:aws:iam::123456789012:role/roleName",
						ECROptions: ECROptions{
							Tags: []Tag{
								{
									Key:   "CreatedBy",
									Value: "k8s-image-swapper",
								},
								{
									Key:   "A",
									Value: "B",
								},
							},
						},
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
