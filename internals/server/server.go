package server

import (
	"net/http"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/account"
)

type server struct {
	http.Handler
}

func NewServer(accountServer account.AccountServer) *server {
	s := new(server)

	router := http.NewServeMux()
	accountServer.RegisterEndpoints(router)
	s.Handler = router

	return s
}
