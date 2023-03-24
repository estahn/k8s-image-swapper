package webhook

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/containers/image/v5/docker/reference"
	ctypes "github.com/containers/image/v5/types"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"

	"github.com/estahn/k8s-image-swapper/pkg/metrics"
)

// struct representing a job of copying an image with its subcontext
type ImageCopier struct {
	sourcePod      *corev1.Pod
	sourceImageRef ctypes.ImageReference
	targetImageRef ctypes.ImageReference

	imagePullPolicy corev1.PullPolicy
	imageSwapper    *ImageSwapper

	context       context.Context
	cancelContext context.CancelFunc
}

type Task struct {
	function    func() error
	description string
}

var ErrImageAlreadyPresent = errors.New("image already present in target registry")

// replace the default context with a new one with a timeout
func (ic *ImageCopier) withDeadline() *ImageCopier {
	imageCopierContext, imageCopierContextCancel := context.WithTimeout(ic.context, ic.imageSwapper.imageCopyDeadline)
	ic.context = imageCopierContext
	ic.cancelContext = imageCopierContextCancel
	return ic
}

// start the image copy job
func (ic *ImageCopier) start() {
	if _, hasDeadline := ic.context.Deadline(); hasDeadline {
		defer ic.cancelContext()
	}

	// list of actions to execute in order to copy an image
	tasks := []*Task{
		{
			function:    ic.taskCheckImage,
			description: "checking image presence in target registry",
		},
		{
			function:    ic.taskCreateRepository,
			description: "creating a new repository in target registry",
		},
		{
			function:    ic.taskCopyImage,
			description: "copying image data to target repository",
		},
	}

	for _, task := range tasks {
		err := ic.run(task.function)

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				metrics.IncrementError("TimeoutDuringImageCopy")
				log.Ctx(ic.context).Err(err).Msg("timeout during image copy")
			} else if errors.Is(err, ErrImageAlreadyPresent) {
				log.Ctx(ic.context).Trace().Msgf("image copy aborted: %s", err.Error())
			} else {
				metrics.IncrementError("ImageCopyError")
				log.Ctx(ic.context).Err(err).Msgf("image copy error while %s", task.description)
			}
			break
		}
	}
}

// run a task function and check for timeout
func (ic *ImageCopier) run(taskFunc func() error) error {
	if err := ic.context.Err(); err != nil {
		return err
	}

	return taskFunc()
}

func (ic *ImageCopier) taskCheckImage() error {
	registryClient := ic.imageSwapper.registryClient

	imageAlreadyExists := registryClient.ImageExists(ic.context, ic.targetImageRef) && ic.imagePullPolicy != corev1.PullAlways

	if err := ic.context.Err(); err != nil {
		return err
	} else if imageAlreadyExists {
		return ErrImageAlreadyPresent
	}

	return nil
}

func (ic *ImageCopier) taskCreateRepository() error {
	createRepoName := reference.TrimNamed(ic.sourceImageRef.DockerReference()).String()
	metrics.IncrementReposCreateRequests(ic.sourceImageRef.DockerReference().Name(), reference.TrimNamed(ic.sourceImageRef.DockerReference()).String(), createRepoName)
	log.Ctx(ic.context).Debug().Str("repository", createRepoName).Msg("create repository")

	err := ic.imageSwapper.registryClient.CreateRepository(ic.context, createRepoName)
	if err != nil {
		metrics.IncrementEcrError(ic.sourcePod.Namespace, ic.sourceImageRef.DockerReference().Name(), createRepoName, "CreateRepositoryFail")
	}
	return err
}

func (ic *ImageCopier) taskCopyImage() error {
	ctx := ic.context
	// Retrieve secrets and auth credentials
	imagePullSecrets, err := ic.imageSwapper.imagePullSecretProvider.GetImagePullSecrets(ctx, ic.sourcePod)
	// not possible at the moment
	if err != nil {
		metrics.IncrementEcrError(ic.sourcePod.Namespace, ic.sourceImageRef.DockerReference().Name(), reference.TrimNamed(ic.sourceImageRef.DockerReference()).String(), "GetImagePullSecretsFail")
		return err
	}

	authFile, err := imagePullSecrets.AuthFile()
	if err != nil {
		metrics.IncrementEcrError(ic.sourcePod.Namespace, ic.sourceImageRef.DockerReference().Name(), reference.TrimNamed(ic.sourceImageRef.DockerReference()).String(), "AuthFileGenerateFail")
		log.Ctx(ctx).Err(err).Msg("failed generating authFile")
	}

	defer func() {
		if err := os.RemoveAll(authFile.Name()); err != nil {
			metrics.IncrementEcrError(ic.sourcePod.Namespace, ic.sourceImageRef.DockerReference().Name(), reference.TrimNamed(ic.sourceImageRef.DockerReference()).String(), "AuthFileRemoveFail")
			log.Ctx(ctx).Err(err).Str("file", authFile.Name()).Msg("failed removing auth file")
		}
	}()

	if err := ctx.Err(); err != nil {
		return err
	}

	// Copy image
	// TODO: refactor to use structure instead of passing file name / string
	//
	//	or transform registryClient creds into auth compatible form, e.g.
	//	{"auths":{"aws_account_id.dkr.ecr.region.amazonaws.com":{"username":"AWS","password":"..."	}}}
	copyStart := time.Now()
	copyErr := ic.imageSwapper.registryClient.CopyImage(ctx, ic.sourceImageRef, authFile.Name(), ic.targetImageRef, ic.imageSwapper.registryClient.Credentials())
	if copyErr != nil {
		metrics.IncrementEcrError(ic.sourcePod.Namespace, ic.sourceImageRef.DockerReference().Name(), reference.TrimNamed(ic.sourceImageRef.DockerReference()).String(), "CopyImageFail")
		log.Ctx(ctx).Err(err).Msg("copying image to target registry failed")
	} else {
		duration := time.Since(copyStart).Seconds()
		metrics.SetImageCopyDuration(ic.sourcePod.Namespace, ic.sourceImageRef.DockerReference().Name(), reference.TrimNamed(ic.sourceImageRef.DockerReference()).String(), duration)
		log.Ctx(ctx).Debug().Float64("duration", duration).Msg("copied image")
	}
	return copyErr
}