package postgres_test

import (
	"context"
	"testing"

	rfbuilder "github.com/dwaynedwards/rss-feed-aggregator-in-go/builder"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/service"
	rfpg "github.com/dwaynedwards/rss-feed-aggregator-in-go/store/postgres"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/testcontainers"
	"github.com/matryer/is"
)

func TestPostgresDBAuthServiceIntegration(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	ctx := context.Background()

	container, err := testcontainers.NewPostgresTestContainer(ctx)
	is.NoErr(err)

	t.Cleanup(func() {
		err := container.Cleanup(ctx)
		is.NoErr(err) // failed to terminate pgContainer
	})

	store := rfpg.NewAuthStore(container.DB)
	service := service.NewAuthService(store)

	authSignUpSuccess := rfbuilder.NewAuthBuilder().
		WithUser(rfbuilder.NewUserBuilder().WithName("Gopher")).
		WithBasicAuth(rfbuilder.NewBasicAuthBuilder().
			WithEmail("gopher1@go.com").
			WithPassword("gogopher1")).
		Build()

	token, err := service.SignUp(ctx, authSignUpSuccess)

	is.NoErr(err)                                       // should sign up
	is.True(len(token) > 0)                             // should receive token
	is.Equal(authSignUpSuccess.UserID, int64(1))        // auth UserID should be 1
	is.Equal(authSignUpSuccess.ID, int64(1))            // auth ID should be 1
	is.True(!authSignUpSuccess.CreatedAt.IsZero())      // auth CreatedAt should be set
	is.True(!authSignUpSuccess.ModifiedAt.IsZero())     // auth ModifiedAt should be set
	is.True(!authSignUpSuccess.LastSignedInAt.IsZero()) // auth LastLoggedInAt should be set

	authSignInSuccess := rfbuilder.NewAuthBuilder().
		WithBasicAuth(rfbuilder.NewBasicAuthBuilder().
			WithEmail("gopher1@go.com").
			WithPassword("gogopher1")).
		Build()

	token, err = service.SignIn(ctx, authSignInSuccess)

	is.NoErr(err)                                // should sign in
	is.True(len(token) > 0)                      // should receive token
	is.Equal(authSignInSuccess.UserID, int64(1)) // auth UserID should be 1

	authSignUpFailure := rfbuilder.NewAuthBuilder().
		WithUser(rfbuilder.NewUserBuilder().WithName("Gopher")).
		WithBasicAuth(rfbuilder.NewBasicAuthBuilder().
			WithEmail("gopher1@go.com").
			WithPassword("gogopher1")).
		Build()

	token, err = service.SignUp(ctx, authSignUpFailure)

	is.True(err != nil)      // should fail to sign up with duplicate email
	is.True(len(token) == 0) // should receive no token

	authSignInEmailFailure := rfbuilder.NewAuthBuilder().
		WithBasicAuth(rfbuilder.NewBasicAuthBuilder().
			WithEmail("gopher2@go.com").
			WithPassword("gogopher1")).
		Build()

	token, err = service.SignIn(ctx, authSignInEmailFailure)

	is.True(err != nil)      // should fail to sign in with incorrect email
	is.True(len(token) == 0) // should receive no token

	authSignInPasswordFailure := rfbuilder.NewAuthBuilder().
		WithBasicAuth(rfbuilder.NewBasicAuthBuilder().
			WithEmail("gopher1@go.com").
			WithPassword("gogopher2")).
		Build()

	token, err = service.SignIn(ctx, authSignInPasswordFailure)

	is.True(err != nil)      // should fail to sign in with incorrect email
	is.True(len(token) == 0) // should receive no token
}
