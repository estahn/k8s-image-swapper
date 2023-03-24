package backend

import (
	"context"
	"fmt"
	"os/exec"

	ctypes "github.com/containers/image/v5/types"
	"github.com/rs/zerolog/log"
)

// Skopeo is the legacy Backend by leveraging execution of the skopeo binary
type Skopeo struct {
	retries int
}

func NewSkopeo() *Skopeo {
	return &Skopeo{
		retries: 3,
	}
}

func (s *Skopeo) credArgs(creds Credentials, prefix string) []string {
	args := make([]string, 0)

	if creds.AuthFile != "" {
		args = append(args, fmt.Sprintf("--%sauthfile", prefix), creds.AuthFile)
	}

	if creds.Creds != "" {
		args = append(args, fmt.Sprintf("--%screds", prefix), creds.Creds)
	}

	if len(args) == 0 {
		args = append(args, fmt.Sprintf("--%sno-creds", prefix), creds.Creds)
	}

	return args
}

func (s *Skopeo) Exists(ctx context.Context, imageRef ctypes.ImageReference, srcCreds Credentials) (bool, error) {
	ref := imageRef.DockerReference().String()

	app := "skopeo"
	args := []string{
		"inspect",
		"--retry-times", "3",
		"docker://" + ref,
	}

	args = append(args, s.credArgs(srcCreds, "")...)

	log.Ctx(ctx).Trace().Str("app", app).Strs("args", args).Msg("executing command to inspect image")
	if err := exec.CommandContext(ctx, app, args...).Run(); err != nil {
		log.Ctx(ctx).Trace().Str("ref", ref).Msg("not found in target repository")
		return false, nil
	}

	log.Ctx(ctx).Trace().Str("ref", ref).Msg("found in target repository")
	return true, nil
}

func (s *Skopeo) Copy(ctx context.Context, srcRef ctypes.ImageReference, srcCreds Credentials, destRef ctypes.ImageReference, destCreds Credentials) error {
	src := srcRef.DockerReference().String()
	dest := destRef.DockerReference().String()
	app := "skopeo"
	args := []string{
		"--override-os", "linux",
		"copy",
		"--multi-arch", "all",
		"--retry-times", "3",
		"docker://" + src,
		"docker://" + dest,
	}

	args = append(args, s.credArgs(srcCreds, "src-")...)
	args = append(args, s.credArgs(destCreds, "dst-")...)

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
		return fmt.Errorf("command error, stderr: %q, stdout: %q", cmdErr.Error(), string(output))
	}

	return nil
}
