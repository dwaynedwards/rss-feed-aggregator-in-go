package mock

import (
	"context"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
)

type AuthFailureCase struct {
	Desc string
	Auth *rf.Auth
	Code string
	Msg  string
}

var AuthWithMissingEmail = rf.NewAuthBuilder().
	WithUser(rf.NewUserBuilder().WithName("Gopher")).
	WithBasicAuth(rf.NewBasicAuthBuilder().WithPassword("gogopher1")).
	Build()

var AuthWithMissingPassword = rf.NewAuthBuilder().
	WithUser(rf.NewUserBuilder().WithName("Gopher")).
	WithBasicAuth(rf.NewBasicAuthBuilder().WithEmail("gopher1@go.com")).
	Build()

var AuthWithMissingName = rf.NewAuthBuilder().
	WithBasicAuth(rf.NewBasicAuthBuilder().
		WithEmail("gopher1@go.com").
		WithPassword("gogopher1")).
	Build()

type (
	optionsStoreFunc   func(c *AuthStore)
	optionsServiceFunc func(c *AuthService)
)

type AuthStore struct {
	CreateFn           func(ctx context.Context, auth *rf.Auth) error
	CreateInvoked      bool
	FindByEmailFn      func(ctx context.Context, email string) (*rf.Auth, error)
	FindByEmailInvoked bool
}

func NewAuthStore(opts ...optionsStoreFunc) *AuthStore {
	store := &AuthStore{}

	for _, opt := range opts {
		opt(store)
	}

	return store
}

func WithCreate(fn func(ctx context.Context, auth *rf.Auth) error) optionsStoreFunc {
	return func(c *AuthStore) {
		c.CreateFn = fn
	}
}

func WithFindByEmail(fn func(ctx context.Context, email string) (*rf.Auth, error)) optionsStoreFunc {
	return func(c *AuthStore) {
		c.FindByEmailFn = fn
	}
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
	SignUpFn func(ctx context.Context, auth *rf.Auth) error
	SignInFn func(ctx context.Context, id int64) error
}

func NewAuthService(opts ...optionsServiceFunc) *AuthService {
	service := &AuthService{}

	for _, opt := range opts {
		opt(service)
	}

	return service
}

func (a *AuthService) SignUp(ctx context.Context, auth *rf.Auth) error {
	return a.SignUpFn(ctx, auth)
}

func (a *AuthService) SignIn(ctx context.Context, id int64) error {
	return a.SignInFn(ctx, id)
}
