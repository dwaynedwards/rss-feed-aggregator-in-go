package http

import (
	"net/http"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/cookie"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/errors"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/request"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/response"
)

func (s *APIServer) registerAuthRoutes(r *http.ServeMux) {
	r.Handle("POST /api/v1/auths/signup", makeHTTPHandlerFunc(s.handleSignUp()))
	r.Handle("POST /api/v1/auths/signin", makeHTTPHandlerFunc(s.handleSignIn()))
}

func (s *APIServer) handleSignUp() APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var req *rf.SignUpRequest

		if err := request.ReadJSON(w, r, &req); err != nil {
			return errors.MalformedDataError(err)
		}

		token, err := s.AuthService.SignUp(r.Context(), req)
		if err != nil {
			return errors.ToAPIError(err)
		}

		c := http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		}

		cookie.Write(w, c)

		err = response.WriteJSON(w, http.StatusCreated, nil)
		if err != nil {
			return err
		}
		return nil
	}
}

func (s *APIServer) handleSignIn() APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var req *rf.SignInRequest

		if err := request.ReadJSON(w, r, &req); err != nil {
			return errors.MalformedDataError(err)
		}

		token, err := s.AuthService.SignIn(r.Context(), req)
		if err != nil {
			return errors.ToAPIError(err)
		}

		c := http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		}

		cookie.Write(w, c)

		err = response.WriteJSON(w, http.StatusOK, nil)
		if err != nil {
			return err
		}
		return nil
	}
}
