package postgres_test

import (
	"context"
	"os"
	"testing"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/service"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/store/postgres"
	"github.com/matryer/is"
	"github.com/pressly/goose/v3"
)

func TestPostgresDBAuthServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestPostgresDBAuthServiceIntegration in short mode")
	}

	is := is.New(t)

	db, mig := MustOpenDB(t, is)
	defer MustCloseDB(t, is, db, mig)

	store := postgres.NewAuthStore(db)
	service := service.NewAuthService(store)

	authSignUpSuccess := rf.NewAuthBuilder().
		WithUser(rf.NewUserBuilder().WithName("Gopher")).
		WithBasicAuth(rf.NewBasicAuthBuilder().
			WithEmail("gopher1@go.com").
			WithPassword("gogopher1")).
		Build()

	token, err := service.SignUp(context.Background(), authSignUpSuccess)

	is.NoErr(err)                                       // should sign up
	is.True(len(token) > 0)                             // should receive token
	is.Equal(authSignUpSuccess.UserID, int64(1))        // auth UserID should be 1
	is.Equal(authSignUpSuccess.ID, int64(1))            // auth ID should be 1
	is.True(!authSignUpSuccess.CreatedAt.IsZero())      // auth CreatedAt should be set
	is.True(!authSignUpSuccess.ModifiedAt.IsZero())     // auth ModifiedAt should be set
	is.True(!authSignUpSuccess.LastSignedInAt.IsZero()) // auth LastLoggedInAt should be set

	authSignInSuccess := rf.NewAuthBuilder().
		WithBasicAuth(rf.NewBasicAuthBuilder().
			WithEmail("gopher1@go.com").
			WithPassword("gogopher1")).
		Build()

	token, err = service.SignIn(context.Background(), authSignInSuccess)

	is.NoErr(err)                                // should sign in
	is.True(len(token) > 0)                      // should receive token
	is.Equal(authSignInSuccess.UserID, int64(1)) // auth UserID should be 1

	authSignUpFailure := rf.NewAuthBuilder().
		WithUser(rf.NewUserBuilder().WithName("Gopher")).
		WithBasicAuth(rf.NewBasicAuthBuilder().
			WithEmail("gopher1@go.com").
			WithPassword("gogopher1")).
		Build()

	_, err = service.SignUp(context.Background(), authSignUpFailure)

	is.True(err != nil) // should fail to sign up with duplicate email

	authSignInEmailFailure := rf.NewAuthBuilder().
		WithBasicAuth(rf.NewBasicAuthBuilder().
			WithEmail("gopher2@go.com").
			WithPassword("gogopher1")).
		Build()

	_, err = service.SignIn(context.Background(), authSignInEmailFailure)

	is.True(err != nil) // should fail to sign in with incorrect email

	authSignInPasswordFailure := rf.NewAuthBuilder().
		WithBasicAuth(rf.NewBasicAuthBuilder().
			WithEmail("gopher1@go.com").
			WithPassword("gogopher2")).
		Build()

	_, err = service.SignIn(context.Background(), authSignInPasswordFailure)

	is.True(err != nil) // should fail to sign in with incorrect email
}

func MustCloseDB(tb testing.TB, is *is.I, db *postgres.DB, migration *rf.Migration) {
	tb.Helper()

	is.NoErr(migration.Reset()) // should reset migration
	is.NoErr(migration.Close()) // should close migration postgres test db connection

	db.Close()
}

func MustOpenDB(tb testing.TB, is *is.I) (*postgres.DB, *rf.Migration) {
	tb.Helper()

	is.NoErr(goose.SetDialect("postgres"))

	dbURL := os.Getenv("TEST_DATABASE_URL")
	db := postgres.NewDB(dbURL)
	is.NoErr(db.Open()) // should open postgres test db connection

	migration, err := postgres.NewMigration(db, false)
	is.NoErr(err) // should create migration

	is.NoErr(migration.Up()) // should up migration
	return db, migration
}
