package secrets

import v1 "k8s.io/api/core/v1"

// DummyImagePullSecretsProvider does nothing
type DummyImagePullSecretsProvider struct {
}

// NewDummyImagePullSecretsProvider initialises a dummy image pull secrets provider
func NewDummyImagePullSecretsProvider() ImagePullSecretsProvider {
	return &DummyImagePullSecretsProvider{}
}

// GetImagePullSecrets returns an empty ImagePullSecretsResult
func (p *DummyImagePullSecretsProvider) GetImagePullSecrets(pod *v1.Pod) (*ImagePullSecretsResult, error) {
	return NewImagePullSecretsResult(), nil
}
