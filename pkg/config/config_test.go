package config

import (
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

// TestConfigParses validates if yaml annotation do not overlap
const defaultConfig = `
source:
  registries:
    - type: "aws"
      aws:
        accountId: "12345678912"
        region: "us-west-1"
    - type: "generic"
      generic:
        repository: "https://12345678912"
        username: "demo"
        password: "pass"
    - type: "gcp"
      gcp:
        location: "us-east"
        projectId: "12345"
        repositoryId: "67890"
  filters:
    - jmespath: "obj.metadata.namespace == 'kube-system'"
    - jmespath: "obj.metadata.namespace != 'playground'"
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
`

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
			name: "should render multiple source registries",
			cfg: `
source:
  registries:
    - type: "aws"
      aws:
        accountId: "12345678912"
        region: "us-west-1"
    - type: "generic"
      generic:
        repository: "https://12345678912"
        username: "demo"
        password: "pass"
    - type: "aws"
      aws:
        accountId: "12345678912"
        region: "us-east-1"
`,
			expCfg: Config{
				Source: Source{
					Registries: []Registry{
						{
							Type: "aws",
							AWS: AWS{
								AccountID: "12345678912",
								Region:    "us-west-1",
							}},
						{
							Type: "generic",
							GENERIC: GENERIC{
								Repository: "https://12345678912",
								Username:   "demo",
								Password:   "pass",
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

func TestSuccess(t *testing.T) {
	cfg := Config{}
	err := yaml.Unmarshal([]byte(defaultConfig), &cfg)
	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	awsRegistry := cfg.Source.Registries[0]
	awsDomain := awsRegistry.Domain()
	assert.Equal(t, "12345678912.dkr.ecr.us-west-1.amazonaws.com", awsDomain)
	assert.Nil(t, CheckRegistryConfiguration(awsRegistry))

	genericRegistry := cfg.Source.Registries[1]
	genericDomain := genericRegistry.Domain()
	assert.Equal(t, "", genericDomain)
	assert.Nil(t, CheckRegistryConfiguration(genericRegistry))

	gcpRegistry := cfg.Source.Registries[2]
	gcpDomain := gcpRegistry.Domain()
	assert.Equal(t, "us-east-docker.pkg.dev/12345/67890", gcpDomain)
	assert.Nil(t, CheckRegistryConfiguration(gcpRegistry))
}

func TestNoRegistryType(t *testing.T) {
	cfg := Config{}
	err := yaml.Unmarshal([]byte(defaultConfig), &cfg)
	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	registry := cfg.Source.Registries[0]
	registry.Type = ""

	err = CheckRegistryConfiguration(registry)
	assert.NotNil(t, err)
	assert.Equal(t, "a registry requires a type", err.Error())
}

func TestUnknownRegistryType(t *testing.T) {
	cfg := Config{}
	err := yaml.Unmarshal([]byte(defaultConfig), &cfg)
	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	registry := cfg.Source.Registries[0]
	registry.Type = "TEST"

	err = CheckRegistryConfiguration(registry)
	assert.Nil(t, err)
}

func TestAWSRegistryNoRegion(t *testing.T) {

	cfg := Config{}
	err := yaml.Unmarshal([]byte(defaultConfig), &cfg)
	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	awsRegistry := cfg.Source.Registries[0]
	awsRegistry.AWS.Region = ""

	err = CheckRegistryConfiguration(awsRegistry)
	assert.NotNil(t, err)
	assert.Equal(t, "registry of type \"aws\" requires a field region", err.Error())

}

func TestAWSRegistryNoAccount(t *testing.T) {

	cfg := Config{}
	err := yaml.Unmarshal([]byte(defaultConfig), &cfg)
	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	awsRegistry := cfg.Source.Registries[0]
	awsRegistry.AWS.AccountID = ""

	err = CheckRegistryConfiguration(awsRegistry)
	assert.NotNil(t, err)
	assert.Equal(t, "registry of type \"aws\" requires a field \"accountdId\"", err.Error())

}

func TestGenericRegistryNoRepository(t *testing.T) {

	cfg := Config{}
	err := yaml.Unmarshal([]byte(defaultConfig), &cfg)
	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	genericRegistry := cfg.Source.Registries[1]
	genericRegistry.GENERIC.Repository = ""

	err = CheckRegistryConfiguration(genericRegistry)
	assert.NotNil(t, err)
	assert.Equal(t, "registry of type \"generic\" requires a field \"repository\"", err.Error())

}

func TestGenericRegistryNoUsername(t *testing.T) {

	cfg := Config{}
	err := yaml.Unmarshal([]byte(defaultConfig), &cfg)
	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	genericRegistry := cfg.Source.Registries[1]
	genericRegistry.GENERIC.Username = ""

	err = CheckRegistryConfiguration(genericRegistry)
	assert.Nil(t, err)
}

func TestGenericRegistryNoPassword(t *testing.T) {

	cfg := Config{}
	err := yaml.Unmarshal([]byte(defaultConfig), &cfg)
	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	genericRegistry := cfg.Source.Registries[1]
	genericRegistry.GENERIC.Password = ""

	err = CheckRegistryConfiguration(genericRegistry)
	assert.Nil(t, err)
}

func TestGCPRegistryNoLocation(t *testing.T) {

	cfg := Config{}
	err := yaml.Unmarshal([]byte(defaultConfig), &cfg)
	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	gcpRegistry := cfg.Source.Registries[2]
	gcpRegistry.GCP.Location = ""

	err = CheckRegistryConfiguration(gcpRegistry)
	assert.NotNil(t, err)
	assert.Equal(t, "registry of type \"gcp\" requires a field \"location\"", err.Error())

}

func TestGCPRegistryNoRepositoryID(t *testing.T) {

	cfg := Config{}
	err := yaml.Unmarshal([]byte(defaultConfig), &cfg)
	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	gcpRegistry := cfg.Source.Registries[2]
	gcpRegistry.GCP.RepositoryID = ""

	err = CheckRegistryConfiguration(gcpRegistry)
	assert.NotNil(t, err)
	assert.Equal(t, "registry of type \"gcp\" requires a field \"repositoryId\"", err.Error())

}

func TestGCPRegistryNoProjectID(t *testing.T) {

	cfg := Config{}
	err := yaml.Unmarshal([]byte(defaultConfig), &cfg)
	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	gcpRegistry := cfg.Source.Registries[2]
	gcpRegistry.GCP.ProjectID = ""

	err = CheckRegistryConfiguration(gcpRegistry)
	assert.NotNil(t, err)
	assert.Equal(t, "registry of type \"gcp\" requires a field \"projectId\"", err.Error())

}
