# Change these variables as necessary.
API_PACKAGE_PATH := ./cmd/api
PROJECT_NAME := rss-feed-aggregator
BINARY_NAME := rss-feed

MIGRATIONS_PATH := ./store/postgres/migrations
DRIVER := postgres
DBSTRING := "postgres://postgres:postgres@localhost:5432/rss_feeds?sslmode=disable"

## build: build the api
.PHONY: build-api
build-api:
	@go build -o=/tmp/${PROJECT_NAME}/bin/${BINARY_NAME}-api ${API_PACKAGE_PATH}

## run: run the api
.PHONY: run-api
run-api: build-api
	@/tmp/${PROJECT_NAME}/bin/${BINARY_NAME}-api

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
	@go run golang.org/x/vuln/cmd/govulncheck@latest -show verbose ./...
	@go test -race -buildvcs -vet=off ./...

## test: run all tests
.PHONY: test
test:
	@go clean -testcache && go test -race -buildvcs ./...

## test: run all tests
.PHONY: test-base
test-base:
	@go clean -testcache && go test ./...

## test-v: run all tests
.PHONY: test-v
test-v:
	@go clean -testcache && go test -v -race -buildvcs  ./...

## test-cover: run all tests and display coverage
.PHONY: test-cover
test-cover:
	@go test -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	@go tool cover -html=/tmp/coverage.out

# docker-up: spins up docker container
.PHONY: docker-up
docker-up:
	@docker compose up -d

# db-status: gets the migration status of the db
.PHONY: db-status
db-status:
	@goose -dir ${MIGRATIONS_PATH} ${DRIVER} ${DBSTRING} status

# db-up: migrates up the db
.PHONY: db-up
db-up:
	@goose -dir ${MIGRATIONS_PATH} ${DRIVER} ${DBSTRING} up

# db-down: migrates down the db
.PHONY: db-down
db-down:
	@goose -dir ${MIGRATIONS_PATH} ${DRIVER} ${DBSTRING} down

# db-reset: resets the db
.PHONY: db-reset
db-reset:
	@goose -dir ${MIGRATIONS_PATH} ${DRIVER} ${DBSTRING} reset

# db-create-sql: createa sql migration file
.PHONY: db-create-sql
db-create-sql:
	@goose -dir ${MIGRATIONS_PATH} create $(file) sql
