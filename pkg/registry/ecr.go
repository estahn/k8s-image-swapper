package registry

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/dgraph-io/ristretto"
)

type ECRClient struct {
	client    *ecr.ECR
	ecrDomain string
	authToken []byte
	cache     *ristretto.Cache
}

func (e *ECRClient) Credentials() string {
	return string(e.authToken)
}

func (e *ECRClient) CreateRepository(name string) error {
	if _, found := e.cache.Get(name); found {
		return nil
	}

	_, err := e.client.CreateRepository(&ecr.CreateRepositoryInput{
		RepositoryName: aws.String(name),
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

	e.cache.Set(name, "", 1)

	return nil
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

func (e *ECRClient) ImageExists() bool {
	panic("implement me")
}

func (e *ECRClient) Endpoint() string {
	return e.ecrDomain
}

func NewECRClient(region string, ecrDomain string) (*ECRClient, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	ecrClient := ecr.New(sess, &aws.Config{Region: aws.String(region)})

	getAuthTokenOutput, err := ecrClient.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return nil, err
	}

	authToken, err := base64.StdEncoding.DecodeString(*getAuthTokenOutput.AuthorizationData[0].AuthorizationToken)
	if err != nil {
		return nil, err
	}

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}

	return &ECRClient{
		client:    ecrClient,
		ecrDomain: ecrDomain,
		authToken: authToken,
		cache: cache,
	}, nil
}
