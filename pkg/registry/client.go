package registry

import (
	"context"
	"fmt"

	"github.com/estahn/k8s-image-swapper/pkg/config"
)

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
	DockerConfig() ([]byte, error)
}

// returns a registry client ready for use without the need to specify an implementation
func NewClient(r config.Registry) (Client, error) {
	if err := r.ValidateConfiguration(); err != nil {
		return nil, err
	}

	switch r.Type {
	case config.Aws:
		aws := r.AWS
		ecrDomain := r.GetServerAddress()
		return newECRClient(aws.Region, ecrDomain, aws.AccountID, aws.Role, aws.ECROptions)
	default:
		return nil, fmt.Errorf(`registry of type "%s" is not supported`, r.Type)
	}
}
