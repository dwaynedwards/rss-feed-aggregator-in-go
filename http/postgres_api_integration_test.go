package http_test

import (
	"context"
	"testing"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/store/postgres"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/testcontainers"
	"github.com/matryer/is"
)

func TestPostgresDBAuthServiceAPIServerIntegration(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	ctx := context.Background()

	container, err := testcontainers.NewPostgresTestContainer(ctx)
	is.NoErr(err)

	migration, err := postgres.NewPostgresMigration(container.DB, "../store/postgres/migrations")
	is.NoErr(err)

	migration.Up()
	is.NoErr(err)

	t.Cleanup(func() {
		err := migration.Reset()
		is.NoErr(err)
		err = migration.Close()
		is.NoErr(err)
		err = container.Cleanup(ctx)
		is.NoErr(err) // failed to terminate pgContainer
	})
	server := makeAPIServer(postgres.NewAuthStore(container.DB))

	APIServerIntegration(t, is, server)
}
