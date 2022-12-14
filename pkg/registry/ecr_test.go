package registry

import (
	"encoding/base64"
	"testing"

	"github.com/estahn/k8s-image-swapper/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestDockerConfig(t *testing.T) {
	fakeToken := []byte("token")
	fakeBase64Token := base64.StdEncoding.EncodeToString(fakeToken)

	expected := []byte("{\"auths\":{\"12345678912.dkr.ecr.us-east-1.amazonaws.com\":{\"auth\":\"" + fakeBase64Token + "\"}}}")

	fakeRegistry := NewDummyECRClient("us-east-1", "12345678912", "", config.ECROptions{}, fakeToken)

	r, _ := fakeRegistry.DockerConfig()

	assert.Equal(t, r, expected)
}
