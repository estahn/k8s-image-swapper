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

		// It actually defaults to the current runtime, ao we may not need to override it
		// OSChoice: "linux",
	}
}

func (n *Native) Exists(ctx context.Context, imageRef ctypes.ImageReference, creds Credentials) (bool, error) {
	srcImage, err := imageRef.NewImageSource(ctx, n.newContext(creds))
	if err != nil {
		return false, err
	}
	defer srcImage.Close()

	var rawManifest []byte
	if err := retry.IfNecessary(ctx, func() error {
		rawManifest, _, err = srcImage.GetManifest(ctx, nil)
		return err
	}, &n.retryOpts); err != nil {
		// TODO: check if error is only client errors or also not found?
		return false, fmt.Errorf("Error retrieving manifest for image: %w", err)
	}

	exists := len(rawManifest) > 0

	return exists, nil

}

func (n *Native) Copy(ctx context.Context, srcRef ctypes.ImageReference, srcCreds Credentials, destRef ctypes.ImageReference, destCreds Credentials) error {
	policy, err := signature.DefaultPolicy(nil)
	if err != nil {
		return fmt.Errorf("unable to get image copy policy: %q", err)
	}
	policyContext, err := signature.NewPolicyContext(policy)
	if err != nil {
		return fmt.Errorf("unable to get image copy policy context: %q", err)
	}
	defer policyContext.Destroy()

	opts := &copy.Options{
		SourceCtx:          n.newContext(srcCreds),
		DestinationCtx:     n.newContext(destCreds),
		ImageListSelection: copy.CopyAllImages, // multi-arch
	}

	_, err = copy.Image(ctx, policyContext, destRef, srcRef, opts)

	return retry.IfNecessary(ctx, func() error {
		_, err := copy.Image(ctx, policyContext, destRef, srcRef, opts)
		if err != nil {
			return fmt.Errorf("failed to copy image: %q", err)
		}
		return nil
	}, &n.retryOpts)
}
