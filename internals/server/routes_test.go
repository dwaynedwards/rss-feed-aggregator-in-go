package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internals/account"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internals/server"
	"github.com/google/go-cmp/cmp"
)

func TestServer_HealthCheck(t *testing.T) {
	accountStore := account.NewAccountStore()
	accountService := account.NewAccountService(accountStore)
	svr := server.NewServer(accountService)

	t.Run("GET /healthz returns 200", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		want := server.HealthCheckResponseMsg
		got := response.Body.String()

		assertResponseStatusCode(t, response.Code, http.StatusOK)
		assertResponseText(t, got, want)
	})
}

func TestServer_Account(t *testing.T) {
	accountStore := account.NewAccountStore()
	accountService := account.NewAccountService(accountStore)
	svr := server.NewServer(accountService)

	t.Run("POST /accounts creates an account and returns an account id and 201", func(t *testing.T) {
		body := jsonBodyReaderFromStruct(account.CreateAccountRequest{
			Email:    "gopher@go.com",
			Password: "password1",
			Name:     "Gopher",
		})

		request, _ := http.NewRequest(http.MethodPost, "/accounts", body)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		assertResponseStatusCode(t, response.Code, http.StatusCreated)
	})

	type BadRequestCase struct {
		Desc string
		Body interface{}
		Want int
	}
	type badRequest struct {
		email    string
		password string
		name     string
		bad      string
	}
	badRequestCases := []BadRequestCase{
		{Desc: "unknown field", Body: badRequest{email: "gopher@go.com", password: "password1", name: "Gopher", bad: "request"}},
		{Desc: "missing email field", Body: account.CreateAccountRequest{Name: "Gopher", Password: "password1"}},
		{Desc: "missing name field", Body: account.CreateAccountRequest{Email: "gopher@go.com", Password: "password1"}},
		{Desc: "missing name password", Body: account.CreateAccountRequest{Email: "gopher@go.com", Name: "Gopher"}},
		{Desc: "missing email, password and name field", Body: struct{}{}},
	}

	for _, test := range badRequestCases {
		t.Run(fmt.Sprintf("POST /accounts tries to create an account with %s and returns 400", test.Desc), func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodPost, "/accounts", jsonBodyReaderFromStruct(test.Body))
			response := httptest.NewRecorder()

			svr.ServeHTTP(response, request)

			assertResponseStatusCode(t, response.Code, http.StatusBadRequest)
		})
	}

	t.Run("POST /accounts trys to create an account with an email that already exists returns 409", func(t *testing.T) {
		body1 := jsonBodyReaderFromStruct(account.CreateAccountRequest{
			Email:    "gopher@go.com",
			Password: "password1",
			Name:     "Gopher",
		})

		request1, _ := http.NewRequest(http.MethodPost, "/accounts", body1)
		response1 := httptest.NewRecorder()

		svr.ServeHTTP(response1, request1)

		body2 := jsonBodyReaderFromStruct(account.CreateAccountRequest{
			Email:    "gopher@go.com",
			Password: "password1",
			Name:     "Gopher",
		})

		request2, _ := http.NewRequest(http.MethodPost, "/accounts", body2)
		response2 := httptest.NewRecorder()

		svr.ServeHTTP(response2, request2)

		assertResponseStatusCode(t, response2.Code, http.StatusConflict)
	})
}

func assertResponseStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if !cmp.Equal(got, want) {
		t.Errorf("Did not get correct status code, got %d, want %d", got, want)
	}
}

func assertResponseText(t testing.TB, got, want string) {
	t.Helper()

	if !cmp.Equal(got, want) {
		t.Errorf("Did not get correct status code, got %s, want %s", got, want)
	}
}

func jsonBodyReaderFromStruct(data interface{}) io.Reader {
	bodyBytes, _ := json.Marshal(data)
	return bytes.NewBuffer(bodyBytes)
}
