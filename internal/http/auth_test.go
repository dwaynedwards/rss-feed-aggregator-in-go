package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/builder"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/errors"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/mock"
	"github.com/matryer/is"
)

func TestAuthAPI_SignUp_Success(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	t.Run("POST /api/v1/auths/signup signs up a user and returns 201", func(t *testing.T) {
		t.Parallel()

		store := &mock.AuthStore{
			CreateAuthAndUserFn: func(ctx context.Context, auth *rf.Auth) error {
				auth.ID = 1
				auth.UserID = 1
				auth.CreatedAt = time.Now()
				auth.ModifiedAt = time.Now()
				auth.LastSignedInAt = time.Now()
				return nil
			},
			FindByEmailFn: func(ctx context.Context, email string) (*rf.Auth, error) {
				return nil, nil
			},
		}
		s := makeAuthAPIServer(store)

		req := builder.NewSignUpRequestBuilder().
			WithEmail("gopher@go.com").
			WithPassword("password1").
			WithName("Gopher").
			Build()
		body := structToJSONReader(is, req)

		request, err := http.NewRequest(http.MethodPost, "/api/v1/auths/signup", body)
		is.NoErr(err) // should be a successful request

		response := httptest.NewRecorder()

		s.ServeHTTP(response, request)

		is.Equal(response.Code, http.StatusCreated) // should sign up user account with a 201 response
	})
}

func TestAuthService_SignUp_Failure(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	signUpFailureCases := []mock.AuthAPIFailureCase{
		{Desc: "with missing email", AuthReq: mock.SignUpAuthAPIWithMissingEmail, StatusCode: http.StatusBadRequest, Err: errors.ErrEmailRequired},
		{Desc: "with missing password", AuthReq: mock.SignUpAuthAPIWithMissingPassword, StatusCode: http.StatusBadRequest, Err: errors.ErrPasswordRequired},
		{Desc: "with missing name", AuthReq: mock.SignUpAuthAPIWithMissingName, StatusCode: http.StatusBadRequest, Err: errors.ErrNameRequired},
	}
	for _, tc := range signUpFailureCases {
		t.Run(fmt.Sprintf("Should fail to sign up %s", tc.Desc), func(t *testing.T) {
			t.Parallel()

			store := &mock.AuthStore{}
			s := makeAuthAPIServer(store)

			body := structToJSONReader(is, tc.AuthReq)

			request, err := http.NewRequest(http.MethodPost, "/api/v1/auths/signup", body)
			is.NoErr(err) // should be a successful request

			response := httptest.NewRecorder()

			s.ServeHTTP(response, request)

			is.Equal(response.Code, tc.StatusCode) // should not sign up user account with a 400 response

			var got errors.Error
			err = json.NewDecoder(response.Body).Decode(&got)

			is.NoErr(err)                           // should have a response
			is.Equal(got.StatusCode, tc.StatusCode) // shoud have error code

			is.Equal(errors.ToStatusCode(got), tc.StatusCode)    // shoud have error code
			is.True(strings.Contains(errors.ToErr(got), tc.Err)) // should have error message
			is.True(!store.CreateInvoked)                        // auth store sign up should not have been invoked
			is.True(!store.FindByEmailInvoked)                   // auth store FindByEmail should not have been invoked
		})
	}
}

func TestAuthAPI_SignIn_Success(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	t.Run("POST /api/v1/auths/signin signs in a user and returns 200", func(t *testing.T) {
		t.Parallel()

		password := "gogopher1"
		hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
		is.NoErr(err) // Should hash password

		store := &mock.AuthStore{
			FindByEmailFn: func(ctx context.Context, email string) (*rf.Auth, error) {
				auth := builder.NewAuthBuilder().
					WithUserID(1).
					WithBasicAuth(builder.NewBasicAuthBuilder().
						WithPassword(hashedPassword)).
					Build()
				return auth, nil
			},
		}

		s := makeAuthAPIServer(store)

		req := builder.NewSignInRequestBuilder().
			WithEmail("gopher@go.com").
			WithPassword(password).
			Build()
		body := structToJSONReader(is, req)

		request, err := http.NewRequest(http.MethodPost, "/api/v1/auths/signin", body)
		is.NoErr(err) // should be a successful request

		response := httptest.NewRecorder()

		s.ServeHTTP(response, request)

		is.Equal(response.Code, http.StatusOK) // should sign in a user account with a 200 response
	})
}

func TestAuthService_SignIn_Failure(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	signUpFailureCases := []mock.AuthAPIFailureCase{
		{Desc: "with missing email", AuthReq: mock.SignInAuthAPIWithMissingEmail, StatusCode: http.StatusBadRequest, Err: errors.ErrEmailRequired},
		{Desc: "with missing password", AuthReq: mock.SignInAuthAPIWithMissingPassword, StatusCode: http.StatusBadRequest, Err: errors.ErrPasswordRequired},
	}
	for _, tc := range signUpFailureCases {
		t.Run(fmt.Sprintf("Should fail to sign up %s", tc.Desc), func(t *testing.T) {
			t.Parallel()

			store := &mock.AuthStore{}
			s := makeAuthAPIServer(store)

			body := structToJSONReader(is, tc.AuthReq)

			request, err := http.NewRequest(http.MethodPost, "/api/v1/auths/signin", body)
			is.NoErr(err) // should be a successful request

			response := httptest.NewRecorder()

			s.ServeHTTP(response, request)

			is.Equal(response.Code, tc.StatusCode) // should not sign up user account with a 400 response

			var got errors.Error
			err = json.NewDecoder(response.Body).Decode(&got)

			is.NoErr(err)                           // should have a response
			is.Equal(got.StatusCode, tc.StatusCode) // shoud have error code

			is.Equal(errors.ToStatusCode(got), tc.StatusCode)    // shoud have error code
			is.True(strings.Contains(errors.ToErr(got), tc.Err)) // should have error message
			is.True(!store.CreateInvoked)                        // auth store sign up should not have been invoked
			is.True(!store.FindByEmailInvoked)                   // auth store FindByEmail should not have been invoked
		})
	}
}

func structToJSONReader(is *is.I, data any) io.Reader {
	bodyBytes, err := json.Marshal(data)
	is.NoErr(err)

	return bytes.NewBuffer(bodyBytes)
}
