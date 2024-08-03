package account

import (
	"net/http"

	"github.com/google/uuid"
)

type Account struct {
	ID       uuid.UUID
	Email    string
	Password string
	Name     string
}

type CreateAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type CreateAccountResponse struct{}

type SigninAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninAccountResponse struct{}

type AccountServer interface {
	RegisterEndpoints(*http.ServeMux)
}

type AccountService interface {
	CreateAccount(*CreateAccountRequest) (*CreateAccountResponse, error)
	SigninAccount(*SigninAccountRequest) (*SigninAccountResponse, error)
}

type AccountStore interface {
	Insert(*Account) bool
	GetByID(uuid.UUID) *Account
	GetByEmail(string) *Account
}

const (
	ErrEmailRequired    = "email is a required field"
	ErrPasswordRequired = "password is a required field"
	ErrNameRequired     = "name is a required field"
)
