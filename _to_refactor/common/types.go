package common

import "net/http"

type APIFunc func(http.ResponseWriter, *http.Request) error
