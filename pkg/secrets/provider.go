package secrets

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

type ImagePullSecretsProvider interface {
	GetImagePullSecrets(ctx context.Context, pod *v1.Pod) (*ImagePullSecretsResult, error)
}
