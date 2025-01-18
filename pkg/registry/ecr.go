package registry

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"os/exec"
	"time"

	"github.com/containers/image/v5/docker/reference"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	ctypes "github.com/containers/image/v5/types"
	"github.com/dgraph-io/ristretto"
	"github.com/estahn/k8s-image-swapper/pkg/config"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

type ECRClient struct {
	client        ecriface.ECRAPI
	ecrDomain     string
	authToken     []byte
	cache         *ristretto.Cache
	scheduler     *gocron.Scheduler
	targetAccount string
	options       config.ECROptions
}

func NewECRClient(clientConfig config.AWS) (*ECRClient, error) {
	ecrDomain := clientConfig.EcrDomain()

	var sess *session.Session
	var cfg *aws.Config
	if clientConfig.Role != "" {
		log.Info().Str("assumedRole", clientConfig.Role).Msg("assuming specified role")
		stsSession, _ := session.NewSession(cfg)
		creds := stscreds.NewCredentials(stsSession, clientConfig.Role)
		cfg = aws.NewConfig().
			WithRegion(clientConfig.Region).
			WithCredentialsChainVerboseErrors(true).
			WithHTTPClient(&http.Client{
				Timeout: 3 * time.Second,
			}).
			WithCredentials(creds)
	} else {
		cfg = aws.NewConfig().
			WithRegion(clientConfig.Region).
			WithCredentialsChainVerboseErrors(true).
			WithHTTPClient(&http.Client{
				Timeout: 3 * time.Second,
			})
	}

	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            *cfg,
	}))
	ecrClient := ecr.New(sess, cfg)

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
		client:        ecrClient,
		ecrDomain:     ecrDomain,
		cache:         cache,
		scheduler:     scheduler,
		targetAccount: clientConfig.AccountID,
		options:       clientConfig.ECROptions,
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
	if _, found := e.cache.Get(name); found {
		return nil
	}

	log.Ctx(ctx).Debug().Str("repository", name).Msg("create repository")

	encryptionConfiguration := &ecr.EncryptionConfiguration{
		EncryptionType: aws.String(e.options.EncryptionConfiguration.EncryptionType),
	}

	if e.options.EncryptionConfiguration.EncryptionType == "KMS" {
		encryptionConfiguration.KmsKey = aws.String(e.options.EncryptionConfiguration.KmsKey)
	}

	_, err := e.client.CreateRepositoryWithContext(ctx, &ecr.CreateRepositoryInput{
		RepositoryName:          aws.String(name),
		EncryptionConfiguration: encryptionConfiguration,
		ImageScanningConfiguration: &ecr.ImageScanningConfiguration{
			ScanOnPush: aws.Bool(e.options.ImageScanningConfiguration.ImageScanOnPush),
		},
		ImageTagMutability: aws.String(e.options.ImageTagMutability),
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

	if len(e.options.AccessPolicy) > 0 {
		log.Ctx(ctx).Debug().Str("repo", name).Str("accessPolicy", e.options.AccessPolicy).Msg("setting access policy on repo")
		_, err := e.client.SetRepositoryPolicyWithContext(ctx, &ecr.SetRepositoryPolicyInput{
			PolicyText:     &e.options.AccessPolicy,
			RegistryId:     &e.targetAccount,
			RepositoryName: aws.String(name),
		})

		if err != nil {
			log.Err(err).Msg(err.Error())
			return err
		}
	}

	if len(e.options.LifecyclePolicy) > 0 {
		log.Ctx(ctx).Debug().Str("repo", name).Str("lifecyclePolicy", e.options.LifecyclePolicy).Msg("setting lifecycle policy on repo")
		_, err := e.client.PutLifecyclePolicyWithContext(ctx, &ecr.PutLifecyclePolicyInput{
			LifecyclePolicyText: &e.options.LifecyclePolicy,
			RegistryId:          &e.targetAccount,
			RepositoryName:      aws.String(name),
		})

		if err != nil {
			log.Err(err).Msg(err.Error())
			return err
		}
	}

	e.cache.SetWithTTL(name, "", 1, time.Duration(24*time.Hour))

	return nil
}

func (e *ECRClient) buildEcrTags() []*ecr.Tag {
	ecrTags := []*ecr.Tag{}

	for _, t := range e.options.Tags {
		tag := ecr.Tag{Key: aws.String(t.Key), Value: aws.String(t.Value)}
		ecrTags = append(ecrTags, &tag)
	}

	return ecrTags
}

func (e *ECRClient) RepositoryExists() bool {
	panic("implement me")
}

func (e *ECRClient) CopyImage(ctx context.Context, srcRef ctypes.ImageReference, srcCreds string, destRef ctypes.ImageReference, destCreds string, additionalTag string) error {
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

	if len(additionalTag) > 0 {
		args = append(args, "--additional-tag", additionalTag)
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

	output, cmdErr := exec.CommandContext(ctx, app, args...).CombinedOutput()

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

func (e *ECRClient) PullImage() error {
	panic("implement me")
}

func (e *ECRClient) PutImage() error {
	panic("implement me")
}

func (e *ECRClient) ImageExists(ctx context.Context, imageRef ctypes.ImageReference) bool {
	ref := imageRef.DockerReference().String()
	if _, found := e.cache.Get(ref); found {
		log.Ctx(ctx).Trace().Str("ref", ref).Msg("found in cache")
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
		log.Ctx(ctx).Trace().Str("ref", ref).Msg("not found in target repository")
		return false
	}

	log.Ctx(ctx).Trace().Str("ref", ref).Msg("found in target repository")

	e.cache.SetWithTTL(ref, "", 1, 24*time.Hour+time.Duration(rand.Intn(180))*time.Minute)

	return true
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
		targetAccount: targetAccount,
		options:       options,
		ecrDomain:     fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", targetAccount, region),
		authToken:     authToken,
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
		options: config.ECROptions{
			ImageTagMutability:         "MUTABLE",
			ImageScanningConfiguration: config.ImageScanningConfiguration{ImageScanOnPush: true},
			EncryptionConfiguration:    config.EncryptionConfiguration{EncryptionType: "AES256"},
			Tags:                       []config.Tag{{Key: "CreatedBy", Value: "k8s-image-swapper"}, {Key: "AnotherTag", Value: "another-tag"}},
		},
	}

	return client, nil
}
