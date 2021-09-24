package secrets

import (
	"context"
	"io/ioutil"
	"os"

	jsonpatch "github.com/evanphx/json-patch"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// KubernetesImagePullSecretsProvider retrieves the secrets holding docker auth information from Kubernetes and merges
// them if necessary. Supports Pod secrets as well as ServiceAccount secrets.
type KubernetesImagePullSecretsProvider struct {
	kubernetesClient kubernetes.Interface
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

// Add adds a secrets to internal list and rebuilds the aggregate
func (r *ImagePullSecretsResult) Add(name string, data []byte) {
	r.Secrets[name] = data
	r.Aggregate, _ = jsonpatch.MergePatch(r.Aggregate, data)
}

// AuthFile provides the aggregate as a file to be used by a docker client
func (r *ImagePullSecretsResult) AuthFile() (*os.File, error) {
	tmpfile, err := ioutil.TempFile("", "auth")
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
		kubernetesClient: clientset,
	}
}

// GetImagePullSecrets returns all secrets with their respective content
func (p *KubernetesImagePullSecretsProvider) GetImagePullSecrets(pod *v1.Pod) (*ImagePullSecretsResult, error) {
	var secrets = make(map[string][]byte)

	// retrieve secret names from pod ServiceAccount (spec.imagePullSecrets)
	serviceAccount, err := p.kubernetesClient.CoreV1().
		ServiceAccounts(pod.Namespace).
		Get(context.TODO(), pod.Spec.ServiceAccountName, metav1.GetOptions{})
	if err != nil {
		// TODO: Handle error gracefully, dont panic
		return nil, err
	}

	imagePullSecrets := append(pod.Spec.ImagePullSecrets, serviceAccount.ImagePullSecrets...)

	result := NewImagePullSecretsResult()
	for _, imagePullSecret := range imagePullSecrets {
		// fetch a secret only once
		if _, exists := secrets[imagePullSecret.Name]; exists {
			continue
		}

		secret, _ := p.kubernetesClient.CoreV1().Secrets(pod.Namespace).Get(context.TODO(), imagePullSecret.Name, metav1.GetOptions{})

		if secret.Type != v1.SecretTypeDockerConfigJson {
			continue
		}

		result.Add(imagePullSecret.Name, secret.Data[v1.DockerConfigJsonKey])
	}

	return result, nil
}
