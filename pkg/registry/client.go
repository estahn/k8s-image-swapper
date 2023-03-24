package registry

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/ristretto"
	"github.com/estahn/k8s-image-swapper/pkg/backend"
	"github.com/estahn/k8s-image-swapper/pkg/config"
	"github.com/estahn/k8s-image-swapper/pkg/types"

	ctypes "github.com/containers/image/v5/types"
)

// Client provides methods required to be implemented by the various target registry clients, e.g. ECR, Docker, Quay.
type Client interface {
	CreateRepository(ctx context.Context, name string) error
	RepositoryExists() bool
	CopyImage(ctx context.Context, src ctypes.ImageReference, srcCreds string, dest ctypes.ImageReference, destCreds string) error
	ImageExists(ctx context.Context, ref ctypes.ImageReference) bool

	// Endpoint returns the domain of the registry
	Endpoint() string
	Credentials() string

	// IsOrigin returns true if the imageRef originates from this registry
	IsOrigin(imageRef ctypes.ImageReference) bool
}

type DockerConfig struct {
	AuthConfigs map[string]AuthConfig `json:"auths"`
}

type AuthConfig struct {
	Auth string `json:"auth,omitempty"`
}

// NewClient returns a registry client ready for use without the need to specify an implementation
func NewClient(r config.Registry, imageBackend backend.Backend) (Client, error) {
	if err := config.CheckRegistryConfiguration(r); err != nil {
		return nil, err
	}

	registry, err := types.ParseRegistry(r.Type)
	if err != nil {
		return nil, err
	}

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		return nil, err
	}

	cachedBackend := backend.NewCached(cache, imageBackend)

	switch registry {
	case types.RegistryAWS:
		return NewECRClient(r.AWS, cachedBackend, cache)
	case types.RegistryGCP:
		return NewGARClient(r.GCP, cachedBackend)
	default:
		return nil, fmt.Errorf(`registry of type "%s" is not supported`, r.Type)
	}
}

func GenerateDockerConfig(c Client) ([]byte, error) {
	dockerConfig := DockerConfig{
		AuthConfigs: map[string]AuthConfig{
			c.Endpoint(): {
				Auth: base64.StdEncoding.EncodeToString([]byte(c.Credentials())),
			},
		},
	}

	dockerConfigJson, err := json.Marshal(dockerConfig)
	if err != nil {
		return []byte{}, err
	}

	return dockerConfigJson, nil
}
