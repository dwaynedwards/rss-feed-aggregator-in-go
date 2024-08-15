package mock

import (
	"context"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/builder"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/errors"
)

type SignUpFailureCase struct {
	Desc    string
	Req     *rf.SignUpRequest
	RefCode errors.ReferenceCode
	Err     string
}

type SignInFailureCase struct {
	Desc    string
	Req     *rf.SignInRequest
	RefCode errors.ReferenceCode
	Err     string
}

var SignUpWithMissingEmail = builder.NewSignUpRequestBuilder().
	WithName("Gopher").
	WithPassword("gogopher1").
	Build()

var SignUpWithMissingPassword = builder.NewSignUpRequestBuilder().
	WithName("Gopher").
	WithEmail("gopher1@go.com").
	Build()

var SignUpWithMissingName = builder.NewSignUpRequestBuilder().
	WithEmail("gopher1@go.com").
	WithPassword("gogopher1").
	Build()

var SignInWithMissingEmail = builder.NewSignInRequestBuilder().
	WithPassword("gogopher1").
	Build()

var SignInWithMissingPassword = builder.NewSignInRequestBuilder().
	WithEmail("gopher1@go.com").
	Build()

type AuthAPIFailureCase struct {
	Desc       string
	AuthReq    any
	StatusCode int
	Err        string
}

var SignUpAuthAPIWithMissingEmail = builder.NewSignUpRequestBuilder().
	WithPassword("password1").
	WithName("Gopher").
	Build()

var SignUpAuthAPIWithMissingPassword = builder.NewSignUpRequestBuilder().
	WithEmail("gopher@go.com").
	WithName("Gopher").
	Build()

var SignUpAuthAPIWithMissingName = builder.NewSignUpRequestBuilder().
	WithEmail("gopher1@go.com").
	WithPassword("gogopher1").
	Build()

var SignInAuthAPIWithMissingEmail = builder.NewSignInRequestBuilder().
	WithPassword("password1").
	Build()

var SignInAuthAPIWithMissingPassword = builder.NewSignInRequestBuilder().
	WithEmail("gopher@go.com").
	Build()

type AuthStore struct {
	CreateAuthAndUserFn func(ctx context.Context, auth *rf.Auth) error
	CreateInvoked       bool
	FindByEmailFn       func(ctx context.Context, email string) (*rf.Auth, error)
	FindByEmailInvoked  bool
}

func (as *AuthStore) CreateAuthAndUser(ctx context.Context, auth *rf.Auth) error {
	as.CreateInvoked = true
	return as.CreateAuthAndUserFn(ctx, auth)
}

func (as *AuthStore) FindByEmail(ctx context.Context, email string) (*rf.Auth, error) {
	as.FindByEmailInvoked = true
	return as.FindByEmailFn(ctx, email)
}
