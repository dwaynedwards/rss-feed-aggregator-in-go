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
	CreateFn           func(ctx context.Context, auth *rf.Auth) error
	CreateInvoked      bool
	FindByEmailFn      func(ctx context.Context, email string) (*rf.Auth, error)
	FindByEmailInvoked bool
}

func (a *AuthStore) Create(ctx context.Context, auth *rf.Auth) error {
	a.CreateInvoked = true
	return a.CreateFn(ctx, auth)
}

func (a *AuthStore) FindByEmail(ctx context.Context, email string) (*rf.Auth, error) {
	a.FindByEmailInvoked = true
	return a.FindByEmailFn(ctx, email)
}

type AuthService struct {
	SignUpFn      func(ctx context.Context, auth *rf.Auth) error
	SignUpInvoked bool
	SignInFn      func(ctx context.Context, id int64) error
	SignInInvoked bool
}

func (a *AuthService) SignUp(ctx context.Context, auth *rf.Auth) error {
	a.SignUpInvoked = true
	return a.SignUpFn(ctx, auth)
}

func (a *AuthService) SignIn(ctx context.Context, id int64) error {
	a.SignInInvoked = true
	return a.SignInFn(ctx, id)
}
