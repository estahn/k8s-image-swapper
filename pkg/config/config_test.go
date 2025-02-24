package config

import (
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
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
			name: "should render empty config with defaults",
			cfg:  "",
			expCfg: Config{
				Target: Registry{
					Type: "aws",
					AWS: AWS{
						ECROptions: ECROptions{
							ImageTagMutability: "MUTABLE",
							ImageScanningConfiguration: ImageScanningConfiguration{
								ImageScanOnPush: true,
							},
							EncryptionConfiguration: EncryptionConfiguration{
								EncryptionType: "AES256",
							},
						},
					},
				},
			},
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
				Target: Registry{
					Type: "aws",
					AWS: AWS{
						ECROptions: ECROptions{
							ImageTagMutability: "MUTABLE",
							ImageScanningConfiguration: ImageScanningConfiguration{
								ImageScanOnPush: true,
							},
							EncryptionConfiguration: EncryptionConfiguration{
								EncryptionType: "AES256",
							},
						},
					},
				},
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
				Target: Registry{
					Type: "aws",
					AWS: AWS{
						AccountID: "123456789",
						Region:    "ap-southeast-2",
						Role:      "arn:aws:iam::123456789012:role/roleName",
						ECROptions: ECROptions{
							ImageTagMutability: "MUTABLE",
							ImageScanningConfiguration: ImageScanningConfiguration{
								ImageScanOnPush: true,
							},
							EncryptionConfiguration: EncryptionConfiguration{
								EncryptionType: "AES256",
							},
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
		{
			name: "should render multiple source registries",
			cfg: `
source:
  registries:
    - type: "aws"
      aws:
        accountId: "12345678912"
        region: "us-west-1"
    - type: "aws"
      aws:
        accountId: "12345678912"
        region: "us-east-1"
`,
			expCfg: Config{
				Target: Registry{
					Type: "aws",
					AWS: AWS{
						ECROptions: ECROptions{
							ImageTagMutability: "MUTABLE",
							ImageScanningConfiguration: ImageScanningConfiguration{
								ImageScanOnPush: true,
							},
							EncryptionConfiguration: EncryptionConfiguration{
								EncryptionType: "AES256",
							},
						},
					},
				},
				Source: Source{
					Registries: []Registry{
						{
							Type: "aws",
							AWS: AWS{
								AccountID: "12345678912",
								Region:    "us-west-1",
							}},
						{
							Type: "aws",
							AWS: AWS{
								AccountID: "12345678912",
								Region:    "us-east-1",
							}},
					},
				},
			},
		},
		{
			name: "should render generic source registry",
			cfg: `
source:
  registries:
    - type: "generic"
      generic:
        name: "dockerio"
        genericOptions:
          domain: "docker.io"
          username: "testuser"
          password: "testpass"
`,
			expCfg: Config{
				Target: Registry{
					Type: "aws",
					AWS: AWS{
						ECROptions: ECROptions{
							ImageTagMutability: "MUTABLE",
							ImageScanningConfiguration: ImageScanningConfiguration{
								ImageScanOnPush: true,
							},
						},
					},
				},
				Source: Source{
					Registries: []Registry{
						{
							Type: "generic",
							Generic: Generic{
								Name: "dockerio",
								GenericOptions: GenericOptions{
									Domain:   "docker.io",
									Username: "testuser",
									Password: "testpass",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "should use previous defaults",
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
				Target: Registry{
					Type: "aws",
					AWS: AWS{
						AccountID: "123456789",
						Region:    "ap-southeast-2",
						Role:      "arn:aws:iam::123456789012:role/roleName",
						ECROptions: ECROptions{
							ImageScanningConfiguration: ImageScanningConfiguration{
								ImageScanOnPush: true,
							},
							EncryptionConfiguration: EncryptionConfiguration{
								EncryptionType: "AES256",
							},
							ImageTagMutability: "MUTABLE",
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

			v := viper.New()
			v.SetConfigType("yaml")
			SetViperDefaults(v)

			readConfigError := v.ReadConfig(strings.NewReader(test.cfg))
			assert.NoError(readConfigError)

			gotCfg := Config{}
			err := v.Unmarshal(&gotCfg)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expCfg, gotCfg)
			}
		})
	}
}
