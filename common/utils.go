package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func MakeHTTPHandlerFunc(a APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer closeBody(r)

		if err := a(w, r); err != nil {
			var mr *MalformedRequestError
			var a *AccountError

			if errors.As(err, &mr) {
				http.Error(w, mr.Error(), mr.Status)
			} else if errors.As(err, &a) {
				http.Error(w, a.Error(), a.Status)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func closeBody(r *http.Request) {
	if r.Body == nil {
		return
	}
	io.Copy(io.Discard, r.Body)
	r.Body.Close()
}

func RespondWithText(w http.ResponseWriter, status int, data string) {
	w.Header().Set("Content-Type", ContentTypePlainText.Value)
	w.WriteHeader(status)
	w.Write([]byte(data))
}

func RespondWithJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", ContentTypeJSON.Value)
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// Source from: https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

type MalformedRequestError struct {
	Status int
	Msg    string
}

func (mr MalformedRequestError) Error() string {
	return mr.Msg
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			return &MalformedRequestError{Status: http.StatusUnsupportedMediaType, Msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &MalformedRequestError{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &MalformedRequestError{Status: http.StatusBadRequest, Msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &MalformedRequestError{Status: http.StatusBadRequest, Msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &MalformedRequestError{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &MalformedRequestError{Status: http.StatusBadRequest, Msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &MalformedRequestError{Status: http.StatusRequestEntityTooLarge, Msg: msg}

		default:
			return &MalformedRequestError{Status: http.StatusBadRequest, Msg: err.Error()}
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		return &MalformedRequestError{Status: http.StatusBadRequest, Msg: msg}
	}

	return nil
}
