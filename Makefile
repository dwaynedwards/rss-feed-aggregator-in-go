# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./cmd/api
PROJECT_NAME := rss-feed-aggregator
BINARY_NAME := rss-feed

## build: build the users api
.PHONY: build-users
build-users:
	@go build -o=/tmp/${PROJECT_NAME}/bin/${BINARY_NAME}-users-api ${MAIN_PACKAGE_PATH}/users

## run: run the users api
.PHONY: run-users
run-users: build-users
	@/tmp/${PROJECT_NAME}/bin/${BINARY_NAME}-users-api

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	@go fmt ./...
	@go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	@go mod verify
	@go vet ./...
	@go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	@go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	@go test -race -buildvcs -vet=off ./tests/...

## test: run all tests
.PHONY: test
test:
	@go test -v -race -buildvcs ./tests/...

## test-cover: run all tests and display coverage
.PHONY: test-cover
test-cover:
	@go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./tests/...
	@go tool cover -html=/tmp/coverage.out
