package registry

import (
	"encoding/base64"
	"testing"

	"github.com/estahn/k8s-image-swapper/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestGenericDockerConfig(t *testing.T) {
	fakeToken := []byte("username:password")
	fakeBase64Token := base64.StdEncoding.EncodeToString(fakeToken)

	expected := []byte("{\"auths\":{\"docker.io\":{\"auth\":\"" + fakeBase64Token + "\"}}}")

	fakeRegistry := NewDummyGenericClient("docker.io", config.GenericOptions{
		Domain:   "docker.io",
		Username: "username",
		Password: "password",
	})

	r, _ := GenerateDockerConfig(fakeRegistry)

	assert.Equal(t, r, expected)
}
