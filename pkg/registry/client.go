package registry

// Client provides methods required to be implemented by the various target registry clients, e.g. ECR, Docker, Quay.
type Client interface {
	CreateRepository(string) error
	RepositoryExists() bool
	CopyImage() error
	PullImage() error
	PutImage() error
	ImageExists(ref string) bool

	// Endpoint returns the domain of the registry
	Endpoint() string
	Credentials() string
}
