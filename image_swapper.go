package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/containers/image/v5/docker/reference"
	"github.com/slok/kubewebhook/pkg/webhook/mutating"
)

// ImageSwapper is a mutator that will set labels on the received pods.
type ImageSwapper struct {
	// targetRegistry contains the domain of the target registry, e.g. docker.io
	targetRegistry string
}

// NewImageSwapper returns a new ImageSwapper initialized.
func NewImageSwapper(targetRegistry string) mutating.Mutator {
	return &ImageSwapper{
		targetRegistry: targetRegistry,
	}
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

	for i := range pod.Spec.Containers {
		srcRef, err := alltransports.ParseImageName("docker://" + pod.Spec.Containers[i].Image)
		if err != nil {
			log.Warn().Msgf("invalid source name %s: %v", pod.Spec.Containers[i].Image, err)
			continue
		}

		// continue if the image is already from the targetRegistry
		if p.skip(srcRef) {
			continue
		}

		targetImage := p.targetName(srcRef)

		pod.Spec.Containers[i].Image = p.targetName(srcRef)

		// logger
		log.Info().Msgf("copy from %s to %s", srcRef.DockerReference().String(), targetImage)
		fmt.Printf("change ns %s pod %s container %s from: %v to: %v", pod.GetNamespace(), pod.GetName(), pod.Spec.Containers[i].Name, srcRef.DockerReference().String(), targetImage)

		//skopeo copy docker://docker.io/library/redis:latest docker://***REMOVED***.dkr.ecr.ap-southeast-2.amazonaws.com/foobar/enrico:redislatest --dest-creds AWS:$(aws ecr get-login-password --output text) --src-no-creds

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			//Config:            aws.Config{Region: aws.String(config.ReplicationConfigs[0].Region)},
		}))
		ecrClient := ecr.New(sess)

		input := &ecr.CreateRepositoryInput{
			RepositoryName: aws.String(reference.TrimNamed(srcRef.DockerReference()).String()),
		}

		// Create repository
		_, err = ecrClient.CreateRepository(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case ecr.ErrCodeServerException:
					fmt.Println(ecr.ErrCodeServerException, aerr.Error())
				case ecr.ErrCodeInvalidParameterException:
					fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
				case ecr.ErrCodeInvalidTagParameterException:
					fmt.Println(ecr.ErrCodeInvalidTagParameterException, aerr.Error())
				case ecr.ErrCodeTooManyTagsException:
					fmt.Println(ecr.ErrCodeTooManyTagsException, aerr.Error())
				case ecr.ErrCodeRepositoryAlreadyExistsException:
					fmt.Println(ecr.ErrCodeRepositoryAlreadyExistsException, aerr.Error())
				case ecr.ErrCodeLimitExceededException:
					fmt.Println(ecr.ErrCodeLimitExceededException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
		}

		app := "/usr/local/bin/skopeo"
		args := []string{
			"--override-os", "linux",
			"copy",
			"docker://" + srcRef.DockerReference().String(),
			"docker://" + targetImage,
			"--src-no-creds",
			"--dest-creds", string(p.destCredentials()),
		}

		cmd := exec.Command(app, args...)
		stdout, err := cmd.Output()

		if err != nil {
			continue
			spew.Dump(err)
			fmt.Println(err.Error())
		}

		spew.Dump(args)
		spew.Dump(stdout)

		// ecrClient.copy(originalImage, transformedImage)
		//rm := pkg.NewRepositoryManager()
		//rm.Copy(originalImage, transformedImage)
	}

	return false, nil
}

// skip returns true if the image should not be mutated
// TODO: Blacklist/Whitelist, exclusion pattern etc
func (p *ImageSwapper) skip(ref types.ImageReference) bool {
	domain := reference.Domain(ref.DockerReference())
	return domain == p.targetRegistry
}

//func (p *ImageSwapper) transform(image string) (original string, transformed string) {
//	// get the normalized form of the image
//	originalImage, _ := reference.ParseAnyReference(image)
//	transformedImage := fmt.Sprintf("%s/%s", p.targetRegistry, originalImage)
//
//	return originalImage.String(), transformedImage
//}

// targetName returns the reference in the target repository
func (p *ImageSwapper) targetName(ref types.ImageReference) string {
	return fmt.Sprintf("%s/%s", p.targetRegistry, ref.DockerReference().String())
}

func (p *ImageSwapper) destCredentials() []byte {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		//Config:            aws.Config{Region: aws.String(config.ReplicationConfigs[0].Region)},
	}))
	ecrClient := ecr.New(sess)

	authToken, err := ecrClient.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		panic(err)
	}

	authData, err := base64.StdEncoding.DecodeString(*authToken.AuthorizationData[0].AuthorizationToken)

	return authData
}
