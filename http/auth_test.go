package http_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/builder"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/mock"
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
		s := makeAPIServer(store)

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

		is.Equal(response.Code, http.StatusCreated) // should sign up user account with a 201 response

		var got rf.SignUpAuthResponse
		err = json.NewDecoder(response.Body).Decode(&got)

		is.NoErr(err)               // should have a response
		is.True(len(got.Token) > 0) // should have a token
	})
}

func TestAuthService_SignUp_Failure(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	signUpFailureCases := []mock.AuthAPIFailureCase{
		{Desc: "with missing email", AuthReq: mock.SignUpAuthAPIWithMissingEmail, StatusCode: http.StatusBadRequest, Msg: rf.EMEmailRequired},
		{Desc: "with missing password", AuthReq: mock.SignUpAuthAPIWithMissingPassword, StatusCode: http.StatusBadRequest, Msg: rf.EMPasswordRequired},
		{Desc: "with missing name", AuthReq: mock.SignUpAuthAPIWithMissingName, StatusCode: http.StatusBadRequest, Msg: rf.EMNameRequired},
	}
	for _, tc := range signUpFailureCases {
		t.Run(fmt.Sprintf("Should fail to sign up %s", tc.Desc), func(t *testing.T) {
			t.Parallel()

			store := &mock.AuthStore{}
			s := makeAPIServer(store)

			body := jsonBodyReaderFromStruct(is, tc.AuthReq)

			request, err := http.NewRequest(http.MethodPost, "/api/v1/auths/signup", body)
			is.NoErr(err) // should be a successful request

			response := httptest.NewRecorder()

			s.ServeHTTP(response, request)

			is.Equal(response.Code, tc.StatusCode) // should not sign up user account with a 400 response

			var got *rf.APIError
			err = json.NewDecoder(response.Body).Decode(&got)

			is.NoErr(err)                           // should have a response
			is.Equal(got.StatusCode, tc.StatusCode) // shoud have error code

			is.Equal(rf.APIErrorCode(got), tc.StatusCode)              // shoud have error code
			is.True(strings.Contains(rf.APIErrorMessage(got), tc.Msg)) // should have error message
			is.True(!store.CreateInvoked)                              // auth store sign up should not have been invoked
			is.True(!store.FindByEmailInvoked)                         // auth store FindByEmail should not have been invoked
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

		s := makeAPIServer(store)

		req := builder.NewSignInAuthRequestBuilder().
			WithEmail("gopher@go.com").
			WithPassword(password).
			Build()
		body := jsonBodyReaderFromStruct(is, req)

		request, err := http.NewRequest(http.MethodPost, "/api/v1/auths/signin", body)
		is.NoErr(err) // should be a successful request

		response := httptest.NewRecorder()

		s.ServeHTTP(response, request)

		is.Equal(response.Code, http.StatusOK) // should sign in a user account with a 200 response

		var got rf.SignInAuthResponse
		err = json.NewDecoder(response.Body).Decode(&got)

		is.NoErr(err)               // should have a response
		is.True(len(got.Token) > 0) // should have a token
	})
}

func TestAuthService_SignIn_Failure(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	signUpFailureCases := []mock.AuthAPIFailureCase{
		{Desc: "with missing email", AuthReq: mock.SignInAuthAPIWithMissingEmail, StatusCode: http.StatusBadRequest, Msg: rf.EMEmailRequired},
		{Desc: "with missing password", AuthReq: mock.SignInAuthAPIWithMissingPassword, StatusCode: http.StatusBadRequest, Msg: rf.EMPasswordRequired},
	}
	for _, tc := range signUpFailureCases {
		t.Run(fmt.Sprintf("Should fail to sign up %s", tc.Desc), func(t *testing.T) {
			t.Parallel()

			store := &mock.AuthStore{}
			s := makeAPIServer(store)

			body := jsonBodyReaderFromStruct(is, tc.AuthReq)

			request, err := http.NewRequest(http.MethodPost, "/api/v1/auths/signin", body)
			is.NoErr(err) // should be a successful request

			response := httptest.NewRecorder()

			s.ServeHTTP(response, request)

			is.Equal(response.Code, tc.StatusCode) // should not sign up user account with a 400 response

			var got *rf.APIError
			err = json.NewDecoder(response.Body).Decode(&got)

			is.NoErr(err)                           // should have a response
			is.Equal(got.StatusCode, tc.StatusCode) // shoud have error code

			is.Equal(rf.APIErrorCode(got), tc.StatusCode)              // shoud have error code
			is.True(strings.Contains(rf.APIErrorMessage(got), tc.Msg)) // should have error message
			is.True(!store.CreateInvoked)                              // auth store sign up should not have been invoked
			is.True(!store.FindByEmailInvoked)                         // auth store FindByEmail should not have been invoked
		})
	}
}
