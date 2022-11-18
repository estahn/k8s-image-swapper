package registry

import "context"

// Client provides methods required to be implemented by the various target registry clients, e.g. ECR, Docker, Quay.
type Client interface {
	CreateRepository(ctx context.Context, name string) error
	RepositoryExists() bool
	CopyImage() error
	PullImage() error
	PutImage() error
	ImageExists(ctx context.Context, ref string) bool

	// Endpoint returns the domain of the registry
	Endpoint() string
	Credentials() string
	Dockerconfig() ([]byte, error)
}
