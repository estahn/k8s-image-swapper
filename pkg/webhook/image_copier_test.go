package webhook

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/estahn/k8s-image-swapper/pkg/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	corev1 "k8s.io/api/core/v1"
)

func TestImageCopier_withDeadline(t *testing.T) {
	var registryClients []registry.Client
	mutator := NewImageSwapperWithOpts(
		registryClients,
		nil,
		ImageCopyDeadline(8*time.Second),
	)

	imageSwapper, _ := mutator.(*ImageSwapper)

	imageCopier := &ImageCopier{
		imageSwapper: imageSwapper,
		context:      context.Background(),
	}

	imageCopier = imageCopier.withDeadline()
	deadline, hasDeadline := imageCopier.context.Deadline()

	// test that a deadline has been set
	assert.Equal(t, true, hasDeadline)

	// test that the deadline is future
	assert.GreaterOrEqual(t, deadline, time.Now())

	// test that the context can be canceled
	assert.NotEqual(t, nil, imageCopier.context.Done())

	imageCopier.cancelContext()

	_, ok := <-imageCopier.context.Done()
	// test that the Done channel is closed, meaning the context is canceled
	assert.Equal(t, false, ok)

}

func TestImageCopier_tasksTimeout(t *testing.T) {
	ecrClient := new(mockECRClient)
	ecrClient.On(
		"CreateRepositoryWithContext",
		mock.AnythingOfType("*context.timerCtx"),
		&ecr.CreateRepositoryInput{
			ImageScanningConfiguration: &ecr.ImageScanningConfiguration{
				ScanOnPush: aws.Bool(true),
			},
			ImageTagMutability: aws.String("MUTABLE"),
			RepositoryName:     aws.String("docker.io/library/init-container"),
			RegistryId:         aws.String("123456789"),
			Tags: []*ecr.Tag{
				{
					Key:   aws.String("CreatedBy"),
					Value: aws.String("k8s-image-swapper"),
				},
			},
		}).Return(mock.Anything)

	targetRegistryClient, _ := registry.NewMockECRClient(ecrClient, "ap-southeast-2", "123456789.dkr.ecr.ap-southeast-2.amazonaws.com", "123456789", "arn:aws:iam::123456789:role/fakerole")
	srcRegistryClients := []registry.Client{}

	// image swapper with an instant timeout for testing purpose
	mutator := NewImageSwapperWithOpts(
		srcRegistryClients,
		targetRegistryClient,
		ImageCopyDeadline(0*time.Second),
	)

	imageSwapper, _ := mutator.(*ImageSwapper)

	srcRef, _ := alltransports.ParseImageName("docker://library/init-container:latest")
	targetRef, _ := alltransports.ParseImageName("docker://123456789.dkr.ecr.ap-southeast-2.amazonaws.com/docker.io/library/init-container:latest")
	imageCopier := &ImageCopier{
		imageSwapper:    imageSwapper,
		context:         context.Background(),
		sourceImageRef:  srcRef,
		targetImageRef:  targetRef,
		imagePullPolicy: corev1.PullAlways,
		sourcePod: &corev1.Pod{
			Spec: corev1.PodSpec{
				ServiceAccountName: "service-account",
				ImagePullSecrets:   []corev1.LocalObjectReference{},
			},
		},
	}
	imageCopier = imageCopier.withDeadline()

	// test that copy steps generate timeout errors
	var timeoutError error

	timeoutError = imageCopier.run(imageCopier.taskCheckImage)
	assert.Equal(t, context.DeadlineExceeded, timeoutError)

	timeoutError = imageCopier.run(imageCopier.taskCreateRepository)
	assert.Equal(t, context.DeadlineExceeded, timeoutError)

	timeoutError = imageCopier.run(imageCopier.taskCopyImage)
	assert.Equal(t, context.DeadlineExceeded, timeoutError)

	timeoutError = imageCopier.taskCheckImage()
	assert.Equal(t, context.DeadlineExceeded, timeoutError)

	timeoutError = imageCopier.taskCreateRepository()
	assert.Equal(t, context.DeadlineExceeded, timeoutError)

	timeoutError = imageCopier.taskCopyImage()
	assert.Equal(t, context.DeadlineExceeded, timeoutError)
}
