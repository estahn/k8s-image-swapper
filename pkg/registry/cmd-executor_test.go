package registry

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCommandExecutor struct {
	CombinedOutputFunc func() ([]byte, error)
	output             []byte
	err                error
}

func (tsc testCommandExecutor) CombinedOutput() ([]byte, error) {
	return tsc.output, tsc.err
}

func (tsc testCommandExecutor) Run() error {
	return tsc.err
}

func TestSuccess(t *testing.T) {

	ctx := context.Background()
	app := "app"
	args := "args"

	shellCmd := commandExecutor(ctx, app, args)
	assert.NotNil(t, shellCmd)

	commandExecutor = func(ctx context.Context, name string, arg ...string) ShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v\n", t.Name(), name, arg)
		return testCommandExecutor{
			output: []byte("command not found"),
			err:    errors.New("copy error"),
		}
	}

	newShellCmd := commandExecutor(ctx, app, args)
	assert.NotNil(t, newShellCmd)

	output, cmdErr := newShellCmd.CombinedOutput()
	assert.Equal(t, "command not found", string(output))
	assert.NotNil(t, cmdErr)
	assert.Equal(t, "copy error", cmdErr.Error())
}
