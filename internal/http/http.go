package http

import (
	"errors"
	"log/slog"
	"net/http"

	rferrors "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/errors"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/response"
)

type APIFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandlerFunc(h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			var werr error
			var e rferrors.Error
			if errors.As(err, &e) {
				werr = response.WriteJSON(w, e.StatusCode, e)
			} else {
				errRes := rferrors.InternalServerError("internal server error")
				werr = response.WriteJSON(w, errRes.StatusCode, errRes)
			}
			slog.Error("HTTP API error", "err", err.Error(), "path", r.URL.Path)
			if werr == nil {
				return
			}
			slog.Error("HTTP API error", "werr", werr.Error(), "path", r.URL.Path)
		}
	}
}
