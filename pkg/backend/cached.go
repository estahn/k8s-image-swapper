package backend

import (
	"context"

	ctypes "github.com/containers/image/v5/types"
	"github.com/dgraph-io/ristretto"
	"github.com/rs/zerolog/log"
)

// Cached backend adds a cache layer in front of a backend
type Cached struct {
	Cache   *ristretto.Cache
	Backend Backend
}

func NewCached(cache *ristretto.Cache, backend Backend) *Cached {
	return &Cached{
		Backend: backend,
		Cache:   cache,
	}
}

func (c *Cached) Exists(ctx context.Context, imageRef ctypes.ImageReference, creds Credentials) (bool, error) {
	ref := imageRef.DockerReference().String()
	if _, found := c.Cache.Get(ref); found {
		log.Ctx(ctx).Trace().Str("ref", ref).Msg("found in cache")
		return true, nil
	}

	exists, err := c.Backend.Exists(ctx, imageRef, creds)
	if err != nil {
		return false, err
	}

	if exists {
		c.Cache.Set(ref, "", 1)
	}
	return exists, nil
}

func (c *Cached) Copy(ctx context.Context, srcRef ctypes.ImageReference, srcCreds Credentials, destRef ctypes.ImageReference, destCreds Credentials) error {
	return c.Backend.Copy(ctx, srcRef, srcCreds, destRef, destCreds)
}
