package users_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/alexedwards/argon2id"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/users"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

type testServer struct {
	http.Handler
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

type dummyUserStore struct{}

func (d *dummyUserStore) InsertUser(a *users.User) bool        { return true }
func (d *dummyUserStore) GetUserByID(id uuid.UUID) *users.User { return nil }
func (d *dummyUserStore) GetUserByEmail(e string) *users.User  { return nil }

var dummyStore = &dummyUserStore{}

func newTestServer(userServer users.UsersServer) *testServer {
	s := new(testServer)

	router := http.NewServeMux()
	userServer.RegisterEndpoints(router)
	s.Handler = router

	return s
}

func makeUser(email string, password string, name string) (*users.User, error) {
	id, err := uuid.NewV6()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	return &users.User{
		ID: id, Email: email, Password: hashedPassword, Name: name,
	}, nil
}

func makeSignUpUserRequest(email string, password string, name string) users.SignUpUserRequest {
	return users.SignUpUserRequest{Email: email, Password: password, Name: name}
}

func makeSignInUserRequest(email string, password string) users.SignInUserRequest {
	return users.SignInUserRequest{Email: email, Password: password}
}

func jsonBodyReaderFromStruct(t testing.TB, data any) io.Reader {
	bodyBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	return bytes.NewBuffer(bodyBytes)
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
