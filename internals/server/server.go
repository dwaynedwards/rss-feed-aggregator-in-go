package server

import (
	"encoding/json"
	"log"
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

func (s *Server) respondWithText(w http.ResponseWriter, r *http.Request, data string, statusCode int) {
	s.respond(w, []byte(data), ContentTypePlainText, statusCode)
}

func (s *Server) respondWithJSON(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	mData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", data)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respond(w, mData, ContentTypeJSON, statusCode)
}

func (s *Server) respond(w http.ResponseWriter, data []byte, contentType ContentType, statusCode int) {
	w.Header().Set("Content-Type", contentType.value)
	w.WriteHeader(statusCode)
	w.Write(data)
}
