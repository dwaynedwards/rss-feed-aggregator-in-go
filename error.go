package rssfeed

import (
	"errors"
	"fmt"
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
	if m, ok := e.Msg.(map[string]string); ok {
		var values []string
		for _, value := range m {
			values = append(values, value)
		}
		return fmt.Sprintf("app error: %v", strings.Join(values, ", "))
	}
	return fmt.Sprintf("app error: %v", e.Msg)
}

func NewAppError(code string, err error) *AppError {
	return &AppError{
		Code: code,
		Msg:  err.Error(),
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
		return fmt.Sprint(e.Msg)
	}
	return "Internal error."
}

func AppErrorf(code string, fmtstring string, args ...any) *AppError {
	return &AppError{
		Code: code,
		Msg:  fmt.Sprintf(fmtstring, args...),
	}
}
