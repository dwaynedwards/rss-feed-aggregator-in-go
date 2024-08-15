package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type ReferenceCode int

const (
	Internal ReferenceCode = iota + 1
	MalformedData
	InvalidData
	Unauthorized
	NotFound
)

const (
	ErrEmailRequired    = "email required."
	ErrPasswordRequired = "password required."
	ErrNameRequired     = "name required."
	ErrURLRequired      = "url required."

	ErrCouldNotProcess    = "could not process request."
	ErrInvalidCredentials = "invalid email and/or password was provided."
	ErrUnauthorized       = "unauthorized to perform this action."

	ErrTokenExpired                 = "token expired"
	ErrTokenClaimsFailed            = "token claims failed"
	ErrTokenGenerationFailed        = "token generation failed"
	ErrTokenParseFailed             = "token parse failed"
	ErrTokenUnexpactedSigningMethod = "token unexpected signing method"
)

type Error struct {
	ReferenceCode ReferenceCode `json:"referenceCode"`
	StatusCode    int           `json:"statusCode"`
	Err           any           `json:"err,omitempty"`
}

func (e Error) Error() string {
	return errToString(e.Err)
}

func Errorf(refCode ReferenceCode, format string, args ...any) Error {
	return Error{
		ReferenceCode: refCode,
		Err:           fmt.Sprintf(format, args...),
	}
}

func InternalError(err any) Error {
	return Error{
		ReferenceCode: Internal,
		Err:           err,
	}
}

func InvalidError(err any) Error {
	return Error{
		ReferenceCode: InvalidData,
		Err:           err,
	}
}

func InternalServerError(err any) Error {
	return Error{
		ReferenceCode: Internal,
		StatusCode:    http.StatusInternalServerError,
		Err:           err,
	}
}

func MalformedDataError(err any) Error {
	return Error{
		ReferenceCode: MalformedData,
		StatusCode:    http.StatusUnprocessableEntity,
		Err:           err,
	}
}

func BadRequestError(err any) Error {
	return Error{
		ReferenceCode: InvalidData,
		StatusCode:    http.StatusBadRequest,
		Err:           err,
	}
}

func UnauthorizedError(err any) Error {
	return Error{
		ReferenceCode: Unauthorized,
		StatusCode:    http.StatusUnauthorized,
		Err:           err,
	}
}

func InternalErrorf(format string, args ...any) Error {
	return Errorf(Internal, format, args...)
}

func MalformedDataf(format string, args ...any) Error {
	return Errorf(MalformedData, format, args...)
}

func InvalidDataf(format string, args ...any) Error {
	return Errorf(InvalidData, format, args...)
}

func Unauthorizedf(format string, args ...any) Error {
	return Errorf(Unauthorized, format, args...)
}

func ToAPIError(err error) error {
	var e Error
	if err == nil {
		return nil
	} else if errors.As(err, &e) {
		switch e.ReferenceCode {
		case InvalidData:
			return BadRequestError(e.Err)
		case Unauthorized:
			return UnauthorizedError(e.Err)
		}
	}
	return err
}

func ToReferenceCode(err error) ReferenceCode {
	var e Error
	if err == nil {
		return 0
	} else if errors.As(err, &e) {
		return e.ReferenceCode
	}
	return Internal
}

func ToStatusCode(err error) int {
	var e Error
	if err == nil {
		return 0
	} else if errors.As(err, &e) {
		return e.StatusCode
	}
	return http.StatusInternalServerError
}

func ToErr(err error) string {
	var e Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return errToString(e.Err)
	}
	return "internal error"
}

func errToString(msg any) string {
	if m, ok := msg.(map[string]any); ok {
		var values []string
		for _, value := range m {
			values = append(values, value.(string))
		}
		return strings.Join(values, ", ")
	}
	return fmt.Sprintf("%v", msg)
}
