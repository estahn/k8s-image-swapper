package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/alitto/pond"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"github.com/dgraph-io/ristretto"
	"github.com/estahn/k8s-image-swapper/pkg"
	"github.com/estahn/k8s-image-swapper/pkg/registry"
	"github.com/jmespath/go-jmespath"
	"github.com/rs/zerolog/log"
	"github.com/slok/kubewebhook/pkg/webhook"
	whcontext "github.com/slok/kubewebhook/pkg/webhook/context"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/containers/image/v5/docker/reference"
	"github.com/slok/kubewebhook/pkg/webhook/mutating"
)

// ImageSwapper is a mutator that will download images and change the image name.
type ImageSwapper struct {
	registryClient registry.Client

	// filters defines a list of expressions to remove objects that should not be processed,
	// by default all objects will be processed
	filters []pkg.JMESPathFilter

	// downloader manages the download pool
	downloader *pond.WorkerPool

	// cache keeps a list of already downloaded images
	cache *ristretto.Cache
}

// NewImageSwapper returns a new ImageSwapper initialized.
func NewImageSwapper(registryClient registry.Client, filters []pkg.JMESPathFilter) mutating.Mutator {
	return &ImageSwapper{
		registryClient: registryClient,
		filters: filters,
		downloader: pond.New(100, 1000),
	}
}

func NewImageSwapperWebhook(registryClient registry.Client, filters []pkg.JMESPathFilter) (webhook.Webhook, error) {
	imageSwapper := NewImageSwapper(registryClient, filters)
	mt := mutating.MutatorFunc(imageSwapper.Mutate)
	mcfg := mutating.WebhookConfig{
		Name: "k8s-image-swapper",
		Obj:  &corev1.Pod{},
	}

	return mutating.NewWebhook(mcfg, mt, nil, nil, nil)
}

// Mutate will set the required labels on the pods. Satisfies mutating.Mutator interface.
func (p *ImageSwapper) Mutate(ctx context.Context, obj metav1.Object) (bool, error) {
	//switch _ := obj.(type) {
	//case *corev1.Pod:
	//	// o is a pod
	//case *v1beta1.Role:
	//	// o is the actual role Object with all fields etc
	//case *v1beta1.RoleBinding:
	//case *v1beta1.ClusterRole:
	//case *v1beta1.ClusterRoleBinding:
	//case *v1.ServiceAccount:
	//default:
	//	//o is unknown for us
	//}

	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return false, nil
	}

	ar := whcontext.GetAdmissionRequest(ctx)
	logger := log.With().
		Str("uid", string(ar.UID)).
		Str("kind", ar.Kind.String()).
		Str("namespace", ar.Namespace).
		Str("name", pod.Name).
		Logger()
	//spew.Dump()

	lctx := logger.
		WithContext(ctx)

	for i := range pod.Spec.Containers {
		srcRef, err := alltransports.ParseImageName("docker://" + pod.Spec.Containers[i].Image)
		if err != nil {
			log.Ctx(lctx).Warn().Msgf("invalid source name %s: %v", pod.Spec.Containers[i].Image, err)
			continue
		}

		// skip if the source and target registry domain are equal (e.g. same ECR registries)
		if domain := reference.Domain(srcRef.DockerReference()); domain == p.registryClient.Endpoint() {
			continue
		}

		filterCtx := NewFilterContext(*ar, pod)
		if filterMatch(filterCtx, p.filters) {
			log.Ctx(lctx).Info().Msg("skip due to filter condition")
			continue
		}

		targetImage := p.targetName(srcRef)

		log.Ctx(lctx).Debug().Str("image", targetImage).Msg("set new container image")
		pod.Spec.Containers[i].Image = targetImage

		// Create repository
		createRepoName := reference.TrimNamed(srcRef.DockerReference()).String()
		log.Ctx(lctx).Debug().Str("repository", createRepoName).Msg("create repository")
		if err := p.registryClient.CreateRepository(createRepoName); err != nil {
			log.Err(err)
		}

		p.downloader.Submit(func() {
			log.Ctx(lctx).Trace().Str("source", srcRef.DockerReference().String()).Str("target", targetImage).Msg("copy image")
			if err := copyImage(srcRef.DockerReference().String(), "", targetImage, p.registryClient.Credentials()); err != nil {
				log.Ctx(lctx).Err(err).Str("source", srcRef.DockerReference().String()).Str("target", targetImage).Msg("copying image to target registry failed")
			}
		})
	}

	return false, nil
}

// filterMatch returns true if one of the filters matches the context
func filterMatch(ctx FilterContext, filters []pkg.JMESPathFilter) bool {
	// Simplify FilterContext to be easier searchable by marshaling it to JSON and back to an interface
	var filterContext interface{}
	jsonBlob, err := json.Marshal(ctx)
	if err != nil {
		log.Err(err).Msg("could not marshal filter context")
		return false
	}

	err = json.Unmarshal(jsonBlob, &filterContext)
	if err != nil {
		log.Err(err).Msg("could not unmarshal json blob")
		return false
	}

	log.Debug().Interface("object", filterContext).Msg("generated filter context")

	for idx, filter := range filters {
		results, err := jmespath.Search(filter.JMESPath, filterContext)
		log.Debug().Str("filter", filter.JMESPath).Interface("results", results).Msg("jmespath search results")

		if err != nil {
			log.Err(err).Str("filter", filter.JMESPath).Msgf("Filter (idx %v) could not be evaluated.", idx)
			return false
		}

		switch results.(type) {
		case bool:
			if results == true {
				return true
			}
		default:
			log.Warn().Str("filter", filter.JMESPath).Msg("filter does not return a bool value")
		}
	}

	return false
}

// targetName returns the reference in the target repository
func (p *ImageSwapper) targetName(ref types.ImageReference) string {
	return fmt.Sprintf("%s/%s", p.registryClient.Endpoint(), ref.DockerReference().String())
}

// FilterContext is being used by JMESPath to search and match
type FilterContext struct {
	Obj metav1.Object `json:"obj,omitempty"`
}

func NewFilterContext(request v1beta1.AdmissionRequest, obj metav1.Object) FilterContext {
	if obj.GetNamespace() == "" {
		obj.SetNamespace(request.Namespace)
	}

	return FilterContext{Obj: obj}
}

func copyImage(src string, srcCeds string, dest string, destCreds string) error {
	app := "skopeo"
	args := []string{
		"--override-os", "linux",
		"copy",
		"docker://" + src,
		"docker://" + dest,
		"--src-no-creds",
		"--dest-creds", destCreds,
	}

	log.Trace().Str("app", app).Strs("args", args).Msg("executing command to copy image")
	cmd := exec.Command(app, args...)
	_, err := cmd.Output()

	return err
}
