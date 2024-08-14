package http

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/service"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/store/postgres"
)

const ShutdownTimeout = 1 * time.Second

type APIServer struct {
	listener net.Listener
	server   *http.Server
	router   *http.ServeMux
	db       rf.DB

	Domain string

	AuthService rf.AuthService
	FeedService rf.FeedService
}

func newAPIServer(db rf.DB) *APIServer {
	s := &APIServer{
		server: &http.Server{},
		router: http.NewServeMux(),
		db:     db,
	}

	s.server.Handler = http.HandlerFunc(s.ServeHTTP)
	s.server.Handler = reportPanic(http.HandlerFunc(s.ServeHTTP))

	s.registerAuthRoutes(s.router)
	s.registerFeedRoutes(s.router)

	return s
}

func NewPostgresAPIServer() *APIServer {
	db := postgres.NewDB(rf.Config.DatabaseURL)

	s := newAPIServer(db)

	authStore := postgres.NewAuthStore(db)
	feedStore := postgres.NewFeedStore(db)
	s.AuthService = service.NewAuthService(authStore)
	s.FeedService = service.NewFeedService(feedStore)

	return s
}

func (s *APIServer) handleAuthRequired(next APIFunc) APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		token, err := rf.Read(r, "token")
		if err != nil {
			return err
		}

		userID, err := rf.ParseAndVerifyUserIDJWT(token)
		if err != nil {
			return err
		}

		if userID == 0 {
			return rf.NewAppError(rf.ECUnautherized, rf.EMUnauthorized)
		}

		r = rf.SetUserIDToRequestContext(r, userID)

		return next(w, r)
	}
}

func reportPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("panic: %v", err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) Open() (err error) {
	if err := s.db.Open(); err != nil {
		return err
	}

	if s.listener, err = net.Listen("tcp", ":"+rf.Config.APIPort); err != nil {
		return err
	}

	go s.server.Serve(s.listener)

	return nil
}

func (s *APIServer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	return s.db.Close()
}

func (s *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *APIServer) Scheme() string {
	return "http"
}

func (s *APIServer) Port() int {
	if s.listener == nil {
		return 0
	}
	return s.listener.Addr().(*net.TCPAddr).Port
}

func (s *APIServer) URL() string {
	scheme, port := s.Scheme(), s.Port()

	domain := "localhost"
	if s.Domain != "" {
		domain = s.Domain
	}

	// Return without port if using standard ports.
	if scheme == "http" && port == 80 {
		return fmt.Sprintf("%s://%s", s.Scheme(), domain)
	}
	return fmt.Sprintf("%s://%s:%d", s.Scheme(), domain, s.Port())
}
