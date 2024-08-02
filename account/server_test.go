package account_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/account"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/common"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internals/server"
	"github.com/google/go-cmp/cmp"
)

func TestServer_HealthCheck(t *testing.T) {
	accountStore := account.NewStore()
	accountService := account.NewService(accountStore)
	accountServer := account.NewServer(accountService)
	svr := server.NewServer(accountServer)

	t.Run("GET /healthz returns 200", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/v1/accounts/healthz", nil)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		want := common.HealthCheckResponseMsg
		got := response.Body.String()

		assertResponseStatusCode(t, response.Code, http.StatusOK)
		assertResponseText(t, got, want)
	})
}

func TestServer_Account(t *testing.T) {
	accountStore := account.NewStore()
	accountService := account.NewService(accountStore)
	accountServer := account.NewServer(accountService)
	svr := server.NewServer(accountServer)

	t.Run("POST /accounts creates an account and returns an account id and 201", func(t *testing.T) {
		body := jsonBodyReaderFromStruct(account.CreateAccountRequest{
			Email:    "gopher@go.com",
			Password: "password1",
			Name:     "Gopher",
		})

		request, _ := http.NewRequest(http.MethodPost, "/api/v1/accounts", body)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		assertResponseStatusCode(t, response.Code, http.StatusCreated)
	})

	type badRequestCase struct {
		desc string
		body any
	}
	type badRequest struct {
		email    string
		password string
		name     string
		bad      string
	}
	badCreateRequestCases := []badRequestCase{
		{desc: "unknown field", body: badRequest{email: "gopher@go.com", password: "password1", name: "Gopher", bad: "request"}},
		{desc: "missing email field", body: account.CreateAccountRequest{Name: "Gopher", Password: "password1"}},
		{desc: "missing name field", body: account.CreateAccountRequest{Email: "gopher@go.com", Password: "password1"}},
		{desc: "missing name password", body: account.CreateAccountRequest{Email: "gopher@go.com", Name: "Gopher"}},
		{desc: "missing email, password and name field", body: struct{}{}},
	}

	for _, test := range badCreateRequestCases {
		t.Run(fmt.Sprintf("POST /accounts tries to create an account with %s and returns 400", test.desc), func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodPost, "/api/v1/accounts", jsonBodyReaderFromStruct(test.body))
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

		request1, _ := http.NewRequest(http.MethodPost, "/api/v1/accounts", body1)
		response1 := httptest.NewRecorder()

		svr.ServeHTTP(response1, request1)

		body2 := jsonBodyReaderFromStruct(account.CreateAccountRequest{
			Email:    "gopher@go.com",
			Password: "password1",
			Name:     "Gopher",
		})

		request2, _ := http.NewRequest(http.MethodPost, "/api/v1/accounts", body2)
		response2 := httptest.NewRecorder()

		svr.ServeHTTP(response2, request2)

		assertResponseStatusCode(t, response2.Code, http.StatusConflict)
	})

	t.Run("POST /accounts/signin returns 200", func(t *testing.T) {
		body := jsonBodyReaderFromStruct(account.SigninAccountRequest{
			Email:    "gopher@go.com",
			Password: "password1",
		})

		request, _ := http.NewRequest(http.MethodPost, "/api/v1/accounts/signin", body)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		assertResponseStatusCode(t, response.Code, http.StatusOK)

		var got *account.SigninAccountResponse

		assertResponseJSON(t, response.Body, &got)

		if got.Token == "" {
			t.Fatal("Expecting there to be a value assigned")
		}
	})

	type incorrectCredentialsCase struct {
		desc string
		body any
	}
	incorrectCredentialsCases := []incorrectCredentialsCase{
		{desc: "incorrect email", body: account.SigninAccountRequest{Email: "incorrectemail@go.com", Password: "password1"}},
		{desc: "incorrect password", body: account.SigninAccountRequest{Email: "gopher@go.com", Password: "incorrectpassword1"}},
		{desc: "incorrect email and password", body: account.SigninAccountRequest{Email: "incorrectemail@go.com", Password: "incorrectpassword1"}},
	}
	for _, test := range incorrectCredentialsCases {
		t.Run(fmt.Sprintf("POST /accounts/signin tries to sign in with %s returns 401", test.desc), func(t *testing.T) {
			body := jsonBodyReaderFromStruct(test.body)

			request, _ := http.NewRequest(http.MethodPost, "/api/v1/accounts/signin", body)
			response := httptest.NewRecorder()

			svr.ServeHTTP(response, request)

			assertResponseStatusCode(t, response.Code, http.StatusUnauthorized)
		})
	}

	badSigninRequestCases := []badRequestCase{
		{desc: "unknown field", body: badRequest{email: "gopher@go.com", password: "password1", bad: "request"}},
		{desc: "missing email field", body: account.SigninAccountRequest{Password: "password1"}},
		{desc: "missing name password", body: account.SigninAccountRequest{Email: "gopher@go.com"}},
		{desc: "missing email and password field", body: struct{}{}},
	}
	for _, test := range badSigninRequestCases {
		t.Run(fmt.Sprintf("POST /accounts/signin tries to sign in with %s returns 400", test.desc), func(t *testing.T) {
			body := jsonBodyReaderFromStruct(test.body)

			request, _ := http.NewRequest(http.MethodPost, "/api/v1/accounts/signin", body)
			response := httptest.NewRecorder()

			svr.ServeHTTP(response, request)

			assertResponseStatusCode(t, response.Code, http.StatusBadRequest)
		})
	}
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

func assertResponseJSON(t testing.TB, body *bytes.Buffer, got any) {
	t.Helper()

	if err := json.NewDecoder(body).Decode(got); err != nil {
		t.Fatalf("Unable to parse response from server %q , '%v'", body, err)
	}
}

func jsonBodyReaderFromStruct(data any) io.Reader {
	bodyBytes, _ := json.Marshal(data)
	return bytes.NewBuffer(bodyBytes)
}
