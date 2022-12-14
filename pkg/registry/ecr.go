package registry

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/dgraph-io/ristretto"
	"github.com/estahn/k8s-image-swapper/pkg/config"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

type ECRClient struct {
	client          ecriface.ECRAPI
	ecrDomain       string
	authToken       []byte
	cache           *ristretto.Cache
	scheduler       *gocron.Scheduler
	targetAccount   string
	accessPolicy    string
	lifecyclePolicy string
	tags            []config.Tag
}

type DockerConfig struct {
	AuthConfigs map[string]AuthConfig `json:"auths"`
}

type AuthConfig struct {
	Auth string `json:"auth,omitempty"`
}

func (e *ECRClient) Credentials() string {
	return string(e.authToken)
}

func (e *ECRClient) DockerConfig() ([]byte, error) {
	dockerConfig := DockerConfig{
		AuthConfigs: map[string]AuthConfig{
			e.ecrDomain: {
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

func (e *ECRClient) CreateRepository(ctx context.Context, name string) error {
	if _, found := e.cache.Get(name); found {
		return nil
	}

	_, err := e.client.CreateRepositoryWithContext(ctx, &ecr.CreateRepositoryInput{
		RepositoryName: aws.String(name),
		ImageScanningConfiguration: &ecr.ImageScanningConfiguration{
			ScanOnPush: aws.Bool(true),
		},
		ImageTagMutability: aws.String(ecr.ImageTagMutabilityMutable),
		RegistryId:         &e.targetAccount,
		Tags:               e.buildEcrTags(),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeRepositoryAlreadyExistsException:
				// We ignore this case as it is valid.
			default:
				return err
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return err
		}
	}

	if len(e.accessPolicy) > 0 {
		log.Debug().Str("repo", name).Str("accessPolicy", e.accessPolicy).Msg("setting access policy on repo")
		_, err := e.client.SetRepositoryPolicyWithContext(ctx, &ecr.SetRepositoryPolicyInput{
			PolicyText:     &e.accessPolicy,
			RegistryId:     &e.targetAccount,
			RepositoryName: aws.String(name),
		})

		if err != nil {
			log.Err(err).Msg(err.Error())
			return err
		}
	}

	if len(e.lifecyclePolicy) > 0 {
		log.Debug().Str("repo", name).Str("lifecyclePolicy", e.lifecyclePolicy).Msg("setting lifecycle policy on repo")
		_, err := e.client.PutLifecyclePolicyWithContext(ctx, &ecr.PutLifecyclePolicyInput{
			LifecyclePolicyText: &e.lifecyclePolicy,
			RegistryId:          &e.targetAccount,
			RepositoryName:      aws.String(name),
		})

		if err != nil {
			log.Err(err).Msg(err.Error())
			return err
		}
	}

	e.cache.Set(name, "", 1)

	return nil
}

func (e *ECRClient) buildEcrTags() []*ecr.Tag {
	ecrTags := []*ecr.Tag{}

	for _, t := range e.tags {
		tag := ecr.Tag{Key: aws.String(t.Key), Value: aws.String(t.Value)}
		ecrTags = append(ecrTags, &tag)
	}

	return ecrTags
}

func (e *ECRClient) RepositoryExists() bool {
	panic("implement me")
}

func (e *ECRClient) CopyImage() error {
	panic("implement me")
}

func (e *ECRClient) PullImage() error {
	panic("implement me")
}

func (e *ECRClient) PutImage() error {
	panic("implement me")
}

func (e *ECRClient) ImageExists(ctx context.Context, ref string) bool {
	if _, found := e.cache.Get(ref); found {
		return true
	}

	app := "skopeo"
	args := []string{
		"inspect",
		"--retry-times", "3",
		"docker://" + ref,
		"--creds", e.Credentials(),
	}

	log.Ctx(ctx).Trace().Str("app", app).Strs("args", args).Msg("executing command to inspect image")
	if err := exec.CommandContext(ctx, app, args...).Run(); err != nil {
		return false
	}

	e.cache.Set(ref, "", 1)

	return true
}

func (e *ECRClient) Endpoint() string {
	return e.ecrDomain
}

// requestAuthToken requests and returns an authentication token from ECR with its expiration date
func (e *ECRClient) requestAuthToken() ([]byte, time.Time, error) {
	getAuthTokenOutput, err := e.client.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{
		RegistryIds: []*string{&e.targetAccount},
	})

	if err != nil {
		return []byte(""), time.Time{}, err
	}

	authToken, err := base64.StdEncoding.DecodeString(*getAuthTokenOutput.AuthorizationData[0].AuthorizationToken)
	if err != nil {
		return []byte(""), time.Time{}, err
	}

	return authToken, *getAuthTokenOutput.AuthorizationData[0].ExpiresAt, nil
}

// scheduleTokenRenewal sets a scheduler to execute token renewal before the token expires
func (e *ECRClient) scheduleTokenRenewal() error {
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

func newECRClient(region string, ecrDomain string, targetAccount string, role string, options config.ECROptions) (*ECRClient, error) {
	var sess *session.Session
	var config *aws.Config
	if role != "" {
		log.Info().Str("assumedRole", role).Msg("assuming specified role")
		stsSession, _ := session.NewSession(config)
		creds := stscreds.NewCredentials(stsSession, role)
		config = aws.NewConfig().
			WithRegion(region).
			WithCredentialsChainVerboseErrors(true).
			WithHTTPClient(&http.Client{
				Timeout: 3 * time.Second,
			}).
			WithCredentials(creds)
	} else {
		config = aws.NewConfig().
			WithRegion(region).
			WithCredentialsChainVerboseErrors(true).
			WithHTTPClient(&http.Client{
				Timeout: 3 * time.Second,
			})
	}

	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            (*config),
	}))
	ecrClient := ecr.New(sess, config)

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

	client := &ECRClient{
		client:          ecrClient,
		ecrDomain:       ecrDomain,
		cache:           cache,
		scheduler:       scheduler,
		targetAccount:   targetAccount,
		accessPolicy:    options.AccessPolicy,
		lifecyclePolicy: options.LifecyclePolicy,
		tags:            options.Tags,
	}

	if err := client.scheduleTokenRenewal(); err != nil {
		return nil, err
	}

	return client, nil
}

// For testing purposes
func NewDummyECRClient(region string, targetAccount string, role string, options config.ECROptions, authToken []byte) *ECRClient {
	return &ECRClient{
		targetAccount:   targetAccount,
		accessPolicy:    options.AccessPolicy,
		lifecyclePolicy: options.LifecyclePolicy,
		ecrDomain:       fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", targetAccount, region),
		authToken:       authToken,
	}
}

func NewMockECRClient(ecrClient ecriface.ECRAPI, region string, ecrDomain string, targetAccount, role string) (*ECRClient, error) {
	client := &ECRClient{
		client:        ecrClient,
		ecrDomain:     ecrDomain,
		cache:         nil,
		scheduler:     nil,
		targetAccount: targetAccount,
		authToken:     []byte("mock-ecr-client-fake-auth-token"),
		tags:          []config.Tag{{Key: "CreatedBy", Value: "k8s-image-swapper"}, {Key: "AnotherTag", Value: "another-tag"}},
	}

	return client, nil
}
