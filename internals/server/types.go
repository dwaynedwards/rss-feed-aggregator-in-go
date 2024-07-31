package server

import "net/http"

type ContentType struct {
	Value string
}

type apiFunc func(http.ResponseWriter, *http.Request) error

var (
	ContentTypePlainText = ContentType{"text/plain"}
	ContentTypeJSON      = ContentType{"application/json"}
)
