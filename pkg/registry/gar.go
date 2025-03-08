package registry

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	artifactregistry "cloud.google.com/go/artifactregistry/apiv1"
	"github.com/containers/image/v5/docker/reference"
	ctypes "github.com/containers/image/v5/types"
	"github.com/dgraph-io/ristretto"
	"github.com/estahn/k8s-image-swapper/pkg/config"
	"github.com/go-co-op/gocron"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"

	"github.com/rs/zerolog/log"
)

type GARAPI interface{}

type GARClient struct {
	client    GARAPI
	garDomain string
	cache     *ristretto.Cache
	scheduler *gocron.Scheduler
	authToken []byte
}

func NewGARClient(clientConfig config.GCP) (*GARClient, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.StartAsync()

	client := &GARClient{
		client:    nil,
		garDomain: clientConfig.GarDomain(),
		cache:     cache,
		scheduler: scheduler,
	}

	if err := client.scheduleTokenRenewal(); err != nil {
		return nil, err
	}

	return client, nil
}

// CreateRepository is empty since repositories are not created for artifact registry
func (e *GARClient) CreateRepository(ctx context.Context, name string) error {
	return nil
}

func (e *GARClient) CopyImage(ctx context.Context, srcRef ctypes.ImageReference, srcCreds string, destRef ctypes.ImageReference, destCreds string) error {
	src := srcRef.DockerReference().String()
	dest := destRef.DockerReference().String()

	creds := []string{"--src-authfile", srcCreds}

	// use client credentials for any source GAR repositories
	if strings.HasSuffix(reference.Domain(srcRef.DockerReference()), "-docker.pkg.dev") {
		creds = []string{"--src-creds", e.Credentials()}
	}

	app := "skopeo"
	args := []string{
		"--override-os", "linux",
		"copy",
		"--multi-arch", "all",
		"--retry-times", "3",
		dockerPrefix + src,
		dockerPrefix + dest,
	}

	if len(creds[1]) > 0 {
		args = append(args, creds...)
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

	output, cmdErr := commandExecutor(ctx, app, args...).CombinedOutput()

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

func (e *GARClient) ImageExists(ctx context.Context, imageRef ctypes.ImageReference) bool {
	ref := imageRef.DockerReference().String()
	if _, found := e.cache.Get(ref); found {
		log.Ctx(ctx).Trace().Str("ref", ref).Msg("found in cache")
		return true
	}

	app := "skopeo"
	args := []string{
		"inspect",
		"--retry-times", "3",
		dockerPrefix + ref,
		"--creds", e.Credentials(),
	}

	log.Ctx(ctx).Trace().Str("app", app).Strs("args", args).Msg("executing command to inspect image")
	if err := commandExecutor(ctx, app, args...).Run(); err != nil {
		log.Trace().Str("ref", ref).Msg("not found in target repository")
		return false
	}

	log.Ctx(ctx).Trace().Str("ref", ref).Msg("found in target repository")

	e.cache.SetWithTTL(ref, "", 1, 24*time.Hour+time.Duration(rand.Intn(180))*time.Minute)

	return true
}

func (e *GARClient) Endpoint() string {
	return e.garDomain
}

// IsOrigin returns true if the references origin is from this registry
func (e *GARClient) IsOrigin(imageRef ctypes.ImageReference) bool {
	return strings.HasPrefix(imageRef.DockerReference().String(), e.Endpoint())
}

// requestAuthToken requests and returns an authentication token from GAR with its expiration date
func (e *GARClient) requestAuthToken() ([]byte, time.Time, error) {
	ctx := context.Background()
	creds, err := transport.Creds(ctx, option.WithScopes(artifactregistry.DefaultAuthScopes()...))
	if err != nil {
		log.Err(err).Msg("generating gcp creds")
		return []byte(""), time.Time{}, err
	}
	token, err := creds.TokenSource.Token()
	if err != nil {
		log.Err(err).Msg("generating token")
		return []byte(""), time.Time{}, err
	}

	return []byte(fmt.Sprintf("oauth2accesstoken:%v", token.AccessToken)), token.Expiry, nil
}

// scheduleTokenRenewal sets a scheduler to execute token renewal before the token expires
func (e *GARClient) scheduleTokenRenewal() error {
	token, expiryAt, err := e.requestAuthToken()
	if err != nil {
		return err
	}

	renewalAt := expiryAt.Add(-2 * time.Minute)
	e.authToken = token

	log.Debug().Time("expiryAt", expiryAt).Time("renewalAt", renewalAt).Msg("auth token set, schedule next token renewal")

	j, _ := e.scheduler.Every(1).StartAt(renewalAt).Do(e.scheduleTokenRenewal)
	j.LimitRunsTo(1)

	return nil
}

func (e *GARClient) Credentials() string {
	return string(e.authToken)
}

func (e *GARClient) DockerConfig() ([]byte, error) {
	dockerConfig := DockerConfig{
		AuthConfigs: map[string]AuthConfig{
			e.garDomain: {
				Auth: base64.StdEncoding.EncodeToString(e.authToken),
			},
		},
	}

	dockerConfigJson, err := json.Marshal(dockerConfig)
	if err != nil {
		return []byte{}, err
	}

	return dockerConfigJson, nil
}

func NewMockGARClient(garClient GARAPI, garDomain string) (*GARClient, error) {
	client := &GARClient{
		client:    garClient,
		garDomain: garDomain,
		cache:     nil,
		scheduler: nil,
		authToken: []byte("oauth2accesstoken:mock-gar-client-fake-auth-token"),
	}

	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10,      // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 1,       // number of keys per Get buffer.
	})
	client.cache = cache

	return client, nil
}
