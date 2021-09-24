package secrets

import v1 "k8s.io/api/core/v1"

type ImagePullSecretsProvider interface {
	GetImagePullSecrets(pod *v1.Pod) (*ImagePullSecretsResult, error)
}
