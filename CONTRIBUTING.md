# Contributing

By participating to this project, you agree to abide our 
[code of conduct](/CODE_OF_CONDUCT.md).

## Setup your machine

`k8s-image-swapper` is written in [Go](https://golang.org/).

Prerequisites:

- `make`
- [Go 1.15+](https://golang.org/doc/install)
- [Docker](https://www.docker.com/) (or [Podman](https://podman.io/))
- [kind](https://kind.sigs.k8s.io/)

Clone `k8s-image-swapper` anywhere:

```sh
git clone git@github.com:estahn/k8s-image-swapper.git
```

Install the build and lint dependencies:

```sh
make setup
```

A good way of making sure everything is all right is running the test suite:

```sh
make test
```

## Test your change

You can create a branch for your changes and try to build from the source as you go:

```sh
make build
```

When you are satisfied with the changes, we suggest you run:

```sh
make ci
```

Which runs all the linters and tests.

## Create a commit

Commit messages should be well formatted, and to make that "standardized", we
are using Conventional Commits.

You can follow the documentation on
[their website](https://www.conventionalcommits.org).

## Submit a pull request

Push your branch to your `k8s-image-swapper` fork and open a pull request against the
main branch.
