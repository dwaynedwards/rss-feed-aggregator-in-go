package account

import (
	"net/http"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/common"
)

type server struct {
	service AccountService
}

func NewServer(service AccountService) *server {
	s := new(server)

	s.service = service

	return s
}

func (s *server) RegisterEndpoints(r *http.ServeMux) {
	s.makeV1Router(r)
}

func (s *server) makeV1Router(r *http.ServeMux) {
	r.Handle("GET /api/v1/accounts/healthz", common.MakeHTTPHandlerFunc(s.handleHealthCheck()))

	r.Handle("POST /api/v1/accounts", common.MakeHTTPHandlerFunc(s.handleAccountCreate()))
	r.Handle("POST /api/v1/accounts/signin", common.MakeHTTPHandlerFunc(s.handleAccountSignin()))
}

func (s *server) handleHealthCheck() common.APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		common.RespondWithText(w, http.StatusOK, common.HealthCheckResponseMsg)
		return nil
	}
}

func (s *server) handleAccountCreate() common.APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		req, err := getCreateRequestFromBody(w, r)
		if err != nil {
			return err
		}

		res, err := s.service.CreateAccount(req)
		if err != nil {
			return err
		}

		return common.RespondWithJSON(w, http.StatusCreated, res)
	}
}

func (s *server) handleAccountSignin() common.APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		req, err := getSigninRequestFromBody(w, r)
		if err != nil {
			return err
		}

		res, err := s.service.SigninAccount(req)
		if err != nil {
			return err
		}

		return common.RespondWithJSON(w, http.StatusOK, res)
	}
}
