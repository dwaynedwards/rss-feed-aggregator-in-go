package mock

import (
	"context"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/builder"
)

type AuthFailureCase struct {
	Desc string
	Auth *rf.Auth
	Code string
	Msg  string
}

var AuthWithMissingEmail = builder.NewAuthBuilder().
	WithUser(builder.NewUserBuilder().WithName("Gopher")).
	WithBasicAuth(builder.NewBasicAuthBuilder().WithPassword("gogopher1")).
	Build()

var AuthWithMissingPassword = builder.NewAuthBuilder().
	WithUser(builder.NewUserBuilder().WithName("Gopher")).
	WithBasicAuth(builder.NewBasicAuthBuilder().WithEmail("gopher1@go.com")).
	Build()

var AuthWithMissingName = builder.NewAuthBuilder().
	WithBasicAuth(builder.NewBasicAuthBuilder().
		WithEmail("gopher1@go.com").
		WithPassword("gogopher1")).
	Build()

type AuthAPIFailureCase struct {
	Desc       string
	AuthReq    any
	StatusCode int
	Msg        string
}

var SignUpAuthAPIWithMissingEmail = builder.NewSignUpAuthRequestBuilder().
	WithPassword("password1").
	WithName("Gopher").
	Build()

var SignUpAuthAPIWithMissingPassword = builder.NewSignUpAuthRequestBuilder().
	WithEmail("gopher@go.com").
	WithName("Gopher").
	Build()

var SignUpAuthAPIWithMissingName = builder.NewSignUpAuthRequestBuilder().
	WithEmail("gopher1@go.com").
	WithPassword("gogopher1").
	Build()

var SignInAuthAPIWithMissingEmail = builder.NewSignInAuthRequestBuilder().
	WithPassword("password1").
	Build()

var SignInAuthAPIWithMissingPassword = builder.NewSignInAuthRequestBuilder().
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
