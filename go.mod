module k8s-image-swapper

go 1.12

require (
	github.com/aws/aws-sdk-go v1.32.3
	github.com/containers/image/v5 v5.8.1
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v1.4.2-0.20191219165747-a9416c67da9f
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/rs/zerolog v1.20.0
	github.com/slok/kubewebhook v0.10.0
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
)
