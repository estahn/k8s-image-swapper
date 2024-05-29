package registry

import (
	"context"
	"fmt"

	ctypes "github.com/containers/image/v5/types"
	"github.com/estahn/k8s-image-swapper/pkg/config"
)

type GenericClient struct {
	options config.GenericOptions
}

func NewGenericClient(clientConfig config.Generic) (*GenericClient, error) {
	client := GenericClient{}

	client.options = clientConfig.GenericOptions

	return &client, nil
}

func (g *GenericClient) CreateRepository(ctx context.Context, name string) error {
	return nil
}

func (g *GenericClient) RepositoryExists() bool {
	return true
}

func (g *GenericClient) CopyImage(ctx context.Context, src ctypes.ImageReference, srcCreds string, dest ctypes.ImageReference, destCreds string) error {
	panic("implement me")
}

func (g *GenericClient) PullImage() error {
	panic("implement me")
}

func (g *GenericClient) PutImage() error {
	panic("implement me")
}

func (g *GenericClient) ImageExists(ctx context.Context, ref ctypes.ImageReference) bool {
	return true
}

// Endpoint returns the domain of the registry
func (g *GenericClient) Endpoint() string {
	return g.options.Domain
}

func (g *GenericClient) Credentials() string {
	return fmt.Sprintf("%s:%s", g.options.Username, g.options.Password)
}

// IsOrigin returns true if the imageRef originates from this registry
func (g *GenericClient) IsOrigin(imageRef ctypes.ImageReference) bool {
	return true
}

// For testing purposes
func NewDummyGenericClient(domain string, options config.GenericOptions) *GenericClient {
	return &GenericClient{
		options: options,
	}
}
