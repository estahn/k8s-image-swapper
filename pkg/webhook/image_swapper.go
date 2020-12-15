package webhook

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"github.com/estahn/k8s-image-swapper/pkg"
	"github.com/jmespath/go-jmespath"
	"github.com/rs/zerolog/log"
	"github.com/slok/kubewebhook/pkg/webhook"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/containers/image/v5/docker/reference"
	"github.com/slok/kubewebhook/pkg/webhook/mutating"
)

// ImageSwapper is a mutator that will set labels on the received pods.
type ImageSwapper struct {
	// targetRegistry contains the domain of the target registry, e.g. docker.io
	targetRegistry string

	// filters defines a list of expressions to remove objects that should not be processed,
	// by default all objects will be processed
	filters []pkg.JMESPathFilter
}

// NewImageSwapper returns a new ImageSwapper initialized.
func NewImageSwapper(targetRegistry string, filters []pkg.JMESPathFilter) mutating.Mutator {
	return &ImageSwapper{
		targetRegistry: targetRegistry,
		filters: filters,
	}
}

func NewImageSwapperWebhook(targetRegistry string, filters []pkg.JMESPathFilter) (webhook.Webhook, error) {
	imageSwapper := NewImageSwapper(targetRegistry, filters)
	mt := mutating.MutatorFunc(imageSwapper.Mutate)
	mcfg := mutating.WebhookConfig{
		Name: "k8s-image-swapper",
		Obj:  &corev1.Pod{},
	}

	return mutating.NewWebhook(mcfg, mt, nil, nil, nil)
}

// Mutate will set the required labels on the pods. Satisfies mutating.Mutator interface.
func (p *ImageSwapper) Mutate(ctx context.Context, obj metav1.Object) (bool, error) {
	//switch o := obj.(type) {
	//case *v1.Pod:
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

	// TODO: Refactor to be outside of Mutate to avoid per request token creation
	// TODO: Implement re-issue auth token if operation fails
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		//Config:            aws.Config{Region: aws.String(config.ReplicationConfigs[0].Region)},
	}))
	ecrClient := ecr.New(sess, &aws.Config{Region: aws.String("ap-southeast-2")})

	getAuthTokenOutput, err := ecrClient.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		log.Err(err).Msg("could not fetch auth token for ECR")
		return false, err
	}

	authToken, err := base64.StdEncoding.DecodeString(*getAuthTokenOutput.AuthorizationData[0].AuthorizationToken)
	if err != nil {
		log.Err(err).Msg("could not decode auth token")
		return false, err
	}

	for i := range pod.Spec.Containers {
		srcRef, err := alltransports.ParseImageName("docker://" + pod.Spec.Containers[i].Image)
		if err != nil {
			log.Warn().Msgf("invalid source name %s: %v", pod.Spec.Containers[i].Image, err)
			continue
		}

		// skip if the source and target registry domain are equal (e.g. same ECR registries)
		if domain := reference.Domain(srcRef.DockerReference()); domain == p.targetRegistry {
			continue
		}

		filterCtx := FilterContext{
			Obj: obj,
		}

		if filterMatch(filterCtx, p.filters) {
			continue
		}

		targetImage := p.targetName(srcRef)

		log.Debug().Str("image", targetImage).Msg("set new container image")
		pod.Spec.Containers[i].Image = targetImage

		// Create repository
		createRepoName := reference.TrimNamed(srcRef.DockerReference()).String()
		log.Debug().Str("repository", createRepoName).Msg("create repository")
		_, err = ecrClient.CreateRepository(&ecr.CreateRepositoryInput{
			RepositoryName: aws.String(createRepoName),
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case ecr.ErrCodeServerException:
					log.Err(aerr).Msg(ecr.ErrCodeServerException)
				case ecr.ErrCodeInvalidParameterException:
					log.Err(aerr).Msg(ecr.ErrCodeInvalidParameterException)
				case ecr.ErrCodeInvalidTagParameterException:
					log.Err(aerr).Msg(ecr.ErrCodeInvalidTagParameterException)
				case ecr.ErrCodeTooManyTagsException:
					log.Err(aerr).Msg(ecr.ErrCodeTooManyTagsException)
				case ecr.ErrCodeRepositoryAlreadyExistsException:
					log.Info().Msg(ecr.ErrCodeRepositoryAlreadyExistsException)
				case ecr.ErrCodeLimitExceededException:
					log.Err(aerr).Msg(ecr.ErrCodeLimitExceededException)
				default:
					log.Err(aerr)
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				log.Err(err)
			}
		}

		// TODO: refactor, moving this into a subroutine for now to speed up response.
		go func() {
			log.Info().Str("source", srcRef.DockerReference().String()).Str("target", targetImage).Msgf("copy image")
			app := "skopeo"
			args := []string{
				"--override-os", "linux",
				"copy",
				"docker://" + srcRef.DockerReference().String(),
				"docker://" + targetImage,
				"--src-no-creds",
				"--dest-creds", string(authToken),
			}

			log.Debug().Str("app", app).Strs("args", args).Msg("executing command to copy image")

			cmd := exec.Command(app, args...)
			stdout, err := cmd.Output()

			if err != nil {
				log.Err(err).Bytes("output", stdout).Msg("copying image to target registry failed")
				//continue
			}
		}()
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
	return fmt.Sprintf("%s/%s", p.targetRegistry, ref.DockerReference().String())
}

// FilterContext is being used by JMESPath to search and match
type FilterContext struct {
	Obj metav1.Object `json:"obj,omitempty"`
}
