package registry

import (
	"context"

	ctypes "github.com/containers/image/v5/types"
)

// Client provides methods required to be implemented by the various target registry clients, e.g. ECR, Docker, Quay.
type Client interface {
	CreateRepository(ctx context.Context, name string) error
	RepositoryExists() bool
	CopyImage(ctx context.Context, src ctypes.ImageReference, srcCreds string, dest ctypes.ImageReference, destCreds string) error
	PullImage() error
	PutImage() error
	ImageExists(ctx context.Context, ref ctypes.ImageReference) bool

	// Endpoint returns the domain of the registry
	Endpoint() string
	Credentials() string
}
