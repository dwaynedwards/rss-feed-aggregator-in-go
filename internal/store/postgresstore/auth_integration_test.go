package postgresstore_test

import (
	"context"
	"testing"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/builder"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/service/authservice"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/store/postgresstore"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/testcontainers"
	"github.com/matryer/is"
)

func TestPostgresDBAuthServiceIntegration(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	ctx := context.Background()

	container, err := testcontainers.NewPostgres(ctx)
	is.NoErr(err)

	migration, err := postgresstore.NewPostgresMigration(container.DB, "migrations")
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

	authStore := postgresstore.NewAuthStore(container.DB)
	authService := authservice.NewAuthService(authStore)

	signUpSuccess := builder.NewSignUpRequestBuilder().
		WithName("Gopher").
		WithEmail("gopher1@go.com").
		WithPassword("gogopher1").
		Build()

	token, err := authService.SignUp(ctx, signUpSuccess)

	is.NoErr(err)           // should sign up
	is.True(len(token) > 0) // should receive token

	signInSuccess := builder.NewSignInRequestBuilder().
		WithEmail("gopher1@go.com").
		WithPassword("gogopher1").
		Build()

	token, err = authService.SignIn(ctx, signInSuccess)

	is.NoErr(err)           // should sign in
	is.True(len(token) > 0) // should receive token

	signUpFailure := builder.NewSignUpRequestBuilder().
		WithName("Gopher").
		WithEmail("gopher1@go.com").
		WithPassword("gogopher1").
		Build()

	token, err = authService.SignUp(ctx, signUpFailure)

	is.True(err != nil)      // should fail to sign up with duplicate email
	is.True(len(token) == 0) // should receive no token

	signInEmailFailure := builder.NewSignInRequestBuilder().
		WithEmail("gopher2@go.com").
		WithPassword("gogopher1").
		Build()

	token, err = authService.SignIn(ctx, signInEmailFailure)

	is.True(err != nil)      // should fail to sign in with incorrect email
	is.True(len(token) == 0) // should receive no token

	signInPasswordFailure := builder.NewSignInRequestBuilder().
		WithEmail("gopher1@go.com").
		WithPassword("gogophe2").
		Build()

	token, err = authService.SignIn(ctx, signInPasswordFailure)

	is.True(err != nil)      // should fail to sign in with incorrect email
	is.True(len(token) == 0) // should receive no token
}
