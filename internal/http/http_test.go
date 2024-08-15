package http_test

import (
	rfhttp "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/http"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/service/authservice"
)

type APIServer struct {
	*rfhttp.APIServer
}

func makeAuthAPIServer(store authservice.AuthStore) *APIServer {
	s := &APIServer{
		APIServer: rfhttp.NewPostgresAPIServer(),
	}

	s.AuthService = authservice.NewAuthService(store)

	return s
}
