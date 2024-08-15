package rf

import (
	"time"
)

type Auth struct {
	ID             int64     `json:"id"`
	Enabled        bool      `json:"enabled"`
	Deleted        bool      `json:"deleted"`
	CreatedAt      time.Time `json:"createdAt"`
	ModifiedAt     time.Time `json:"modifiedAt"`
	LastSignedInAt time.Time `json:"lastSignedInAt"`

	UserID int64 `json:"userID"`
	User   *User `json:"user"`

	BasicAuth *BasicAuth `json:"basicAuth"`

	Token string `json:"token"`
}

type BasicAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
