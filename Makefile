# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./cmd/api
PROJECT_NAME := rss-feed-aggregator
BINARY_NAME := rss-feed

MIGRATIONS_PATH := ./store/postgres/migrations
DRIVER := postgres
DBSTRING := "postgres://postgres:postgres@localhost:5432/rss_feeds?sslmode=disable"

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
	@go test -race -buildvcs -vet=off ./...

## test: run all tests
.PHONY: test
test:
	@go test -race -buildvcs ./...

## test-cover: run all tests and display coverage
.PHONY: test-cover
test-cover:
	@go test -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	@go tool cover -html=/tmp/coverage.out

# db-status: gets the migration status of the db
.PHONY: db-status
db-status:
	@goose -dir ${MIGRATIONS_PATH} ${DRIVER} ${DBSTRING} status

# db-status: migrates up the db
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

# db-create-go: createa go sql migration file
.PHONY: db-create-go
db-create-go:
	@goose -dir ${MIGRATIONS_PATH} create $(file) go