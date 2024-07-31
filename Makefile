# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./cmd/api
PROJECT_NAME := rss-feed-aggregator
BINARY_NAME := rss-feed

## build: build the application
.PHONY: build
build:
	@go build -o=/tmp/${PROJECT_NAME}/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## run: run the  application
.PHONY: run
run: build
	@/tmp/${PROJECT_NAME}/bin/${BINARY_NAME}


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
	@go run golang.org/x/lint/golint@latest ./...
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
