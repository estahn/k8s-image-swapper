package webhook

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/containers/image/v5/docker/reference"
	ctypes "github.com/containers/image/v5/types"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
)

// struct representing a job of copying an image with its subcontext
type ImageCopier struct {
	sourcePod      *corev1.Pod
	sourceImageRef ctypes.ImageReference
	targetImage    string

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
				log.Ctx(ic.context).Err(err).Msg("timeout during image copy")
			} else if errors.Is(err, ErrImageAlreadyPresent) {
				log.Ctx(ic.context).Trace().Msgf("image copy aborted: %s", err.Error())
			} else {
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

	imageAlreadyExists := registryClient.ImageExists(ic.context, ic.targetImage) && ic.imagePullPolicy != corev1.PullAlways

	if err := ic.context.Err(); err != nil {
		return err
	} else if imageAlreadyExists {
		return ErrImageAlreadyPresent
	}

	return nil
}

func (ic *ImageCopier) taskCreateRepository() error {
	createRepoName := reference.TrimNamed(ic.sourceImageRef.DockerReference()).String()

	log.Ctx(ic.context).Debug().Str("repository", createRepoName).Msg("create repository")

	return ic.imageSwapper.registryClient.CreateRepository(ic.context, createRepoName)
}

func (ic *ImageCopier) taskCopyImage() error {
	ctx := ic.context
	sourceImage := ic.sourceImageRef.DockerReference().String()
	targetImage := ic.targetImage

	// Retrieve secrets and auth credentials
	imagePullSecrets, err := ic.imageSwapper.imagePullSecretProvider.GetImagePullSecrets(ctx, ic.sourcePod)
	// not possible at the moment
	if err != nil {
		return err
	}

	authFile, err := imagePullSecrets.AuthFile()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed generating authFile")
	}

	defer func() {
		if err := os.RemoveAll(authFile.Name()); err != nil {
			log.Ctx(ctx).Err(err).Str("file", authFile.Name()).Msg("failed removing auth file")
		}
	}()

	// Copy image
	// TODO: refactor to use structure instead of passing file name / string
	//
	//	or transform registryClient creds into auth compatible form, e.g.
	//	{"auths":{"aws_account_id.dkr.ecr.region.amazonaws.com":{"username":"AWS","password":"..."	}}}
	return skopeoCopyImage(ctx, sourceImage, authFile.Name(), targetImage, ic.imageSwapper.registryClient.Credentials())
}

func skopeoCopyImage(ctx context.Context, src string, srcCeds string, dest string, destCreds string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	app := "skopeo"
	args := []string{
		"--override-os", "linux",
		"copy",
		"--all",
		"--retry-times", "3",
		"docker://" + src,
		"docker://" + dest,
	}

	if len(srcCeds) > 0 {
		args = append(args, "--src-authfile", srcCeds)
	} else {
		args = append(args, "--src-no-creds")
	}

	if len(destCreds) > 0 {
		args = append(args, "--dest-creds", destCreds)
	} else {
		args = append(args, "--dest-no-creds")
	}

	log.Ctx(ctx).
		Trace().
		Str("app", app).
		Strs("args", args).
		Msg("execute command to copy image")

	output, cmdErr := exec.CommandContext(ctx, app, args...).CombinedOutput()

	// check if the command timed out during execution for proper logging
	if err := ctx.Err(); err != nil {
		return err
	}

	// enrich error with output from the command which may contain the actual reason
	if cmdErr != nil {
		return fmt.Errorf("Command error, stderr: %s, stdout: %s", cmdErr.Error(), string(output))
	}

	return nil
}
