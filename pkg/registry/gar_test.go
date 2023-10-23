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

func TestGARIsOrigin(t *testing.T) {
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
			input:    "us-central1-docker.pkg.dev/gcp-project-123/main/k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713",
			expected: true,
		},
	}

	fakeRegistry, _ := NewMockGARClient(nil, "us-central1-docker.pkg.dev/gcp-project-123/main")

	for _, testcase := range testcases {
		imageRef, err := alltransports.ParseImageName("docker://" + testcase.input)

		assert.NoError(t, err)

		result := fakeRegistry.IsOrigin(imageRef)

		assert.Equal(t, testcase.expected, result)
	}
}

func TestGARClientCopyImageSuccess(t *testing.T) {

	garClient, _ := NewMockGARClient(nil, "us-central1-docker.pkg.dev/gcp-project-123/main")

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

	err := garClient.CopyImage(context.Background(), srcRef, srcCreds, destRef, destCreds)
	assert.Nil(t, err)

}

func TestGARClientCopyImageWithSuffixSuccess(t *testing.T) {

	garClient, _ := NewMockGARClient(nil, "us-central1-docker.pkg.dev/gcp-project-123/main")

	curCommandExecutor := commandExecutor
	defer func() { commandExecutor = curCommandExecutor }()

	commandExecutor = func(ctx context.Context, name string, arg ...string) ShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v\n", t.Name(), name, arg)
		return testCommandExecutor{
			output: []byte(""),
			err:    nil,
		}
	}

	srcRef, _ := alltransports.ParseImageName("docker://12345678912-docker.pkg.dev/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")
	destRef, _ := alltransports.ParseImageName("docker://12345678912.dkr.ecr.us-east-1.amazonaws.com/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	srcCreds := "user:pass"
	destCreds := "user:pass"

	err := garClient.CopyImage(context.Background(), srcRef, srcCreds, destRef, destCreds)
	assert.Nil(t, err)

}

func TestGARClientCopyImageSuccessNoCreds(t *testing.T) {

	garClient, _ := NewMockGARClient(nil, "us-central1-docker.pkg.dev/gcp-project-123/main")

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

	err := garClient.CopyImage(context.Background(), srcRef, srcCreds, destRef, destCreds)
	assert.Nil(t, err)

}

func TestGARClientCopyImageFailure(t *testing.T) {

	garClient, _ := NewMockGARClient(nil, "us-central1-docker.pkg.dev/gcp-project-123/main")

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

	err := garClient.CopyImage(context.Background(), srcRef, srcCreds, destRef, destCreds)
	assert.NotNil(t, err)
	assert.Equal(t, "Command error, stderr: Command Failed, stdout: missing", err.Error())

}

func TestGARClientIsNotOrigin(t *testing.T) {
	garClient, _ := NewMockGARClient(nil, "us-central1-docker.pkg.dev/gcp-project-123/main")

	imageRef, _ := alltransports.ParseImageName("docker://test-ecr.us-east-1.amazonaws.com/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	isOrigin := garClient.IsOrigin(imageRef)
	assert.False(t, isOrigin)
}

func TestGARClientIsOrigin(t *testing.T) {
	garClient, _ := NewMockGARClient(nil, "us-central1-docker.pkg.dev/gcp-project-123/main")

	imageRef, _ := alltransports.ParseImageName("docker://us-central1-docker.pkg.dev/gcp-project-123/main/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	isOrigin := garClient.IsOrigin(imageRef)
	assert.True(t, isOrigin)
}

func TestGARClientImageExistsSuccess(t *testing.T) {

	garClient, _ := NewMockGARClient(nil, "us-central1-docker.pkg.dev/gcp-project-123/main")

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

	exists := garClient.ImageExists(context.Background(), imageRef)
	assert.Equal(t, true, exists)
}

func TestGARClientImageExistsInCacheSuccess(t *testing.T) {

	garClient, _ := NewMockGARClient(nil, "us-central1-docker.pkg.dev/gcp-project-123/main")

	imageRef, _ := alltransports.ParseImageName("docker://us-central1-docker.pkg.dev/gcp-project-123/main/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")
	garClient.cache.Set(imageRef.DockerReference().String(), "123", 123)
	time.Sleep(time.Millisecond * 10)

	exists := garClient.ImageExists(context.Background(), imageRef)
	assert.Equal(t, true, exists)
}

func TestGARClientImageDoesNotExistsSuccess(t *testing.T) {

	garClient, _ := NewMockGARClient(nil, "us-central1-docker.pkg.dev/gcp-project-123/main")

	curCommandExecutor := commandExecutor
	defer func() { commandExecutor = curCommandExecutor }()

	commandExecutor = func(ctx context.Context, name string, arg ...string) ShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v\n", t.Name(), name, arg)
		return testCommandExecutor{
			err: errors.New("Image not found"),
		}
	}

	imageRef, _ := alltransports.ParseImageName("docker://12345678912.dkr.ecr.us-east-1.amazonaws.com/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713")

	exists := garClient.ImageExists(context.Background(), imageRef)
	assert.Equal(t, false, exists)
}

func TestGARClientDockerConfigSuccess(t *testing.T) {
	garClient, _ := NewMockGARClient(nil, "us-central1-docker.pkg.dev/gcp-project-123/main")

	data, err := garClient.DockerConfig()
	assert.Nil(t, err)
	assert.NotNil(t, data)

	dockerConfig := &DockerConfig{}
	err = json.Unmarshal(data, dockerConfig)
	assert.Nil(t, err)

	for key, authConfig := range dockerConfig.AuthConfigs {
		assert.Equal(t, "us-central1-docker.pkg.dev/gcp-project-123/main", key)
		assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("oauth2accesstoken:mock-gar-client-fake-auth-token")), authConfig.Auth)
	}
}
func TestGARClientCreateRegistryNilResponse(t *testing.T) {

	garClient, _ := NewMockGARClient(nil, "us-central1-docker.pkg.dev/gcp-project-123/main")
	err := garClient.CreateRepository(context.Background(), "repo")
	assert.Nil(t, err)
}

func TestGARClientScheduleTokenRenewalFailure(t *testing.T) {

	garClient, _ := NewMockGARClient(nil, "us-central1-docker.pkg.dev/gcp-project-123/main")
	err := garClient.scheduleTokenRenewal()
	assert.NotNil(t, err)
}

func TestGARClientNewClientFailure(t *testing.T) {

	cfg := config.GCP{
		Location:     "us-central1-docker.pkg.dev/gcp-project-123/main",
		ProjectID:    "123",
		RepositoryID: "456",
	}

	garClient, err := NewGARClient(cfg)
	assert.NotNil(t, err)
	assert.Nil(t, garClient)
}
