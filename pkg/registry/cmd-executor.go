package registry

import (
	"context"
	"os/exec"
)

type ShellCommand interface {
	CombinedOutput() ([]byte, error)
	Run() error
}

type execShellCommand struct {
	*exec.Cmd
}

func newCommandExecutor(ctx context.Context, name string, arg ...string) ShellCommand {
	execCmd := exec.CommandContext(ctx, name, arg...)
	return execShellCommand{Cmd: execCmd}
}

var commandExecutor = newCommandExecutor
