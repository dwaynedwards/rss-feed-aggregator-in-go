package server

import (
	"net/http"
)

func (s *Server) routes() {
	router := http.NewServeMux()

	router.Handle("GET /healthz", s.handleHealthCheck())

	s.Handler = router
}

const HealthCheckResponseMsg = "Health check ok!"

func (s *Server) handleHealthCheck() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, HealthCheckResponseMsg, PlainText, http.StatusOK)
	})
}
