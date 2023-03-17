package registry

import (
	"encoding/base64"
	"github.com/containers/image/v5/transports/alltransports"
	"testing"

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
