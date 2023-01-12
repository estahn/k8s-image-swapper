package webhook

import (
	"context"
	"errors"
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
	tasks := []func() error{
		ic.taskCheckImage,
		ic.taskCreateRepository,
		ic.taskCopyImage,
	}

	for _, task := range tasks {
		err := ic.run(task)

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				log.Ctx(ic.context).Err(err).Msg("timeout during image copy")
			}
			if errors.Is(err, ErrImageAlreadyPresent) {
				log.Ctx(ic.context).Trace().Msgf("image copy aborted: %s", err.Error())
			}
			break
		}
	}
}

// run a task and update the copy status
func (ic *ImageCopier) run(task func() error) error {
	if err := ic.context.Err(); err != nil {
		return err
	}

	return task()
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

	log.Ctx(ctx).Trace().Msg("copy image")

	// Copy image
	// TODO: refactor to use structure instead of passing file name / string
	//
	//	or transform registryClient creds into auth compatible form, e.g.
	//	{"auths":{"aws_account_id.dkr.ecr.region.amazonaws.com":{"username":"AWS","password":"..."	}}}
	if err := skopeoCopyImage(ctx, sourceImage, authFile.Name(), targetImage, ic.imageSwapper.registryClient.Credentials()); err != nil {
		return err
	}

	return nil
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

	output, err := exec.CommandContext(ctx, app, args...).CombinedOutput()

	log.Ctx(ctx).
		Trace().
		Str("app", app).
		Strs("args", args).
		Bytes("output", output).
		Msg("executed command to copy image")

	return err
}
