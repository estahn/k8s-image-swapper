package secrets

import (
	"context"

	"github.com/estahn/k8s-image-swapper/pkg/registry"
	v1 "k8s.io/api/core/v1"
)

// DummyImagePullSecretsProvider does nothing
type DummyImagePullSecretsProvider struct {
}

// NewDummyImagePullSecretsProvider initialises a dummy image pull secrets provider
func NewDummyImagePullSecretsProvider() ImagePullSecretsProvider {
	return &DummyImagePullSecretsProvider{}
}

func (p *DummyImagePullSecretsProvider) SetAuthenticatedRegistries(registries []registry.Client) {
	//empty
}

// GetImagePullSecrets returns an empty ImagePullSecretsResult
func (p *DummyImagePullSecretsProvider) GetImagePullSecrets(ctx context.Context, pod *v1.Pod) (*ImagePullSecretsResult, error) {
	return NewImagePullSecretsResult(), nil
}
