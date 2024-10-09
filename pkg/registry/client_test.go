package registry

import (
	"context"
	"fmt"
	"testing"

	"github.com/estahn/k8s-image-swapper/pkg/config"
	"github.com/stretchr/testify/assert"
)

var genConfig = config.Generic{
	Repository: "localhost",
	Username:   "user",
	Password:   "password",
	IgnoreCert: true,
}

func TestNewClientSuccess(t *testing.T) {

	genConfig = config.Generic{
		Repository: "localhost",
		Username:   "user",
		Password:   "password",
		IgnoreCert: true,
	}
	r := config.Registry{
		Type:    "generic",
		Generic: genConfig,
	}

	curCommandExecutor := commandExecutor
	defer func() { commandExecutor = curCommandExecutor }()

	commandExecutor = func(ctx context.Context, name string, arg ...string) ShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v\n", t.Name(), name, arg)
		return testCommandExecutor{
			output: []byte("login successful"),
			err:    nil,
		}
	}

	client, err := NewClient(r)
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestNewClientFailureNoType(t *testing.T) {

	genConfig = config.Generic{
		Repository: "localhost",
		Username:   "user",
		Password:   "password",
		IgnoreCert: true,
	}
	r := config.Registry{
		Type:    "",
		Generic: genConfig,
	}

	client, err := NewClient(r)
	assert.NotNil(t, err)
	assert.Nil(t, client)
	assert.Equal(t, "a registry requires a type", err.Error())
}

func TestNewClientFailureInvalidType(t *testing.T) {

	genConfig = config.Generic{
		Repository: "localhost",
		Username:   "user",
		Password:   "password",
		IgnoreCert: true,
	}
	r := config.Registry{
		Type:    "badType",
		Generic: genConfig,
	}

	client, err := NewClient(r)
	assert.NotNil(t, err)
	assert.Nil(t, client)
	assert.Equal(t, "unknown target registry string: 'badType', defaulting to unknown", err.Error())
}
