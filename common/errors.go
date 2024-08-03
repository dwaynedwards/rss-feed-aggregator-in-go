package common

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int `json:"statusCode"`
	Msg        any `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d", e.StatusCode)
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
}

func InvalidRequestData(errs map[string]string) APIError {
	return APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Msg:        errs,
	}
}

func InvalidJSON(err error) APIError {
	return NewAPIError(http.StatusBadRequest, err)
}

func InvalidCredentials() APIError {
	return NewAPIError(http.StatusUnauthorized, fmt.Errorf("invalid email or password was provided"))
}

// InvalidAccountExists nornally this workflow would be handled with a status 201 and a message saying an email was sent to
// verify the account. When this error is hit, an email would be sent saying if you're trying to create
// an you can trying executing the forgot password workflow instead of leaking internal info to the user
// that an account already exists with the email provided, but this is outside of the scope of this project
func InvalidAccountExists() APIError {
	return NewAPIError(http.StatusConflict, fmt.Errorf("account already exists"))
}
