package registry

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/containers/image/v5/transports/alltransports"
	"github.com/estahn/k8s-image-swapper/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestDockerConfig(t *testing.T) {
	fakeToken := []byte("token")
	fakeBase64Token := base64.StdEncoding.EncodeToString(fakeToken)

	expected := []byte("{\"auths\":{\"12345678912.dkr.ecr.us-east-1.amazonaws.com\":{\"auth\":\"" + fakeBase64Token + "\"}}}")

	fakeRegistry := NewDummyECRClient("us-east-1", "12345678912", "", config.ECROptions{}, fakeToken)

	r, _ := GenerateDockerConfig(fakeRegistry)

	assert.Equal(t, r, expected)
}

func TestECRIsOrigin(t *testing.T) {
	type testCase struct {
		input    string
		expected bool
	}
	testcases := []testCase{
		{
			input:    "k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713",
			expected: false,
		},
		{
			input:    "12345678912.dkr.ecr.us-east-1.amazonaws.com/k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713",
			expected: true,
		},
	}

	fakeRegistry := NewDummyECRClient("us-east-1", "12345678912", "", config.ECROptions{}, []byte(""))

	for _, testcase := range testcases {
		imageRef, err := alltransports.ParseImageName("docker://" + testcase.input)

		assert.NoError(t, err)

		result := fakeRegistry.IsOrigin(imageRef)

		assert.Equal(t, testcase.expected, result)
	}
}

func TestECRClientCopyImageSuccess(t *testing.T) {

	ecrClient := NewDummyECRClient("us-east-1", "12345678912", "", config.ECROptions{}, []byte(""))

	curCommandExecutor := commandExecutor
	defer func() { commandExecutor = curCommandExecutor }()

	commandExecutor = func(ctx context.Context, name string, arg ...string) ShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v\n", t.Name(), name, arg)
		return testCommandExecutor{
			output: []byte(""),
			err:    nil,
		}
	}

	srcRef, _ := alltransports.ParseImageName("docker://12345678912.dkr.ecr.us-east-1.amazonaws.com/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")
	destRef, _ := alltransports.ParseImageName("docker://12345678912.dkr.ecr.us-east-1.amazonaws.com/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	srcCreds := "user:pass"
	destCreds := "user:pass"

	err := ecrClient.CopyImage(context.Background(), srcRef, srcCreds, destRef, destCreds)
	assert.Nil(t, err)

}

func TestECRClientCopyImageSuccessNoCreds(t *testing.T) {

	ecrClient := NewDummyECRClient("us-east-1", "12345678912", "", config.ECROptions{}, []byte(""))

	curCommandExecutor := commandExecutor
	defer func() { commandExecutor = curCommandExecutor }()

	commandExecutor = func(ctx context.Context, name string, arg ...string) ShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v\n", t.Name(), name, arg)
		return testCommandExecutor{
			output: []byte(""),
			err:    nil,
		}
	}

	srcRef, _ := alltransports.ParseImageName("docker://12345678912.dkr.ecr.us-east-1.amazonaws.com/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")
	destRef, _ := alltransports.ParseImageName("docker://12345678912.dkr.ecr.us-east-1.amazonaws.com/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	srcCreds := ""
	destCreds := ""

	err := ecrClient.CopyImage(context.Background(), srcRef, srcCreds, destRef, destCreds)
	assert.Nil(t, err)

}

func TestECRClientCopyImageFailure(t *testing.T) {

	ecrClient := NewDummyECRClient("us-east-1", "12345678912", "", config.ECROptions{}, []byte(""))

	curCommandExecutor := commandExecutor
	defer func() { commandExecutor = curCommandExecutor }()

	commandExecutor = func(ctx context.Context, name string, arg ...string) ShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v\n", t.Name(), name, arg)
		return testCommandExecutor{
			output: []byte("missing"),
			err:    errors.New("Command Failed"),
		}
	}

	srcRef, _ := alltransports.ParseImageName("docker://12345678912.dkr.ecr.us-east-1.amazonaws.com/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")
	destRef, _ := alltransports.ParseImageName("docker://12345678912.dkr.ecr.us-east-1.amazonaws.com/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	srcCreds := ""
	destCreds := ""

	err := ecrClient.CopyImage(context.Background(), srcRef, srcCreds, destRef, destCreds)
	assert.NotNil(t, err)
	assert.Equal(t, "Command error, stderr: Command Failed, stdout: missing", err.Error())

}

func TestECRClientIsNotOrigin(t *testing.T) {
	ecrClient := NewDummyECRClient("us-east-1", "12345678912", "", config.ECROptions{}, []byte(""))

	imageRef, _ := alltransports.ParseImageName("docker://test-ecr.us-east-1.amazonaws.com/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	isOrigin := ecrClient.IsOrigin(imageRef)
	assert.False(t, isOrigin)
}

func TestECRClientIsOrigin(t *testing.T) {
	ecrClient := NewDummyECRClient("us-east-1", "12345678912", "", config.ECROptions{}, []byte(""))

	imageRef, _ := alltransports.ParseImageName("docker://12345678912.dkr.ecr.us-east-1.amazonaws.com/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	isOrigin := ecrClient.IsOrigin(imageRef)
	assert.True(t, isOrigin)
}

func TestECRClientImageExistsSuccess(t *testing.T) {

	ecrClient := NewDummyECRClient("us-east-1", "12345678912", "", config.ECROptions{}, []byte(""))

	curCommandExecutor := commandExecutor
	defer func() { commandExecutor = curCommandExecutor }()

	commandExecutor = func(ctx context.Context, name string, arg ...string) ShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v\n", t.Name(), name, arg)
		return testCommandExecutor{
			output: []byte(""),
			err:    nil,
		}
	}

	imageRef, _ := alltransports.ParseImageName("docker://12345678912.dkr.ecr.us-east-1.amazonaws.com/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	exists := ecrClient.ImageExists(context.Background(), imageRef)
	assert.Equal(t, true, exists)
}

func TestECRClientImageExistsInCacheSuccess(t *testing.T) {

	ecrClient := NewDummyECRClient("us-east-1", "12345678912", "", config.ECROptions{}, []byte(""))

	imageRef, _ := alltransports.ParseImageName("docker://12345678912.dkr.ecr.us-east-1.amazonaws.com/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")
	ecrClient.cache.Set(imageRef.DockerReference().String(), "123", 123)
	time.Sleep(time.Millisecond * 10)

	exists := ecrClient.ImageExists(context.Background(), imageRef)
	assert.Equal(t, true, exists)
}

func TestECRClientImageDoesNotExistsSuccess(t *testing.T) {

	ecrClient := NewDummyECRClient("us-east-1", "12345678912", "", config.ECROptions{}, []byte(""))

	curCommandExecutor := commandExecutor
	defer func() { commandExecutor = curCommandExecutor }()

	commandExecutor = func(ctx context.Context, name string, arg ...string) ShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v\n", t.Name(), name, arg)
		return testCommandExecutor{
			err: errors.New("Image not found"),
		}
	}

	imageRef, _ := alltransports.ParseImageName("docker://12345678912.dkr.ecr.us-east-1.amazonaws.com/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	exists := ecrClient.ImageExists(context.Background(), imageRef)
	assert.Equal(t, false, exists)
}

func TestInitClientSuccess(t *testing.T) {

	config := config.AWS{
		AccountID:  "123",
		Region:     "us-east-1",
		Role:       "",
		ECROptions: config.ECROptions{},
	}

	client := initClient(config)

	assert.NotNil(t, client)
	assert.Equal(t, "123.dkr.ecr.us-east-1.amazonaws.com", client.Endpoint())
	assert.Equal(t, "123", client.targetAccount)
}
func TestInitClientWithRoleSuccess(t *testing.T) {

	config := config.AWS{
		AccountID:  "123",
		Region:     "us-east-1",
		Role:       "admin",
		ECROptions: config.ECROptions{},
	}

	client := initClient(config)

	assert.NotNil(t, client)
	assert.Equal(t, "123.dkr.ecr.us-east-1.amazonaws.com", client.Endpoint())
	assert.Equal(t, "123", client.targetAccount)
}

func TestCreateRepositoryInCache(t *testing.T) {
	ecrClient := NewDummyECRClient("us-east-1", "12345678912", "", config.ECROptions{}, []byte(""))
	ecrClient.cache.Set("reg1", "123", 1)
	time.Sleep(time.Millisecond * 10)
	err := ecrClient.CreateRepository(context.Background(), "reg1")
	assert.Nil(t, err)

}

func TestBuildEcrTagsSuccess(t *testing.T) {
	registryClient, _ := NewMockECRClient(nil, "ap-southeast-2", "123456789.dkr.ecr.ap-southeast-2.amazonaws.com", "123456789", "arn:aws:iam::123456789:role/fakerole")
	assert.NotNil(t, registryClient)
	ecrTags := registryClient.buildEcrTags()
	assert.NotNil(t, ecrTags)
	assert.Len(t, ecrTags, 2)
	assert.Equal(t, "CreatedBy", *ecrTags[0].Key)
	assert.Equal(t, "k8s-image-swapper", *ecrTags[0].Value)
	assert.Equal(t, "AnotherTag", *ecrTags[1].Key)
	assert.Equal(t, "another-tag", *ecrTags[1].Value)
}
