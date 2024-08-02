package common

import "net/http"

type ContentType struct {
	Value string
}

type APIFunc func(http.ResponseWriter, *http.Request) error

var (
	ContentTypePlainText = ContentType{"text/plain"}
	ContentTypeJSON      = ContentType{"application/json"}
)

const HealthCheckResponseMsg = "Health check ok!"
