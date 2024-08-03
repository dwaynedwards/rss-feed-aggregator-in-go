# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./cmd/api
PROJECT_NAME := rss-feed-aggregator
BINARY_NAME := rss-feed

## build: build the account api
.PHONY: build/account
build/account:
	@go build -o=/tmp/${PROJECT_NAME}/bin/${BINARY_NAME}-account-api ${MAIN_PACKAGE_PATH}/account

## run: run the account api
.PHONY: run/account
run/account: build/account
	@/tmp/${PROJECT_NAME}/bin/${BINARY_NAME}-account-api


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
	@go test -race -buildvcs -vet=off ./...

## test: run all tests
.PHONY: test
test:
	@go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	@go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	@go tool cover -html=/tmp/coverage.out
