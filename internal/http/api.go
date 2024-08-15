package http

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal"
	rfcontext "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/context"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/cookie"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/errors"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/jwt"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/service/authservice"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/service/feedservice"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/store/postgresstore"
)

const ShutdownTimeout = 1 * time.Second

type AuthService interface {
	SignUp(ctx context.Context, req *rf.SignUpRequest) (string, error)
	SignIn(ctx context.Context, req *rf.SignInRequest) (string, error)
}

type FeedService interface {
	AddFeed(ctx context.Context, req *rf.AddFeedRequest) (int64, error)
	RemoveFeed(ctx context.Context, feedID int64) error
	GetFeeds(ctx context.Context) ([]rf.Feed, error)
	GetFeed(ctx context.Context, feedID int64) (*rf.Feed, error)
}

type DB interface {
	Open() error
	Close() error
}

type APIServer struct {
	listener net.Listener
	server   *http.Server
	router   *http.ServeMux
	db       DB

	Domain string

	AuthService AuthService
	FeedService FeedService
}

func NewAPIServer(db DB) *APIServer {
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
	db := postgresstore.NewDB(rf.Config.DatabaseURL)

	s := NewAPIServer(db)

	authStore := postgresstore.NewAuthStore(db)
	feedStore := postgresstore.NewFeedStore(db)
	s.AuthService = authservice.NewAuthService(authStore)
	s.FeedService = feedservice.NewFeedService(feedStore)

	return s
}

func (s *APIServer) handleAuthRequired(next APIFunc) APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		token, err := cookie.Read(r, "token")
		if err != nil {
			return err
		}

		userID, err := jwt.ParseAndVerifyUserID(token)
		if err != nil {
			return err
		}

		if userID == 0 {
			return errors.Unauthorizedf(errors.ErrUnauthorized)
		}

		r = rfcontext.SetUserIDToRequestContext(r, userID)

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
