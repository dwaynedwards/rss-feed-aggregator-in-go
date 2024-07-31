package data

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internals/util"
)

// User struct
type User struct {
	ID       int
	Email    string
	Password string
	Name     string
}

// CreateUserRequest struct
type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// CreateUserResponse struct
type CreateUserResponse struct {
	ID int `json:"id"`
}

// Error constants
const (
	ErrEmailRequired          = "email is a required field"
	ErrPasswordRequired       = "password is a required field"
	ErrNameRequired           = "name is a required field"
	ErrUnableToProcessRequest = "unable to process request body: %s"
)

// GetCreateUserRequestFromBody gets a CreateUserRequest
func GetCreateUserRequestFromBody(w http.ResponseWriter, r *http.Request) (CreateUserRequest, error) {
	var requestData CreateUserRequest

	if err := util.DecodeJSONBody(w, r, &requestData); err != nil {
		return requestData, err
	}

	if err := validateCreateUserRequest(requestData); err != nil {
		return requestData, err
	}

	return requestData, nil
}

// GetCreateUserResponseFromUser gets a CreateUserResponse
func GetCreateUserResponseFromUser(user User) CreateUserResponse {
	return CreateUserResponse{
		ID: user.ID,
	}
}

// GetUserFromCreateUserRequestWithID gets a User
func GetUserFromCreateUserRequestWithID(req CreateUserRequest, id int) User {
	return User{
		ID:       id,
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
	}
}

func validateCreateUserRequest(req CreateUserRequest) error {
	var errs []string

	if req.Email == "" {
		errs = append(errs, ErrEmailRequired)
	}

	if req.Password == "" {
		errs = append(errs, ErrPasswordRequired)
	}

	if req.Name == "" {
		errs = append(errs, ErrNameRequired)
	}

	if len(errs) > 0 {
		return fmt.Errorf(ErrUnableToProcessRequest, strings.Join(errs, ", "))
	}

	return nil
}
