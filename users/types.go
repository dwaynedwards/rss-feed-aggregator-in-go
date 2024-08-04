package users

import (
	"net/http"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Email    string
	Password string
	Name     string
}

type SignUpUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type SignUpUserResponse struct{}

type SignInUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInUserResponse struct{}

type UsersServer interface {
	RegisterEndpoints(*http.ServeMux)
}

type UsersService interface {
	SignUpUser(*SignUpUserRequest) (*SignUpUserResponse, error)
	SignInUser(*SignInUserRequest) (*SignInUserResponse, error)
}

type UsersStore interface {
	InsertUser(*User) bool
	GetUserByID(uuid.UUID) *User
	GetUserByEmail(string) *User
}

const (
	ErrEmailRequired    = "email is a required field"
	ErrPasswordRequired = "password is a required field"
	ErrNameRequired     = "name is a required field"
)
