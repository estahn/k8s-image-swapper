package backend

import (
	"context"

	ctypes "github.com/containers/image/v5/types"
)

type Credentials struct {
	// AuthFile is the optional path of the containers authentication file
	AuthFile string
	// Creds optional USERNAME[:PASSWORD] for accessing the registry
	Creds string
}

// Backend describes a image handler
type Backend interface {
	Exists(ctx context.Context, imageRef ctypes.ImageReference, srcCreds Credentials) (bool, error)
	Copy(ctx context.Context, srcRef ctypes.ImageReference, srcCreds Credentials, destRef ctypes.ImageReference, destCreds Credentials) error
}
