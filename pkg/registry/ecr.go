package registry

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/containers/image/v5/docker/reference"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	ctypes "github.com/containers/image/v5/types"
	"github.com/estahn/k8s-image-swapper/pkg/backend"
	"github.com/estahn/k8s-image-swapper/pkg/config"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

type ECRClient struct {
	client          ecriface.ECRAPI
	ecrDomain       string
	authToken       []byte
	scheduler       *gocron.Scheduler
	targetAccount   string
	accessPolicy    string
	lifecyclePolicy string
	tags            []config.Tag
	backend         backend.Backend
}

func NewECRClient(clientConfig config.AWS, imageBackend backend.Backend) (*ECRClient, error) {
	ecrDomain := clientConfig.EcrDomain()

	var sess *session.Session
	var config *aws.Config
	if clientConfig.Role != "" {
		log.Info().Str("assumedRole", clientConfig.Role).Msg("assuming specified role")
		stsSession, _ := session.NewSession(config)
		creds := stscreds.NewCredentials(stsSession, clientConfig.Role)
		config = aws.NewConfig().
			WithRegion(clientConfig.Region).
			WithCredentialsChainVerboseErrors(true).
			WithHTTPClient(&http.Client{
				Timeout: 3 * time.Second,
			}).
			WithCredentials(creds)
	} else {
		config = aws.NewConfig().
			WithRegion(clientConfig.Region).
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

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.StartAsync()

	client := &ECRClient{
		client:          ecrClient,
		ecrDomain:       ecrDomain,
		scheduler:       scheduler,
		targetAccount:   clientConfig.AccountID,
		accessPolicy:    clientConfig.ECROptions.AccessPolicy,
		lifecyclePolicy: clientConfig.ECROptions.LifecyclePolicy,
		tags:            clientConfig.ECROptions.Tags,
		backend:         imageBackend,
	}

	if err := client.scheduleTokenRenewal(); err != nil {
		return nil, err
	}

	return client, nil
}

func (e *ECRClient) Credentials() string {
	return string(e.authToken)
}

func (e *ECRClient) CreateRepository(ctx context.Context, name string) error {
	log.Ctx(ctx).Debug().Str("repository", name).Msg("create repository")

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
		log.Ctx(ctx).Debug().Str("repo", name).Str("accessPolicy", e.accessPolicy).Msg("setting access policy on repo")
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
		log.Ctx(ctx).Debug().Str("repo", name).Str("lifecyclePolicy", e.lifecyclePolicy).Msg("setting lifecycle policy on repo")
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

func (e *ECRClient) CopyImage(ctx context.Context, srcRef ctypes.ImageReference, srcCreds string, destRef ctypes.ImageReference, destCreds string) error {
	srcCredentials := backend.Credentials{
		AuthFile: srcCreds,
	}
	dstCredentials := backend.Credentials{
		Creds: destCreds,
	}

	return e.backend.Copy(ctx, srcRef, srcCredentials, destRef, dstCredentials)
}

func (e *ECRClient) ImageExists(ctx context.Context, imageRef ctypes.ImageReference) bool {
	creds := backend.Credentials{
		Creds: e.Credentials(),
	}

	exists, err := e.backend.Exists(ctx, imageRef, creds)
	if err != nil {
		log.Error().Err(err).Msg("unable to check existence of image")
		return false
	}

	return exists
}

func (e *ECRClient) Endpoint() string {
	return e.ecrDomain
}

// IsOrigin returns true if the references origin is from this registry
func (e *ECRClient) IsOrigin(imageRef ctypes.ImageReference) bool {
	domain := reference.Domain(imageRef.DockerReference())
	return domain == e.Endpoint()
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

// For testing purposes
func NewDummyECRClient(region string, targetAccount string, role string, options config.ECROptions, authToken []byte) *ECRClient {
	return &ECRClient{
		targetAccount:   targetAccount,
		accessPolicy:    options.AccessPolicy,
		lifecyclePolicy: options.LifecyclePolicy,
		ecrDomain:       fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", targetAccount, region),
		authToken:       authToken,
		backend:         backend.NewSkopeo(),
	}
}

func NewMockECRClient(ecrClient ecriface.ECRAPI, region string, ecrDomain string, targetAccount, role string) (*ECRClient, error) {
	client := &ECRClient{
		client:        ecrClient,
		ecrDomain:     ecrDomain,
		scheduler:     nil,
		targetAccount: targetAccount,
		authToken:     []byte("mock-ecr-client-fake-auth-token"),
		tags:          []config.Tag{{Key: "CreatedBy", Value: "k8s-image-swapper"}, {Key: "AnotherTag", Value: "another-tag"}},
		backend:       backend.NewSkopeo(),
	}

	return client, nil
}
