package server

import (
	"net/http"
)

// Server struct
type Server struct {
	http.Handler
}

// NewServer creates a new server and returns pointer to it
func NewServer() *Server {
	s := new(Server)

	s.routes()

	return s
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, data string, contentType ContentType, statusCode int) {
	w.Header().Set("Content-Type", contentType.value)
	w.WriteHeader(statusCode)
	w.Write([]byte(data))
}
