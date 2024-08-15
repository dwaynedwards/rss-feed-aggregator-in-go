package authservice_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/builder"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/errors"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/mock"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/service/authservice"
	"github.com/matryer/is"
)

func TestAuthService_SignUp_Success(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	t.Run("Should succeed with sign up", func(t *testing.T) {
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

		service := authservice.NewAuthService(store)

		req := builder.NewSignUpRequestBuilder().
			WithName("Gopher").
			WithEmail("gopher1@go.com").
			WithPassword("gogopher1").
			Build()

		token, err := service.SignUp(context.Background(), req)

		is.NoErr(err)                     // should be signed up
		is.True(len(token) > 0)           // should receive token
		is.True(store.CreateInvoked)      // auth store Create should have been invoked
		is.True(store.FindByEmailInvoked) // auth store FindByEmail should have been invoked
	})
}

func TestAuthService_SignUp_Failure(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	signUpFailureCases := []mock.SignUpFailureCase{
		{Desc: "with missing email", Req: mock.SignUpWithMissingEmail, RefCode: errors.InvalidData, Err: errors.ErrEmailRequired},
		{Desc: "with missing password", Req: mock.SignUpWithMissingPassword, RefCode: errors.InvalidData, Err: errors.ErrPasswordRequired},
		{Desc: "with missing name", Req: mock.SignUpWithMissingName, RefCode: errors.InvalidData, Err: errors.ErrNameRequired},
	}
	for _, tc := range signUpFailureCases {
		t.Run(fmt.Sprintf("Should fail to sign up %s", tc.Desc), func(t *testing.T) {
			t.Parallel()

			store := &mock.AuthStore{}
			service := authservice.NewAuthService(store)

			_, err := service.SignUp(context.Background(), tc.Req)

			is.True(err != nil)                                  // should be an error
			is.Equal(errors.ToReferenceCode(err), tc.RefCode)    // shoud have error code
			is.True(strings.Contains(errors.ToErr(err), tc.Err)) // should have error message
			is.True(!store.CreateInvoked)                        // auth store Create should not have been invoked
			is.True(!store.FindByEmailInvoked)                   // auth store FindByEmail should not have been invoked
		})
	}
}

func TestAuthService_SignIn_Success(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	t.Run("Should succeed with sign in", func(t *testing.T) {
		t.Parallel()

		password := "gogopher1"
		hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
		is.NoErr(err) // Should hash password

		store := &mock.AuthStore{
			FindByEmailFn: func(ctx context.Context, email string) (*rf.Auth, error) {
				auth := builder.NewAuthBuilder().
					WithUserID(1).
					WithBasicAuth(builder.NewBasicAuthBuilder().WithPassword(hashedPassword)).
					Build()
				return auth, nil
			},
		}

		service := authservice.NewAuthService(store)

		req := builder.NewSignInRequestBuilder().
			WithEmail("gopher1@go.com").
			WithPassword(password).
			Build()

		token, err := service.SignIn(context.Background(), req)

		is.NoErr(err)                     // should be signed in
		is.True(len(token) > 0)           // should receive token
		is.True(!store.CreateInvoked)     // auth store Create should not have been invoked
		is.True(store.FindByEmailInvoked) // auth store FindByEmail should  have been invoked
	})
}

func TestAuthService_SignIn_Failure(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	signUpFailureCases := []mock.SignInFailureCase{
		{Desc: "with missing email", Req: mock.SignInWithMissingEmail, RefCode: errors.InvalidData, Err: errors.ErrEmailRequired},
		{Desc: "with missing password", Req: mock.SignInWithMissingPassword, RefCode: errors.InvalidData, Err: errors.ErrPasswordRequired},
	}
	for _, tc := range signUpFailureCases {
		t.Run(fmt.Sprintf("Should fail to sign in %s", tc.Desc), func(t *testing.T) {
			t.Parallel()

			store := &mock.AuthStore{}
			service := authservice.NewAuthService(store)

			_, err := service.SignIn(context.Background(), tc.Req)

			is.True(err != nil)                                  // should be an error
			is.Equal(errors.ToReferenceCode(err), tc.RefCode)    // shoud have  error code
			is.True(strings.Contains(errors.ToErr(err), tc.Err)) // should have error message
			is.True(!store.CreateInvoked)                        // auth store Create should not have been invoked
			is.True(!store.FindByEmailInvoked)                   // auth store FindByEmail should not have been invoked
		})
	}
}
