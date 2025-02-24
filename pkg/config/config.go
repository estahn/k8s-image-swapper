/*
Copyright © 2020 Enrico Stahn <enrico.stahn@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"

	"github.com/estahn/k8s-image-swapper/pkg/types"
)

const DefaultImageCopyDeadline = 8 * time.Second

type Config struct {
	LogLevel  string `yaml:"logLevel" validate:"oneof=trace debug info warn error fatal"`
	LogFormat string `yaml:"logFormat" validate:"oneof=json console"`

	ListenAddress string

	DryRun            bool          `yaml:"dryRun"`
	ImageSwapPolicy   string        `yaml:"imageSwapPolicy" validate:"oneof=always exists"`
	ImageCopyPolicy   string        `yaml:"imageCopyPolicy" validate:"oneof=delayed immediate force none"`
	ImageCopyDeadline time.Duration `yaml:"imageCopyDeadline"`

	Source Source   `yaml:"source"`
	Target Registry `yaml:"target"`

	TLSCertFile string
	TLSKeyFile  string
}

type JMESPathFilter struct {
	JMESPath string `yaml:"jmespath"`
}

type Source struct {
	Registries []Registry       `yaml:"registries"`
	Filters    []JMESPathFilter `yaml:"filters"`
}

type Registry struct {
	Type    string  `yaml:"type"`
	AWS     AWS     `yaml:"aws"`
	GCP     GCP     `yaml:"gcp"`
	Generic Generic `yaml:"generic"`
}

type Generic struct {
	Name           string         `yaml:"name"`
	GenericOptions GenericOptions `yaml:"genericOptions"`
}

type GenericOptions struct {
	Domain   string `yaml:"domain"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type AWS struct {
	AccountID  string     `yaml:"accountId"`
	Region     string     `yaml:"region"`
	Role       string     `yaml:"role"`
	ECROptions ECROptions `yaml:"ecrOptions"`
}

type GCP struct {
	Location     string `yaml:"location"`
	ProjectID    string `yaml:"projectId"`
	RepositoryID string `yaml:"repositoryId"`
}

type ECROptions struct {
	AccessPolicy               string                     `yaml:"accessPolicy"`
	LifecyclePolicy            string                     `yaml:"lifecyclePolicy"`
	Tags                       []Tag                      `yaml:"tags"`
	ImageTagMutability         string                     `yaml:"imageTagMutability"  validate:"oneof=MUTABLE IMMUTABLE"`
	ImageScanningConfiguration ImageScanningConfiguration `yaml:"imageScanningConfiguration"`
	EncryptionConfiguration    EncryptionConfiguration    `yaml:"encryptionConfiguration"`
}

type Tag struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type ImageScanningConfiguration struct {
	ImageScanOnPush bool `yaml:"imageScanOnPush"`
}

type EncryptionConfiguration struct {
	EncryptionType string `yaml:"encryptionType" validate:"oneof=KMS AES256"`
	KmsKey         string `yaml:"kmsKey"`
}

func (a *AWS) EcrDomain() string {
	return fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", a.AccountID, a.Region)
}

func (g *GCP) GarDomain() string {
	return fmt.Sprintf("%s-docker.pkg.dev/%s/%s", g.Location, g.ProjectID, g.RepositoryID)
}

func (g *Generic) GenericDomain() string {
	return g.GenericOptions.Domain
}

func (r Registry) Domain() string {
	registry, _ := types.ParseRegistry(r.Type)
	switch registry {
	case types.RegistryAWS:
		return r.AWS.EcrDomain()
	case types.RegistryGCP:
		return r.GCP.GarDomain()
	case types.RegistryGeneric:
		return r.Generic.GenericDomain()
	default:
		return ""
	}
}

// provides detailed information about wrongly provided configuration
func CheckRegistryConfiguration(r Registry) error {
	if r.Type == "" {
		return fmt.Errorf("a registry requires a type")
	}

	errorWithType := func(info string) error {
		return fmt.Errorf(`registry of type "%s" %s`, r.Type, info)
	}

	registry, _ := types.ParseRegistry(r.Type)
	switch registry {
	case types.RegistryAWS:
		if r.AWS.Region == "" {
			return errorWithType(`requires a field "region"`)
		}
		if r.AWS.AccountID == "" {
			return errorWithType(`requires a field "accountdId"`)
		}
		if r.AWS.ECROptions.EncryptionConfiguration.EncryptionType == "KMS" && r.AWS.ECROptions.EncryptionConfiguration.KmsKey == "" {
			return errorWithType(`requires a field "kmsKey" if encryptionType is set to "KMS"`)
		}
	case types.RegistryGCP:
		if r.GCP.Location == "" {
			return errorWithType(`requires a field "location"`)
		}
		if r.GCP.ProjectID == "" {
			return errorWithType(`requires a field "projectId"`)
		}
		if r.GCP.RepositoryID == "" {
			return errorWithType(`requires a field "repositoryId"`)
		}
	}

	return nil
}

// provides detailed information about wrongly provided configuration (target specific)
func CheckTargetRegistryConfiguration(r Registry) error {
	registryType, err := types.ParseRegistry(r.Type)
	if err != nil {
		return fmt.Errorf("couldn't parse target registry type")
	} else if registryType == types.RegistryGeneric {
		return fmt.Errorf("generic registry not allowed as target: %s", r.Generic.Name)
	}
	return nil
}

// SetViperDefaults configures default values for config items that are not set.
func SetViperDefaults(v *viper.Viper) {
	v.SetDefault("Target.Type", "aws")
	v.SetDefault("Target.AWS.ECROptions.ImageScanningConfiguration.ImageScanOnPush", true)
	v.SetDefault("Target.AWS.ECROptions.ImageTagMutability", "MUTABLE")
	v.SetDefault("Target.AWS.ECROptions.EncryptionConfiguration.EncryptionType", "AES256")
}
