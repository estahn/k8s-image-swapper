package registry

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/containers/image/v5/transports/alltransports"
	"github.com/dgraph-io/ristretto"

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

func TestEcrImageExistsCaching(t *testing.T) {
	// Setup a test cache
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	assert.NoError(t, err)

	tests := []struct {
		name            string
		cacheTtlMinutes int
		expectCached    bool
	}{
		{
			name:            "cache disabled when TTL is 0",
			cacheTtlMinutes: 0,
			expectCached:    false,
		},
		{
			name:            "cache enabled with TTL and jitter",
			cacheTtlMinutes: 60,
			expectCached:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			client := NewDummyECRClient("us-east-1", "12345678912", "", config.ECROptions{}, []byte(""))

			// Setup cache
			client.cache = cache
			client.cacheTtlMinutes = tc.cacheTtlMinutes

			// Create a test image reference and add to cache. Use 100ms as TTL
			imageRef, err := alltransports.ParseImageName("docker://12345678912.dkr.ecr.us-east-1.amazonaws.com/test-project/repo/test-image:latest")
			cache.SetWithTTL(imageRef.DockerReference().String(), true, 1, 100*time.Millisecond)
			assert.NoError(t, err)

			// Cache should be a hit
			exists := client.ImageExists(ctx, imageRef)
			assert.Equal(t, tc.expectCached, exists)

			if tc.expectCached {
				// Verify cache expiry
				time.Sleep(time.Duration(150 * time.Millisecond)) // Use milliseconds for testing
				_, found := client.cache.Get(imageRef.DockerReference().String())
				assert.False(t, found, "cache entry should have expired")
			}
		})
	}
}
