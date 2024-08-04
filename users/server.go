package users

import (
	"net/http"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/common"
)

type server struct {
	service UsersService
}

func NewServer(service UsersService) *server {
	s := new(server)

	s.service = service

	return s
}

func (s *server) RegisterEndpoints(r *http.ServeMux) {
	s.makeV1Router(r)
}

func (s *server) makeV1Router(r *http.ServeMux) {
	r.Handle("GET /api/v1/users/status", common.MakeHTTPHandlerFunc(s.handleStatusCheck()))

	r.Handle("POST /api/v1/users/signup", common.MakeHTTPHandlerFunc(s.handleUserSignUp()))
	r.Handle("POST /api/v1/users/signin", common.MakeHTTPHandlerFunc(s.handleUserSignIn()))
}

func (s *server) handleStatusCheck() common.APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		common.WriteJSON(w, http.StatusOK, map[string]string{"msg": "Health check ok!"})
		return nil
	}
}

func (s *server) handleUserSignUp() common.APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		req, err := getCreateRequestFromBody(w, r)
		if err != nil {
			return err
		}

		res, err := s.service.SignUpUser(req)
		if err != nil {
			return err
		}

		return common.WriteJSON(w, http.StatusCreated, res)
	}
}

func (s *server) handleUserSignIn() common.APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		req, err := getSigninRequestFromBody(w, r)
		if err != nil {
			return err
		}

		res, err := s.service.SignInUser(req)
		if err != nil {
			return err
		}

		return common.WriteJSON(w, http.StatusOK, res)
	}
}
