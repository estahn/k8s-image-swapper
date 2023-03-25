package backend

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/containers/common/pkg/retry"
	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/signature"
	ctypes "github.com/containers/image/v5/types"
	"github.com/rs/zerolog/log"
)

type Native struct {
	retryOpts retry.Options
}

func NewNative() *Native {
	return &Native{
		retryOpts: retry.Options{
			MaxRetry: 3,
			Delay:    time.Millisecond * 100,
		},
	}
}

func (n *Native) newContext(creds Credentials) *ctypes.SystemContext {
	// default is no creds
	dockerAuth := &ctypes.DockerAuthConfig{}

	if creds.Creds != "" {
		username, password, _ := strings.Cut(creds.Creds, ":")
		dockerAuth = &ctypes.DockerAuthConfig{
			Username: username,
			Password: password,
		}
	}

	return &ctypes.SystemContext{
		AuthFilePath:     creds.AuthFile,
		DockerAuthConfig: dockerAuth,

		// It actually defaults to the current runtime, so we may not need to override it
		// OSChoice: "linux",
	}
}

func (n *Native) Exists(ctx context.Context, imageRef ctypes.ImageReference, creds Credentials) (bool, error) {
	var rawManifest []byte

	if err := retry.IfNecessary(ctx, func() error {
		srcImage, err := imageRef.NewImageSource(ctx, n.newContext(creds))
		if err != nil {
			log.Debug().Err(err).Msg("failed to read image source")
			// There is no proper error type we can check, so check for existence of specific message :-(
			// it will fail with something like:
			// reading manifest <tag> in <image>: name unknown: The repository with name '<repo>' does not exist in the registry with id '<id>'
			// reading manifest <tag> in <image>: manifest unknown: Requested image not found
			if strings.Contains(strings.ToLower(err.Error()), "name unknown:") {
				return nil
			}
			if strings.Contains(strings.ToLower(err.Error()), "manifest unknown:") {
				return nil
			}
			return err
		}
		defer srcImage.Close()

		rawManifest, _, err = srcImage.GetManifest(ctx, nil)
		return err
	}, &n.retryOpts); err != nil {
		return false, fmt.Errorf("unable to retrieve manifest for image: %w", err)
	}

	exists := len(rawManifest) > 0

	return exists, nil
}

func (n *Native) Copy(ctx context.Context, srcRef ctypes.ImageReference, srcCreds Credentials, destRef ctypes.ImageReference, destCreds Credentials) error {
	policy, err := signature.DefaultPolicy(nil)
	if err != nil {
		return fmt.Errorf("unable to get image copy policy: %w", err)
	}
	policyContext, err := signature.NewPolicyContext(policy)
	if err != nil {
		return fmt.Errorf("unable to get image copy policy context: %w", err)
	}

	defer func() {
		if err := policyContext.Destroy(); err != nil {
			log.Err(err).Msg("failed to destroy policy context")
		}
	}()

	opts := &copy.Options{
		SourceCtx:          n.newContext(srcCreds),
		DestinationCtx:     n.newContext(destCreds),
		ImageListSelection: copy.CopyAllImages, // multi-arch
	}

	return retry.IfNecessary(ctx, func() error {
		log.Debug().
			Str("dst", destRef.StringWithinTransport()).
			Str("src", srcRef.StringWithinTransport()).
			Msg("copy image started")

		_, err := copy.Image(ctx, policyContext, destRef, srcRef, opts)

		log.Debug().
			Err(err).
			Str("dst", destRef.StringWithinTransport()).
			Str("src", srcRef.StringWithinTransport()).
			Msg("copy image finished")

		if err != nil {
			return fmt.Errorf("failed to copy image: %w", err)
		}
		return nil
	}, &n.retryOpts)
}
