package registry

import (
	"context"

	ctypes "github.com/containers/image/v5/types"
	"github.com/dgraph-io/ristretto"
	"github.com/rs/zerolog/log"
)

// Cached registry cache requests
type Cached struct {
	Cache    *ristretto.Cache
	Registry Client
}

func NewCachedClient(cache *ristretto.Cache, registry Client) (*Cached, error) {
	return &Cached{
		Registry: registry,
		Cache:    cache,
	}, nil
}

func (c *Cached) CreateRepository(ctx context.Context, name string) error {
	if _, found := c.Cache.Get(name); found {
		log.Ctx(ctx).Trace().Str("name", name).Str("method", "CreateRepository").Msg("found in cache")
		return nil
	}

	err := c.Registry.CreateRepository(ctx, name)

	if err == nil {
		c.Cache.Set(name, "", 1)
	}

	return err
}

func (c *Cached) ImageExists(ctx context.Context, imageRef ctypes.ImageReference) bool {
	ref := imageRef.DockerReference().String()
	if _, found := c.Cache.Get(ref); found {
		log.Ctx(ctx).Trace().Str("ref", ref).Str("method", "ImageExists").Msg("found in cache")
		return true
	}

	exists := c.Registry.ImageExists(ctx, imageRef)

	if exists {
		c.Cache.Set(ref, "", 1)
	}
	return exists
}

func (c *Cached) CopyImage(ctx context.Context, src ctypes.ImageReference, srcCreds string, dest ctypes.ImageReference, destCreds string) error {
	return c.Registry.CopyImage(ctx, src, srcCreds, dest, destCreds)
}

func (c *Cached) Endpoint() string {
	return c.Registry.Endpoint()
}

func (c *Cached) Credentials() string {
	return c.Registry.Credentials()
}

func (c *Cached) IsOrigin(imageRef ctypes.ImageReference) bool {
	return c.Registry.IsOrigin(imageRef)
}
