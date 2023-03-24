package registry

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	artifactregistry "cloud.google.com/go/artifactregistry/apiv1"
	"github.com/containers/image/v5/docker/reference"
	ctypes "github.com/containers/image/v5/types"
	"github.com/estahn/k8s-image-swapper/pkg/backend"
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
	scheduler *gocron.Scheduler
	authToken []byte
	backend   backend.Backend
}

func NewGARClient(clientConfig config.GCP, imageBackend backend.Backend) (*GARClient, error) {

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.StartAsync()

	client := &GARClient{
		client:    nil,
		garDomain: clientConfig.GarDomain(),
		scheduler: scheduler,
		backend:   imageBackend,
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

func (e *GARClient) RepositoryExists() bool {
	panic("implement me")
}

func (e *GARClient) CopyImage(ctx context.Context, srcRef ctypes.ImageReference, srcCreds string, destRef ctypes.ImageReference, destCreds string) error {
	srcCredentials := backend.Credentials{
		AuthFile: srcCreds,
	}
	dstCredentials := backend.Credentials{
		Creds: destCreds,
	}

	// use client credentials for any source GAR repositories
	if strings.HasSuffix(reference.Domain(srcRef.DockerReference()), "-docker.pkg.dev") {
		srcCredentials = backend.Credentials{
			Creds: e.Credentials(),
		}
	}

	return e.backend.Copy(ctx, srcRef, srcCredentials, destRef, dstCredentials)
}

func (e *GARClient) ImageExists(ctx context.Context, imageRef ctypes.ImageReference) bool {
	creds := backend.Credentials{
		Creds: e.Credentials(),
	}

	exists, err := e.backend.Exists(ctx, imageRef, creds)
	if err != nil {
		log.Error().Err(err).Msg("unable to check existence of image")
		return false
	}
<<<<<<< HEAD

	log.Ctx(ctx).Trace().Str("ref", ref).Msg("found in target repository")

	e.cache.SetWithTTL(ref, "", 1, 24*time.Hour+time.Duration(rand.Intn(180))*time.Minute)

	return true
=======
	return exists
>>>>>>> 359ee16 (feat: Add native image handler backend)
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
		scheduler: nil,
		backend:   backend.NewSkopeo(),
		authToken: []byte("oauth2accesstoken:mock-gar-client-fake-auth-token"),
	}

	return client, nil
}
