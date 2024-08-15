package testcontainers

import (
	"context"
	"time"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/store/postgresstore"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Postgres struct {
	DB        *postgresstore.DB
	container testcontainers.Container
}

func NewPostgres(ctx context.Context) (*Postgres, error) {
	container, err := postgres.Run(ctx,
		"postgres:15.3-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	dbURL, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	db := postgresstore.NewDB(dbURL)
	if err := db.Open(); err != nil {
		return nil, err
	}

	return &Postgres{
		DB:        db,
		container: container,
	}, nil
}

func (tc *Postgres) Cleanup(ctx context.Context) error {
	if err := tc.DB.Close(); err != nil {
		return err
	}

	return tc.container.Terminate(ctx)
}
