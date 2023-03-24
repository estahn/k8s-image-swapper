package secrets

import (
	"context"

	"github.com/estahn/k8s-image-swapper/pkg/registry"
	v1 "k8s.io/api/core/v1"
)

type ImagePullSecretsProvider interface {
	GetImagePullSecrets(ctx context.Context, pod *v1.Pod) (*ImagePullSecretsResult, error)
	SetAuthenticatedRegistries(privateRegistries []registry.Client)
}
