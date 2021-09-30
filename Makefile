SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=

.PHONY: help $(MAKECMDGOALS)
.DEFAULT_GOAL := help

export GO111MODULE := on
export GOPROXY = https://proxy.golang.org,direct

help: ## List targets & descriptions
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

setup: ## Install dependencies
	go mod download
	go mod tidy

test: ## Run tests
	LC_ALL=C go test $(TEST_OPTIONS) -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=5m

cover: test ## Run tests and open coverage report
	go tool cover -html=coverage.txt

fmt: ## gofmt and goimports all go files
	gofmt -l -w .
	goimports -l -w .

lint: ## Run linters
	golangci-lint run

e2e: ## Run end-to-end tests
	go test -v -run TestHelmDeployment ./test
