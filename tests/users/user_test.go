package users_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/users"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/users/store"
)

func TestUsersHealthCheck(t *testing.T) {
	log.SetOutput(io.Discard)

	userStore, err := store.NewMapUserStore()
	if err != nil {
		t.Fatal(err)
	}

	userService := users.NewService(userStore)
	userServer := users.NewServer(userService)
	server := newTestServer(userServer)

	t.Run("GET /status returns 200", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/api/v1/users/status", nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatusCode(t, response.Code, http.StatusOK)
	})
}

func TestUsersSignUp(t *testing.T) {
	t.Run("POST /users/signup creates a user and returns 201", func(t *testing.T) {
		usersStore, err := store.NewMapUserStore()
		if err != nil {
			t.Fatal(err)
		}

		userService := users.NewService(usersStore)
		userServer := users.NewServer(userService)
		server := newTestServer(userServer)

		body := jsonBodyReaderFromStruct(t, makeSignUpUserRequest("gopher@go.com", "password1", "Gopher"))

		request, err := http.NewRequest(http.MethodPost, "/api/v1/users/signup", body)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatusCode(t, response.Code, http.StatusCreated)

		var got *users.SignUpUserResponse

		assertResponseJSON(t, response.Body, &got)
	})

	badCreateRequestCases := []requestCase{
		{desc: "unknown field", body: badRequest{email: "gopher@go.com", password: "password1", name: "Gopher", bad: "request"}},
		{desc: "missing email field", body: makeSignUpUserRequest("", "password1", "Gopher")},
		{desc: "missing name field", body: makeSignUpUserRequest("gopher@go.com", "password1", "")},
		{desc: "missing name password", body: makeSignUpUserRequest("gopher@go.com", "", "Gopher")},
		{desc: "missing email, password and name field", body: struct{}{}},
	}

	for _, test := range badCreateRequestCases {
		userService := users.NewService(dummyStore)
		userServer := users.NewServer(userService)
		server := newTestServer(userServer)

		t.Run(fmt.Sprintf("POST /users tries to create a user with %s and returns 422", test.desc), func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, "/api/v1/users/signup", jsonBodyReaderFromStruct(t, test.body))
			if err != nil {
				t.Fatal(err)
			}

			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertResponseStatusCode(t, response.Code, http.StatusUnprocessableEntity)
		})
	}

	t.Run("POST /users tries to create a user with an email that already exists returns 409", func(t *testing.T) {
		password := "password1"
		user, err := makeUser("gopher@go.com", password, "Gopher")
		if err != nil {
			t.Fatal(err)
		}

		initialDB := store.MapUsersDB{
			user.ID: user,
		}

		usersStore, err := store.NewMapUserStore(store.WithDB(initialDB))
		if err != nil {
			t.Fatal(err)
		}

		userService := users.NewService(usersStore)
		userServer := users.NewServer(userService)
		server := newTestServer(userServer)

		body := jsonBodyReaderFromStruct(t, makeSignUpUserRequest("gopher@go.com", "password1", "Gopher"))

		request, err := http.NewRequest(http.MethodPost, "/api/v1/users/signup", body)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatusCode(t, response.Code, http.StatusConflict)
	})
}

func TestUsersSignin(t *testing.T) {
	t.Run("POST /users/signin signs in a user and returns 200", func(t *testing.T) {
		password := "password1"
		user, err := makeUser("gopher@go.com", password, "Gopher")
		if err != nil {
			t.Fatal(err)
		}

		initialDB := store.MapUsersDB{
			user.ID: user,
		}

		usersStore, err := store.NewMapUserStore(store.WithDB(initialDB))
		if err != nil {
			t.Fatal(err)
		}

		userService := users.NewService(usersStore)
		userServer := users.NewServer(userService)
		server := newTestServer(userServer)

		body := jsonBodyReaderFromStruct(t, makeSignInUserRequest("gopher@go.com", password))

		request, err := http.NewRequest(http.MethodPost, "/api/v1/users/signin", body)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatusCode(t, response.Code, http.StatusOK)

		var got *users.SignInUserResponse

		assertResponseJSON(t, response.Body, &got)
	})

	incorrectCredentialsCases := []requestCase{
		{desc: "incorrect email", body: makeSignInUserRequest("incorrectemail@go.com", "password1")},
		{desc: "incorrect password", body: makeSignInUserRequest("gopher@go.com", "incorrectpassword1")},
		{desc: "incorrect email and password", body: makeSignInUserRequest("incorrectemail@go.com", "incorrectpassword1")},
	}

	for _, test := range incorrectCredentialsCases {
		t.Run(fmt.Sprintf("POST /users/signin tries to sign in a user with %s returns 401", test.desc), func(t *testing.T) {
			password := "password1"
			user, err := makeUser("gopher@go.com", password, "Gopher")
			if err != nil {
				t.Fatal(err)
			}

			initialDB := store.MapUsersDB{
				user.ID: user,
			}

			usersStore, err := store.NewMapUserStore(store.WithDB(initialDB))
			if err != nil {
				t.Fatal(err)
			}

			userService := users.NewService(usersStore)
			userServer := users.NewServer(userService)
			server := newTestServer(userServer)

			body := jsonBodyReaderFromStruct(t, test.body)

			request, err := http.NewRequest(http.MethodPost, "/api/v1/users/signin", body)
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
		{desc: "missing email field", body: makeSignInUserRequest("", "password1")},
		{desc: "missing name password", body: makeSignInUserRequest("gopher@go.com", "")},
		{desc: "missing email and password field", body: struct{}{}},
	}

	for _, test := range badSigninRequestCases {
		t.Run(fmt.Sprintf("POST /users/signin tries to sign in with %s returns 422", test.desc), func(t *testing.T) {
			userService := users.NewService(dummyStore)
			userServer := users.NewServer(userService)
			server := newTestServer(userServer)

			body := jsonBodyReaderFromStruct(t, test.body)

			request, err := http.NewRequest(http.MethodPost, "/api/v1/users/signin", body)
			if err != nil {
				t.Fatal(err)
			}

			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertResponseStatusCode(t, response.Code, http.StatusUnprocessableEntity)
		})
	}
}
