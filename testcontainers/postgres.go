package testcontainers

import (
	"context"
	"time"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	rfpg "github.com/dwaynedwards/rss-feed-aggregator-in-go/store/postgres"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresTestContainer struct {
	DB        *rfpg.DB
	container testcontainers.Container
	migration *rf.Migration
}

func NewPostgresTestContainer(ctx context.Context) (*PostgresTestContainer, error) {
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

	db := rfpg.NewDB(dbURL)
	if err := db.Open(); err != nil {
		return nil, err
	}

	migration, err := rfpg.NewMigration(db, false)
	if err != nil {
		return nil, err
	}

	if err := migration.Up(); err != nil {
		return nil, err
	}

	return &PostgresTestContainer{
		DB:        db,
		container: container,
		migration: migration,
	}, nil
}

func (tc *PostgresTestContainer) Cleanup(ctx context.Context) error {
	if err := tc.migration.Reset(); err != nil {
		return err
	}
	if err := tc.migration.Close(); err != nil {
		return err
	}

	return tc.container.Terminate(ctx)
}
