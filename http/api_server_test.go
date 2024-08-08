package http_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	rfhttp "github.com/dwaynedwards/rss-feed-aggregator-in-go/http"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/service"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/store/postgres"
	"github.com/matryer/is"
	"github.com/pressly/goose/v3"
)

type APIServer struct {
	*rfhttp.APIServer
}

func TestPostgresDBAuthServiceAPIServerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestPostgresDBAuthServiceAPIServerIntegration in short mode")
	}

	is := is.New(t)

	db, mig := mustOpenDB(t, is)
	defer mustCloseDB(t, is, db, mig)

	s := makeAPIServer(postgres.NewAuthStore(db))

	req := rf.NewSignUpAuthRequestBuilder().
		WithEmail("gopher@go.com").
		WithPassword("password1").
		WithName("Gopher").
		Build()
	body := jsonBodyReaderFromStruct(is, req)

	request, err := http.NewRequest(http.MethodPost, "/api/v1/auths/signup", body)
	is.NoErr(err) // should be a successful request

	response := httptest.NewRecorder()

	s.ServeHTTP(response, request)

	req2 := rf.NewSignInAuthRequestBuilder().
		WithEmail("gopher@go.com").
		WithPassword("password1").
		Build()
	body2 := jsonBodyReaderFromStruct(is, req2)

	request2, err := http.NewRequest(http.MethodPost, "/api/v1/auths/signin", body2)
	is.NoErr(err) // should be a successful request

	response2 := httptest.NewRecorder()

	s.ServeHTTP(response2, request2)

	is.Equal(response2.Code, http.StatusOK) // should sign in a user account with a 200 response

	var got rf.SignInAuthResponse
	err = json.NewDecoder(response2.Body).Decode(&got)

	is.NoErr(err)               // should have a response
	is.True(len(got.Token) > 0) // should have a token
}

func mustOpenDB(tb testing.TB, is *is.I) (*postgres.DB, *rf.Migration) {
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

func mustCloseDB(tb testing.TB, is *is.I, db *postgres.DB, migration *rf.Migration) {
	tb.Helper()

	is.NoErr(migration.Reset()) // should reset migration
	is.NoErr(migration.Close()) // should close migration postgres test db connection

	db.Close()
}

func makeAPIServer(store rf.AuthStore) *APIServer {
	s := &APIServer{
		APIServer: rfhttp.NewAPIServer(),
	}
	authService := service.NewAuthService(store)

	s.AuthService = authService

	return s
}

func jsonBodyReaderFromStruct(is *is.I, data any) io.Reader {
	bodyBytes, err := json.Marshal(data)
	is.NoErr(err)

	return bytes.NewBuffer(bodyBytes)
}
