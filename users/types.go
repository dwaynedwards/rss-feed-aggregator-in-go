package users

import (
	"net/http"
)

type User struct {
	ID       int64
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
	InsertUser(*User) error
	GetUserByID(int64) (*User, error)
	GetUserByEmail(string) (*User, error)
}

const (
	ErrEmailRequired    = "email is a required field"
	ErrPasswordRequired = "password is a required field"
	ErrNameRequired     = "name is a required field"
)
