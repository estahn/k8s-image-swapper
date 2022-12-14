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
)

const DefaultImageCopyDeadline = 8 * time.Second

type Config struct {
	LogLevel  string `yaml:"logLevel" validate:"oneof=trace debug info warn error fatal"`
	LogFormat string `yaml:"logFormat" validate:"oneof=json console"`

	ListenAddress string

	DryRun            bool          `yaml:"dryRun"`
	ImageSwapPolicy   string        `yaml:"imageSwapPolicy" validate:"oneof=always exists"`
	ImageCopyPolicy   string        `yaml:"imageCopyPolicy" validate:"oneof=delayed immediate force"`
	ImageCopyDeadline time.Duration `yaml:"imageCopyDeadline"`

	Source Source `yaml:"source"`
	Target Target `yaml:"target"`

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

type Target struct {
	Registry Registry `yaml:"registry"`
}

type RegistryType string

type Registry struct {
	Type RegistryType `yaml:"type"`
	AWS  AWS          `yaml:"aws,omitempty"`
	// TODO add other possible types of registry
	// example:
	// DockerIO  DockerIO    `yaml:"dockerio,omitempty"`
}

type AWS struct {
	AccountID  string     `yaml:"accountId"`
	Region     string     `yaml:"region"`
	Role       string     `yaml:"role"`
	ECROptions ECROptions `yaml:"ecrOptions"`
}

type ECROptions struct {
	AccessPolicy               string                     `yaml:"accessPolicy"`
	LifecyclePolicy            string                     `yaml:"lifecyclePolicy"`
	Tags                       []Tag                      `yaml:"tags"`
	ImageTagMutability         string                     `yaml:"imageTagMutability"`
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
	EncryptionType string `yaml:"encryptionType"`
	KmsKey         string `yaml:"kmsKey"`
}

// TODO add additional structs for different types of registries

// enum for the supported registry types
const (
	Aws RegistryType = "aws"
	// TODO add other possible types of registry
	// example:
	// DockerIO RegistryType = "dockerio"
)

// provides detailed information about wrongly provided configuration
func (r Registry) ValidateConfiguration() error {
	switch r.Type {
	case "":
		return fmt.Errorf("a registry requires a type")
	case Aws:
		if r.AWS.Region == "" {
			return fmt.Errorf(`registry of type "%s" requires a field "region"`, r.Type)
		}
		if r.AWS.AccountID == "" {
			return fmt.Errorf(`registry of type "%s" requires a field "accountdId"`, r.Type)
		}
	}
	return nil
}

func (r Registry) GetServerAddress() string {
	switch r.Type {
	case Aws:
		aws := r.AWS
		return fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", aws.AccountID, aws.Region)
	default:
		return ""
	}
}
