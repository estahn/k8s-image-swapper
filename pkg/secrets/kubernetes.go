package secrets

import (
	"context"
	"fmt"
	"os"

	"github.com/estahn/k8s-image-swapper/pkg/registry"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// KubernetesImagePullSecretsProvider retrieves the secrets holding docker auth information from Kubernetes and merges
// them if necessary. Supports Pod secrets as well as ServiceAccount secrets.
type KubernetesImagePullSecretsProvider struct {
	kubernetesClient        kubernetes.Interface
	authenticatedRegistries []registry.Client
}

// ImagePullSecretsResult contains the result of GetImagePullSecrets
type ImagePullSecretsResult struct {
	Secrets   map[string][]byte
	Aggregate []byte
}

// NewImagePullSecretsResult initialises ImagePullSecretsResult
func NewImagePullSecretsResult() *ImagePullSecretsResult {
	return &ImagePullSecretsResult{
		Secrets:   map[string][]byte{},
		Aggregate: []byte("{}"),
	}
}

// Initialiaze an ImagePullSecretsResult and registers image pull secrets from the given registries
func NewImagePullSecretsResultWithDefaults(defaultImagePullSecrets []registry.Client) *ImagePullSecretsResult {
	imagePullSecretsResult := NewImagePullSecretsResult()
	for index, reg := range defaultImagePullSecrets {
		dockerConfig, err := registry.GenerateDockerConfig(reg)
		if err != nil {
			log.Err(err)
		} else {
			imagePullSecretsResult.Add(fmt.Sprintf("source-registry-%d", index), dockerConfig)
		}
	}
	return imagePullSecretsResult
}

// Add a secrets to internal list and rebuilds the aggregate
func (r *ImagePullSecretsResult) Add(name string, data []byte) {
	r.Secrets[name] = data
	r.Aggregate, _ = jsonpatch.MergePatch(r.Aggregate, data)
}

// AuthFile provides the aggregate as a file to be used by a docker client
func (r *ImagePullSecretsResult) AuthFile() (*os.File, error) {
	tmpfile, err := os.CreateTemp("", "auth")
	if err != nil {
		return nil, err
	}

	if _, err := tmpfile.Write(r.Aggregate); err != nil {
		return nil, err
	}
	if err := tmpfile.Close(); err != nil {
		return nil, err
	}

	return tmpfile, nil
}

func NewKubernetesImagePullSecretsProvider(clientset kubernetes.Interface) ImagePullSecretsProvider {
	return &KubernetesImagePullSecretsProvider{
		kubernetesClient:        clientset,
		authenticatedRegistries: []registry.Client{},
	}
}

func (p *KubernetesImagePullSecretsProvider) SetAuthenticatedRegistries(registries []registry.Client) {
	p.authenticatedRegistries = registries
}

// GetImagePullSecrets returns all secrets with their respective content
func (p *KubernetesImagePullSecretsProvider) GetImagePullSecrets(ctx context.Context, pod *v1.Pod) (*ImagePullSecretsResult, error) {
	var secrets = make(map[string][]byte)

	imagePullSecrets := pod.Spec.ImagePullSecrets

	// retrieve secret names from pod ServiceAccount (spec.imagePullSecrets)
	serviceAccount, err := p.kubernetesClient.CoreV1().
		ServiceAccounts(pod.Namespace).
		Get(ctx, pod.Spec.ServiceAccountName, metav1.GetOptions{})
	if err != nil {
		log.Ctx(ctx).Warn().Msg("error fetching referenced service account, continue without service account imagePullSecrets")
	}

	if serviceAccount != nil {
		imagePullSecrets = append(imagePullSecrets, serviceAccount.ImagePullSecrets...)
	}

	result := NewImagePullSecretsResultWithDefaults(p.authenticatedRegistries)
	for _, imagePullSecret := range imagePullSecrets {
		// fetch a secret only once
		if _, exists := secrets[imagePullSecret.Name]; exists {
			continue
		}

		secret, err := p.kubernetesClient.CoreV1().Secrets(pod.Namespace).Get(ctx, imagePullSecret.Name, metav1.GetOptions{})
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("error fetching secret, continue without imagePullSecrets")
		}

		if secret == nil || secret.Type != v1.SecretTypeDockerConfigJson {
			continue
		}

		result.Add(imagePullSecret.Name, secret.Data[v1.DockerConfigJsonKey])
	}

	return result, nil
}
