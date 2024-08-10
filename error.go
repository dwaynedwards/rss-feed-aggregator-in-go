package rf

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	ECIntenal      = "internal"
	ECInvalid      = "invalid"
	ECUnautherized = "unauthorized"
)

const (
	EMEmailRequired     = "email required."
	EMPasswordRequired  = "password required."
	EMNameRequired      = "name required."
	EMUserExists        = "user exists."
	EMInvlidCredentials = "invalid email and/or password was provided."
)

type AppError struct {
	Code string
	Msg  any
}

func (e AppError) Error() string {
	return msgToErrorString("app", e.Msg)
}

func NewAppError(code string, errs any) *AppError {
	return &AppError{
		Code: code,
		Msg:  errs,
	}
}

func AppErrorCode(err error) string {
	var e *AppError
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	}
	return ECIntenal
}

func AppErrorMessage(err error) string {
	var e *AppError
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return msgToString(e.Msg)
	}
	return "internal error"
}

type APIError struct {
	StatusCode int `json:"statusCode"`
	Msg        any `json:"msg"`
}

func (e APIError) Error() string {
	return msgToErrorString("api", e.Msg)
}

func NewAPIError(statusCode int, errs any) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Msg:        errs,
	}
}

func APIErrorCode(err error) int {
	var e *APIError
	if err == nil {
		return 0
	} else if errors.As(err, &e) {
		return e.StatusCode
	}
	return http.StatusInternalServerError
}

func APIErrorMessage(err error) string {
	var e *APIError
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return msgToString(e.Msg)
	}
	return "internal error"
}

func msgToErrorString(errorType string, msg any) string {
	return fmt.Sprintf("%s error: %v", errorType, msgToString(msg))
}

func msgToString(msg any) string {
	if m, ok := msg.(map[string]any); ok {
		var values []string
		for _, value := range m {
			values = append(values, value.(string))
		}
		return strings.Join(values, ", ")
	}
	return fmt.Sprintf("%v", msg)
}
