package secrets

import v1 "k8s.io/api/core/v1"

// DummyImagePullSecretsProvider does nothing
type DummyImagePullSecretsProvider struct {
}

func NewDummyImagePullSecretsProvider() ImagePullSecretsProvider {
	return &DummyImagePullSecretsProvider{}
}

func (p *DummyImagePullSecretsProvider) GetImagePullSecrets(pod *v1.Pod) (*ImagePullSecretsResult, error) {
	return NewImagePullSecretsResult(), nil
}
