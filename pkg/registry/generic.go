package registry

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	ctypes "github.com/containers/image/v5/types"
	"github.com/dgraph-io/ristretto"
	"github.com/estahn/k8s-image-swapper/pkg/config"
	"github.com/rs/zerolog/log"
)

type GenericAPI interface{}

type GenericClient struct {
	repository string
	username   string
	password   string
	ignoreCert bool
	cache      *ristretto.Cache
}

func NewGenericClient(clientConfig config.Generic) (*GenericClient, error) {

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		return nil, err
	}

	var genericClient = &GenericClient{
		repository: clientConfig.Repository,
		username:   clientConfig.Username,
		password:   clientConfig.Password,
		ignoreCert: clientConfig.IgnoreCert,
		cache:      cache,
	}

	// Only call login if username and password are provided
	if genericClient.username != "" || genericClient.password != "" {
		err = genericClient.login()
		if err != nil {
			return nil, err
		}
	}

	return genericClient, nil
}

func (g *GenericClient) login() error {

	ctx := context.Background()
	app := "skopeo"
	args := []string{
		"login",
		"-u", g.username,
		"--password", g.password,
		g.repository,
	}

	if g.ignoreCert {
		args = append(args, "--tls-verify=false")
	}

	log.Ctx(ctx).
		Trace().
		Str("app", app).
		Strs("args", args).
		Msg("execute command to login to repository")

	log.Trace().Msgf("GenericClient:login - app args %v", args)

	command := commandExecutor(ctx, app, args...)
	output, cmdErr := command.CombinedOutput()

	// enrich error with output from the command which may contain the actual reason
	if cmdErr != nil {
		log.Trace().Msgf("GenericClient:login - Command error, stderr: %s, stdout: %s", cmdErr.Error(), string(output))
		return fmt.Errorf("Command error, stderr: %s, stdout: %s", cmdErr.Error(), string(output))
	}

	return nil
}

func (g *GenericClient) CopyImage(ctx context.Context, srcRef ctypes.ImageReference, srcCreds string, destRef ctypes.ImageReference, destCreds string) error {
	src := srcRef.DockerReference().String()
	dest := destRef.DockerReference().String()

	app := "skopeo"
	args := []string{
		"--override-os", "linux",
		"copy",
		"--multi-arch", "all",
		"--retry-times", "3",
		dockerPrefix + src,
		dockerPrefix + dest,
	}

	//ignore both certs if destination cert is ignored
	if g.ignoreCert {
		args = append(args, "--src-tls-verify=false")
		args = append(args, "--dest-tls-verify=false")
	}

	if len(srcCreds) > 0 {
		args = append(args, "--src-authfile", srcCreds)
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

	log.Trace().Msgf("GenericClient:CopyImage - app args %v", args)
	output, cmdErr := commandExecutor(ctx, app, args...).CombinedOutput()

	// check if the command timed out during execution for proper logging
	if err := ctx.Err(); err != nil {
		return err
	}

	// enrich error with output from the command which may contain the actual reason
	if cmdErr != nil {
		return fmt.Errorf("Command error, stderr: %s, stdout: %s", cmdErr.Error(), string(output))
	}

	log.Info().Msgf("Image copied to target: %s", dest)
	return nil
}

// CreateRepository is empty since repositories are not created for artifact registry
func (g *GenericClient) CreateRepository(ctx context.Context, name string) error {
	return nil
}

func (g *GenericClient) ImageExists(ctx context.Context, imageRef ctypes.ImageReference) bool {
	ref := imageRef.DockerReference().String()
	if _, found := g.cache.Get(ref); found {
		log.Ctx(ctx).Trace().Str("ref", ref).Msg("found in cache")
		return true
	}

	app := "skopeo"
	args := []string{
		"inspect",
		"--retry-times", "3",
		dockerPrefix + ref,
	}

	creds := g.Credentials()
	if creds == "" {
		args = append(args, "--no-creds")
	} else {
		args = append(args, "--creds", creds)
	}

	if g.ignoreCert {
		args = append(args, "--tls-verify=false")
	}

	log.Ctx(ctx).Trace().Str("app", app).Strs("args", args).Msg("executing command to inspect image")
	if err := commandExecutor(ctx, app, args...).Run(); err != nil {
		log.Trace().Str("ref", ref).Msg("not found in repository")
		return false
	}

	log.Ctx(ctx).Trace().Str("ref", ref).Msg("found in repository")
	g.cache.Set(ref, "", 1)

	return true
}

func (g *GenericClient) Endpoint() string {
	return g.repository
}

// IsOrigin returns true if the references origin is from this registry
func (g *GenericClient) IsOrigin(imageRef ctypes.ImageReference) bool {
	return strings.HasPrefix(imageRef.DockerReference().String(), g.Endpoint())
}

func (g *GenericClient) Credentials() string {
	if g.username == "" && g.password == "" {
		return ""
	}
	return g.username + ":" + g.password
}

func (g *GenericClient) DockerConfig() ([]byte, error) {
	var authConfig AuthConfig

	// Use the Credentials method to determine if credentials are present
	creds := g.Credentials()
	if creds != "" {
		authConfig = AuthConfig{
			Auth: base64.StdEncoding.EncodeToString([]byte(creds)),
		}
	}

	// either we generate an empty config (no auth passed) or we use the provided one (username and password given)
	dockerConfig := DockerConfig{
		AuthConfigs: map[string]AuthConfig{
			g.repository: authConfig,
		},
	}

	dockerConfigJson, err := json.Marshal(dockerConfig)
	if err != nil {
		return nil, err
	}

	return dockerConfigJson, nil
}
