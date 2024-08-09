package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/builder"
	rfhttp "github.com/dwaynedwards/rss-feed-aggregator-in-go/http"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/service"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/store/postgres"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/testcontainers"
	"github.com/matryer/is"
)

type APIServer struct {
	*rfhttp.APIServer
}

func TestPostgresDBAuthServiceAPIServerIntegration(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	ctx := context.Background()

	container, err := testcontainers.NewPostgresTestContainer(ctx)
	is.NoErr(err)

	t.Cleanup(func() {
		err := container.Cleanup(ctx)
		is.NoErr(err) // failed to terminate pgContainer
	})

	s := makeAPIServer(postgres.NewAuthStore(container.DB))

	req := builder.NewSignUpAuthRequestBuilder().
		WithEmail("gopher@go.com").
		WithPassword("password1").
		WithName("Gopher").
		Build()
	body := jsonBodyReaderFromStruct(is, req)

	request, err := http.NewRequest(http.MethodPost, "/api/v1/auths/signup", body)
	is.NoErr(err) // should be a successful request

	response := httptest.NewRecorder()

	s.ServeHTTP(response, request)

	req2 := builder.NewSignInAuthRequestBuilder().
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
