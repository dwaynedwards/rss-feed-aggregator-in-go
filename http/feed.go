package http

import (
	"net/http"
)

func (s *APIServer) registerFeedRoutes(r *http.ServeMux) {
	r.Handle("POST /api/v1/feeds/new", makeHTTPHandlerFunc(s.handleAuthRequired(s.handleFeedNew())))
}

func (s *APIServer) handleFeedNew() APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}
