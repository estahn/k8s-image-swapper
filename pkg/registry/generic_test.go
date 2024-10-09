package registry

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/containers/image/v5/transports/alltransports"
	"github.com/estahn/k8s-image-swapper/pkg/config"
	"github.com/stretchr/testify/assert"
)

var defaultGenericCfg = config.Generic{
	Repository: "localhost",
	Username:   "user",
	Password:   "password",
	IgnoreCert: true,
}

func createGenericClient(config config.Generic, testName string) (*GenericClient, error) {

	curCommandExecutor := commandExecutor
	defer func() { commandExecutor = curCommandExecutor }()

	commandExecutor = func(ctx context.Context, name string, arg ...string) ShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v\n", testName, name, arg)
		return testCommandExecutor{
			output: []byte("login successful"),
			err:    nil,
		}
	}
	return NewGenericClient(config)
}

func TestNewGenericClientSuccess(t *testing.T) {

	genericClient, err := createGenericClient(defaultGenericCfg, t.Name())
	assert.Nil(t, err)
	assert.NotNil(t, genericClient)
	assert.Equal(t, "localhost", genericClient.repository)
	assert.Equal(t, "user", genericClient.username)
	assert.Equal(t, "password", genericClient.password)
}

func TestLoginFailure(t *testing.T) {

	curCommandExecutor := commandExecutor
	defer func() { commandExecutor = curCommandExecutor }()

	commandExecutor = func(ctx context.Context, name string, arg ...string) ShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v\n", t.Name(), name, arg)
		return testCommandExecutor{
			output: []byte(""),
			err:    errors.New("login failure"),
		}
	}

	_, err := NewGenericClient(defaultGenericCfg)
	assert.NotNil(t, err)
}

func TestImageExistsSuccess(t *testing.T) {

	genericClient, err := createGenericClient(defaultGenericCfg, t.Name())
	assert.Nil(t, err)
	assert.NotNil(t, genericClient)

	curCommandExecutor := commandExecutor
	defer func() { commandExecutor = curCommandExecutor }()

	commandExecutor = func(ctx context.Context, name string, arg ...string) ShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v\n", t.Name(), name, arg)
		return testCommandExecutor{
			output: []byte(""),
			err:    nil,
		}
	}

	imageRef, _ := alltransports.ParseImageName("docker://k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	exists := genericClient.ImageExists(context.Background(), imageRef)
	assert.Equal(t, true, exists)
}

func TestImageExistsInCacheSuccess(t *testing.T) {

	genericClient, err := createGenericClient(defaultGenericCfg, t.Name())
	assert.Nil(t, err)
	assert.NotNil(t, genericClient)

	imageRef, _ := alltransports.ParseImageName("docker://k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	genericClient.cache.Set(imageRef.DockerReference().String(), "123", 123)
	time.Sleep(time.Millisecond * 10)
	exists := genericClient.ImageExists(context.Background(), imageRef)
	assert.Equal(t, true, exists)
}

func TestImageDoesNotExistsSuccess(t *testing.T) {

	genericClient, err := createGenericClient(defaultGenericCfg, t.Name())
	assert.Nil(t, err)
	assert.NotNil(t, genericClient)

	curCommandExecutor := commandExecutor
	defer func() { commandExecutor = curCommandExecutor }()

	commandExecutor = func(ctx context.Context, name string, arg ...string) ShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v\n", t.Name(), name, arg)
		return testCommandExecutor{
			err: errors.New("Image not found"),
		}
	}

	imageRef, _ := alltransports.ParseImageName("docker://k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	exists := genericClient.ImageExists(context.Background(), imageRef)
	assert.Equal(t, false, exists)
}

func TestCopyImageSuccess(t *testing.T) {

	genericClient, err := createGenericClient(defaultGenericCfg, t.Name())
	assert.Nil(t, err)
	assert.NotNil(t, genericClient)

	curCommandExecutor := commandExecutor
	defer func() { commandExecutor = curCommandExecutor }()

	commandExecutor = func(ctx context.Context, name string, arg ...string) ShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v\n", t.Name(), name, arg)
		return testCommandExecutor{
			output: []byte(""),
			err:    nil,
		}
	}

	srcRef, _ := alltransports.ParseImageName("docker://k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")
	destRef, _ := alltransports.ParseImageName("docker://k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	srcCreds := "user:pass"
	destCreds := "user:pass"

	err = genericClient.CopyImage(context.Background(), srcRef, srcCreds, destRef, destCreds)
	assert.Nil(t, err)
}

func TestCopyImageSuccessNoCredentials(t *testing.T) {

	genericClient, err := createGenericClient(defaultGenericCfg, t.Name())
	assert.Nil(t, err)
	assert.NotNil(t, genericClient)

	curCommandExecutor := commandExecutor
	defer func() { commandExecutor = curCommandExecutor }()

	commandExecutor = func(ctx context.Context, name string, arg ...string) ShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v\n", t.Name(), name, arg)
		return testCommandExecutor{
			output: []byte(""),
			err:    nil,
		}
	}

	srcRef, _ := alltransports.ParseImageName("docker://k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")
	destRef, _ := alltransports.ParseImageName("docker://k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	srcCreds := ""
	destCreds := ""

	err = genericClient.CopyImage(context.Background(), srcRef, srcCreds, destRef, destCreds)
	assert.Nil(t, err)
}

func TestCopyImageFailure(t *testing.T) {

	genericClient, err := createGenericClient(defaultGenericCfg, t.Name())
	assert.Nil(t, err)
	assert.NotNil(t, genericClient)

	curCommandExecutor := commandExecutor
	defer func() { commandExecutor = curCommandExecutor }()

	commandExecutor = func(ctx context.Context, name string, arg ...string) ShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v\n", t.Name(), name, arg)
		return testCommandExecutor{
			output: []byte("command not found"),
			err:    errors.New("copy error"),
		}
	}

	srcRef, _ := alltransports.ParseImageName("docker://k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")
	destRef, _ := alltransports.ParseImageName("docker://k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	srcCreds := "user:pass"
	destCreds := "user:pass"

	err = genericClient.CopyImage(context.Background(), srcRef, srcCreds, destRef, destCreds)
	assert.NotNil(t, err)
	assert.Equal(t, "Command error, stderr: copy error, stdout: command not found", err.Error())
}

func TestDockerConfigSuccess(t *testing.T) {
	genericClient, err := createGenericClient(defaultGenericCfg, t.Name())
	assert.Nil(t, err)
	assert.NotNil(t, genericClient)

	data, err := genericClient.DockerConfig()
	assert.Nil(t, err)
	assert.NotNil(t, data)

	dockerConfig := &DockerConfig{}
	err = json.Unmarshal(data, dockerConfig)
	assert.Nil(t, err)

	for key, authConfig := range dockerConfig.AuthConfigs {
		assert.Equal(t, "localhost", key)
		assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("user:password")), authConfig.Auth)
	}
}

func TestIsOrigin(t *testing.T) {
	genericClient, err := createGenericClient(defaultGenericCfg, t.Name())
	assert.Nil(t, err)
	assert.NotNil(t, genericClient)

	imageRef, _ := alltransports.ParseImageName("docker://localhost/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	isOrigin := genericClient.IsOrigin(imageRef)
	assert.True(t, isOrigin)
}

func TestIsNotOrigin(t *testing.T) {
	genericClient, err := createGenericClient(defaultGenericCfg, t.Name())
	assert.Nil(t, err)
	assert.NotNil(t, genericClient)

	imageRef, _ := alltransports.ParseImageName("docker://k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	isOrigin := genericClient.IsOrigin(imageRef)
	assert.False(t, isOrigin)
}

func TestCreateRegistryNilResponse(t *testing.T) {

	genericClient, err := createGenericClient(defaultGenericCfg, t.Name())
	assert.Nil(t, err)
	assert.NotNil(t, genericClient)

	err = genericClient.CreateRepository(context.Background(), "repo")
	assert.Nil(t, err)
}
