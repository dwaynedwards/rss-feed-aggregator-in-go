package http

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
)

const ShutdownTimeout = 1 * time.Second

type APIServer struct {
	listener net.Listener
	server   *http.Server
	router   *http.ServeMux

	AuthService rf.AuthService
}

func NewAPIServer() *APIServer {
	s := &APIServer{
		server: &http.Server{},
		router: http.NewServeMux(),
	}

	s.registerAuthRoutes(s.router)

	return s
}

func (s *APIServer) Open() (err error) {
	if s.listener, err = net.Listen("tcp", ":"+os.Getenv("API_PORT")); err != nil {
		return err
	}

	go s.server.Serve(s.listener)

	return nil
}

func (s *APIServer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func (s *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
