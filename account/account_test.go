package account_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexedwards/argon2id"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/account"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/account/store"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

type testServer struct {
	http.Handler
}

func NewTestServer(accountServer account.AccountServer) *testServer {
	s := new(testServer)

	router := http.NewServeMux()
	accountServer.RegisterEndpoints(router)
	s.Handler = router

	return s
}

type requestCase struct {
	desc string
	body any
}

type badRequest struct {
	email    string
	password string
	name     string
	bad      string
}

type dummyAccountStore struct{}

func (d *dummyAccountStore) Insert(a *account.Account) bool        { return true }
func (d *dummyAccountStore) GetByID(id uuid.UUID) *account.Account { return nil }
func (d *dummyAccountStore) GetByEmail(e string) *account.Account  { return nil }

var dummyStore = &dummyAccountStore{}

func TestAccount_HealthCheck(t *testing.T) {
	accountStore, err := store.NewMapAccountStore()
	if err != nil {
		t.Fatal(err)
	}

	accountService := account.NewService(accountStore)
	accountServer := account.NewServer(accountService)
	server := NewTestServer(accountServer)

	t.Run("GET /healthz returns 200", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/api/v1/accounts/healthz", nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatusCode(t, response.Code, http.StatusOK)
	})
}

func TestAccount_Create(t *testing.T) {
	t.Run("POST /accounts creates an account and returns an account id and 201", func(t *testing.T) {
		accountStore, err := store.NewMapAccountStore()
		if err != nil {
			t.Fatal(err)
		}

		accountService := account.NewService(accountStore)
		accountServer := account.NewServer(accountService)
		server := NewTestServer(accountServer)

		body := jsonBodyReaderFromStruct(t, account.CreateAccountRequest{
			Email:    "gopher@go.com",
			Password: "password1",
			Name:     "Gopher",
		})

		request, err := http.NewRequest(http.MethodPost, "/api/v1/accounts", body)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatusCode(t, response.Code, http.StatusCreated)

		var got *account.CreateAccountResponse

		assertResponseJSON(t, response.Body, &got)
	})

	badCreateRequestCases := []requestCase{
		{desc: "unknown field", body: badRequest{email: "gopher@go.com", password: "password1", name: "Gopher", bad: "request"}},
		{desc: "missing email field", body: account.CreateAccountRequest{Name: "Gopher", Password: "password1"}},
		{desc: "missing name field", body: account.CreateAccountRequest{Email: "gopher@go.com", Password: "password1"}},
		{desc: "missing name password", body: account.CreateAccountRequest{Email: "gopher@go.com", Name: "Gopher"}},
		{desc: "missing email, password and name field", body: struct{}{}},
	}

	for _, test := range badCreateRequestCases {
		accountService := account.NewService(dummyStore)
		accountServer := account.NewServer(accountService)
		server := NewTestServer(accountServer)

		t.Run(fmt.Sprintf("POST /accounts tries to create an account with %s and returns 422", test.desc), func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, "/api/v1/accounts", jsonBodyReaderFromStruct(t, test.body))
			if err != nil {
				t.Fatal(err)
			}

			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertResponseStatusCode(t, response.Code, http.StatusUnprocessableEntity)
		})
	}

	t.Run("POST /accounts tries to create an account with an email that already exists returns 409", func(t *testing.T) {
		id, err := uuid.NewV6()
		if err != nil {
			t.Fatal(err)
		}

		initialDB := map[uuid.UUID]*account.Account{
			id: {ID: id, Email: "gopher@go.com", Password: "password1", Name: "Gopher"},
		}

		accountStore, err := store.NewMapAccountStore(store.WithDB(initialDB))
		if err != nil {
			t.Fatal(err)
		}

		accountService := account.NewService(accountStore)
		accountServer := account.NewServer(accountService)
		server := NewTestServer(accountServer)

		body := jsonBodyReaderFromStruct(t, account.CreateAccountRequest{
			Email:    "gopher@go.com",
			Password: "password1",
			Name:     "Gopher",
		})

		request, err := http.NewRequest(http.MethodPost, "/api/v1/accounts", body)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatusCode(t, response.Code, http.StatusConflict)
	})
}

func TestAccount_Signin(t *testing.T) {
	t.Run("POST /accounts/signin returns 200", func(t *testing.T) {
		id, err := uuid.NewV6()
		if err != nil {
			t.Fatal(err)
		}

		password := "password1"
		hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
		if err != nil {
			t.Fatal(err)
		}

		initialDB := map[uuid.UUID]*account.Account{
			id: {ID: id, Email: "gopher@go.com", Password: hashedPassword, Name: "Gopher"},
		}

		accountStore, err := store.NewMapAccountStore(store.WithDB(initialDB))
		if err != nil {
			t.Fatal(err)
		}

		accountService := account.NewService(accountStore)
		accountServer := account.NewServer(accountService)
		server := NewTestServer(accountServer)

		body := jsonBodyReaderFromStruct(t, account.SigninAccountRequest{
			Email:    "gopher@go.com",
			Password: password,
		})

		request, err := http.NewRequest(http.MethodPost, "/api/v1/accounts/signin", body)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatusCode(t, response.Code, http.StatusOK)

		var got *account.SigninAccountResponse

		assertResponseJSON(t, response.Body, &got)
	})

	incorrectCredentialsCases := []requestCase{
		{desc: "incorrect email", body: account.SigninAccountRequest{Email: "incorrectemail@go.com", Password: "password1"}},
		{desc: "incorrect password", body: account.SigninAccountRequest{Email: "gopher@go.com", Password: "incorrectpassword1"}},
		{desc: "incorrect email and password", body: account.SigninAccountRequest{Email: "incorrectemail@go.com", Password: "incorrectpassword1"}},
	}

	for _, test := range incorrectCredentialsCases {
		t.Run(fmt.Sprintf("POST /accounts/signin tries to sign in with %s returns 401", test.desc), func(t *testing.T) {
			id, err := uuid.NewV6()
			if err != nil {
				t.Fatal(err)
			}

			password := "password1"
			hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
			if err != nil {
				t.Fatal(err)
			}

			initialDB := map[uuid.UUID]*account.Account{
				id: {ID: id, Email: "gopher@go.com", Password: hashedPassword, Name: "Gopher"},
			}

			accountStore, err := store.NewMapAccountStore(store.WithDB(initialDB))
			if err != nil {
				t.Fatal(err)
			}

			accountService := account.NewService(accountStore)
			accountServer := account.NewServer(accountService)
			server := NewTestServer(accountServer)

			body := jsonBodyReaderFromStruct(t, test.body)

			request, err := http.NewRequest(http.MethodPost, "/api/v1/accounts/signin", body)
			if err != nil {
				t.Fatal(err)
			}

			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertResponseStatusCode(t, response.Code, http.StatusUnauthorized)
		})
	}

	badSigninRequestCases := []requestCase{
		{desc: "unknown field", body: badRequest{email: "gopher@go.com", password: "password1", bad: "request"}},
		{desc: "missing email field", body: account.SigninAccountRequest{Password: "password1"}},
		{desc: "missing name password", body: account.SigninAccountRequest{Email: "gopher@go.com"}},
		{desc: "missing email and password field", body: struct{}{}},
	}

	for _, test := range badSigninRequestCases {
		t.Run(fmt.Sprintf("POST /accounts/signin tries to sign in with %s returns 422", test.desc), func(t *testing.T) {
			accountService := account.NewService(dummyStore)
			accountServer := account.NewServer(accountService)
			server := NewTestServer(accountServer)

			body := jsonBodyReaderFromStruct(t, test.body)

			request, err := http.NewRequest(http.MethodPost, "/api/v1/accounts/signin", body)
			if err != nil {
				t.Fatal(err)
			}

			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertResponseStatusCode(t, response.Code, http.StatusUnprocessableEntity)
		})
	}
}

func assertResponseStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if !cmp.Equal(got, want) {
		t.Errorf("Did not get correct status code, got %d, want %d", got, want)
	}
}

func assertResponseJSON(t testing.TB, body *bytes.Buffer, got any) {
	t.Helper()

	if err := json.NewDecoder(body).Decode(got); err != nil {
		t.Fatalf("Unable to parse response from server %q , '%v'", body, err)
	}
}

func jsonBodyReaderFromStruct(t testing.TB, data any) io.Reader {
	bodyBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	return bytes.NewBuffer(bodyBytes)
}
