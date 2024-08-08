package service_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/mock"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/service"
	"github.com/matryer/is"
)

func TestAuthService_SignUp_Success(t *testing.T) {
	is := is.New(t)

	t.Run("Should succeed with sign up", func(t *testing.T) {
		store := &mock.AuthStore{
			CreateFn: func(ctx context.Context, auth *rf.Auth) error {
				auth.ID = 1
				auth.UserID = 1
				auth.CreatedAt = time.Now()
				auth.ModifiedAt = time.Now()
				auth.LastSignedInAt = time.Now()
				return nil
			},
		}

		service := service.NewAuthService(store)

		auth := rf.NewAuthBuilder().
			WithUser(rf.NewUserBuilder().WithName("Gopher")).
			WithBasicAuth(rf.NewBasicAuthBuilder().
				WithEmail("gopher1@go.com").
				WithPassword("gogopher1")).
			Build()

		token, err := service.SignUp(context.Background(), auth)

		is.NoErr(err)                          // should be signed up
		is.True(len(token) > 0)                // should receive token
		is.Equal(auth.UserID, int64(1))        // auth UserID should be 1
		is.Equal(auth.ID, int64(1))            // auth ID should be 1
		is.True(!auth.CreatedAt.IsZero())      // auth CreatedAt should be set
		is.True(!auth.ModifiedAt.IsZero())     // auth ModifiedAt should be set
		is.True(!auth.LastSignedInAt.IsZero()) // auth LastLoggedInAt should be set
		is.True(store.CreateInvoked)           // auth store Create should have been invoked
		is.True(!store.FindByEmailInvoked)     // auth store FindByEmail should not have been invoked
	})
}

func TestAuthService_SignUp_Failure(t *testing.T) {
	is := is.New(t)

	signUpFailureCases := []mock.AuthFailureCase{
		{Desc: "with missing email", Auth: mock.AuthWithMissingEmail, Code: rf.ECInvalid, Msg: rf.EMEmailRequired},
		{Desc: "with missing password", Auth: mock.AuthWithMissingPassword, Code: rf.ECInvalid, Msg: rf.EMPasswordRequired},
		{Desc: "with missing name", Auth: mock.AuthWithMissingName, Code: rf.ECInvalid, Msg: rf.EMNameRequired},
	}
	for _, tc := range signUpFailureCases {
		t.Run(fmt.Sprintf("Should fail to sign up %s", tc.Desc), func(t *testing.T) {
			store := &mock.AuthStore{}
			service := service.NewAuthService(store)

			_, err := service.SignUp(context.Background(), tc.Auth)

			is.True(err != nil)                                        // should be an error
			is.Equal(rf.AppErrorCode(err), tc.Code)                    // shoud have error code
			is.True(strings.Contains(rf.AppErrorMessage(err), tc.Msg)) // should have error message
			is.True(!store.CreateInvoked)                              // auth store Create should not have been invoked
			is.True(!store.FindByEmailInvoked)                         // auth store FindByEmail should not have been invoked
		})
	}
}

func TestAuthService_SignIn_Success(t *testing.T) {
	is := is.New(t)

	t.Run("Should succeed with sign in", func(t *testing.T) {
		password := "gogopher1"
		hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
		is.NoErr(err) // Should hash password

		store := &mock.AuthStore{
			FindByEmailFn: func(ctx context.Context, email string) (*rf.Auth, error) {
				auth := rf.NewAuthBuilder().
					WithUserID(1).
					WithBasicAuth(rf.NewBasicAuthBuilder().WithPassword(hashedPassword)).
					Build()
				return auth, nil
			},
		}

		service := service.NewAuthService(store)

		authSignIn := rf.NewAuthBuilder().
			WithBasicAuth(rf.NewBasicAuthBuilder().
				WithEmail("gopher1@go.com").
				WithPassword(password)).
			Build()

		token, err := service.SignIn(context.Background(), authSignIn)

		is.NoErr(err)                         // should be signed in
		is.True(len(token) > 0)               // should receive token
		is.Equal(authSignIn.UserID, int64(1)) // auth UserID should be 1
		is.True(!store.CreateInvoked)         // auth store Create should not have been invoked
		is.True(store.FindByEmailInvoked)     // auth store FindByEmail should  have been invoked
	})
}

func TestAuthService_SignIn_Failure(t *testing.T) {
	is := is.New(t)

	signUpFailureCases := []mock.AuthFailureCase{
		{Desc: "with missing email", Auth: mock.AuthWithMissingEmail, Code: rf.ECInvalid, Msg: rf.EMEmailRequired},
		{Desc: "with missing password", Auth: mock.AuthWithMissingPassword, Code: rf.ECInvalid, Msg: rf.EMPasswordRequired},
	}
	for _, tc := range signUpFailureCases {
		t.Run(fmt.Sprintf("Should fail to sign in %s", tc.Desc), func(t *testing.T) {
			store := &mock.AuthStore{}
			service := service.NewAuthService(store)

			_, err := service.SignIn(context.Background(), tc.Auth)

			is.True(err != nil)                                        // should be an error
			is.Equal(rf.AppErrorCode(err), tc.Code)                    // shoud have  error code
			is.True(strings.Contains(rf.AppErrorMessage(err), tc.Msg)) // should have error message
			is.True(!store.CreateInvoked)                              // auth store Create should not have been invoked
			is.True(!store.FindByEmailInvoked)                         // auth store FindByEmail should not have been invoked
		})
	}
}
