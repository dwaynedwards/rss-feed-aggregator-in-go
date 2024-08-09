package rf

import (
	"context"
	"time"
)

type Auth struct {
	ID int64 `json:"id"`

	UserID int64 `json:"userID"`
	User   *User `json:"user"`

	BasicAuth *BasicAuth `json:"basicAuth"`

	Enabled        bool      `json:"enabled"`
	Deleted        bool      `json:"deleted"`
	CreatedAt      time.Time `json:"createdAt"`
	ModifiedAt     time.Time `json:"modifiedAt"`
	LastSignedInAt time.Time `json:"lastSignedInAt"`
}

type BasicAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpAuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type SignUpAuthResponse struct {
	Token string `json:"token"`
}

type SignInAuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInAuthResponse struct {
	Token string `json:"token"`
}

type AuthStore interface {
	Create(ctx context.Context, auth *Auth) error
	FindByEmail(ctx context.Context, email string) (*Auth, error)
}

type AuthService interface {
	SignUp(ctx context.Context, auth *Auth) (string, error)
	SignIn(ctx context.Context, auth *Auth) (string, error)
}
