package server

import (
	"net/http"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internals/account"
)

type server struct {
	http.Handler
	accountService account.AccountService
}

func NewServer(accountService account.AccountService) *server {
	s := new(server)

	s.accountService = accountService

	s.routes()

	return s
}
