package http_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/builder"
	rfhttp "github.com/dwaynedwards/rss-feed-aggregator-in-go/http"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/service"
	"github.com/matryer/is"
)

type APIServer struct {
	*rfhttp.APIServer
}

func APIServerIntegration(t *testing.T, is *is.I, server *APIServer) {
	req := builder.NewSignUpAuthRequestBuilder().
		WithEmail("gopher@go.com").
		WithPassword("password1").
		WithName("Gopher").
		Build()
	body := structToJSONReader(is, req)

	request, err := http.NewRequest(http.MethodPost, "/api/v1/auths/signup", body)
	is.NoErr(err) // should be a successful request

	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)
	req2 := builder.NewSignInAuthRequestBuilder().
		WithEmail("gopher@go.com").
		WithPassword("password1").
		Build()
	body2 := structToJSONReader(is, req2)

	request2, err := http.NewRequest(http.MethodPost, "/api/v1/auths/signin", body2)
	is.NoErr(err) // should be a successful request

	response2 := httptest.NewRecorder()

	server.ServeHTTP(response2, request2)

	is.Equal(response2.Code, http.StatusOK) // should sign in a user account with a 200 response
}

func makeAPIServer(store rf.AuthStore) *APIServer {
	s := &APIServer{
		APIServer: rfhttp.NewPostgresAPIServer(),
	}
	authService := service.NewAuthService(store)

	s.AuthService = authService

	return s
}

func structToJSONReader(is *is.I, data any) io.Reader {
	bodyBytes, err := json.Marshal(data)
	is.NoErr(err)

	return bytes.NewBuffer(bodyBytes)
}
