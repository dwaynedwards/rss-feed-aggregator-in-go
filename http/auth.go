package http

import (
	"net/http"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/builder"
)

func (s *APIServer) registerAuthRoutes(r *http.ServeMux) {
	r.Handle("POST /api/v1/auths/signup", makeHTTPHandlerFunc(s.handleAuthSignUp()))
	r.Handle("POST /api/v1/auths/signin", makeHTTPHandlerFunc(s.handleAuthSignIn()))
}

func (s *APIServer) handleAuthSignUp() APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		auth, err := getSignUpAuth(w, r)
		if err != nil {
			return err
		}

		token, err := s.AuthService.SignUp(r.Context(), auth)
		if err != nil {
			return err
		}

		cookie := http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		}

		rf.Write(w, cookie)

		err = writeJSON(w, http.StatusCreated, nil)
		if err != nil {
			return err
		}
		return nil
	}
}

func (s *APIServer) handleAuthSignIn() APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		auth, err := getSignInAuth(w, r)
		if err != nil {
			return err
		}

		token, err := s.AuthService.SignIn(r.Context(), auth)
		if err != nil {
			return err
		}

		cookie := http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		}

		rf.Write(w, cookie)

		err = writeJSON(w, http.StatusOK, nil)
		if err != nil {
			return err
		}
		return nil
	}
}

func getSignUpAuth(w http.ResponseWriter, r *http.Request) (*rf.Auth, error) {
	var req *rf.SignUpAuthRequest

	if err := decodeJSON(w, r, &req); err != nil {
		return nil, rf.NewAPIError(http.StatusUnprocessableEntity, err)
	}

	if err := validateSignUpRequest(req); err != nil {
		return nil, err
	}

	return builder.NewAuthBuilder().
		WithUser(builder.NewUserBuilder().
			WithName(req.Name)).
		WithBasicAuth(builder.NewBasicAuthBuilder().
			WithEmail(req.Email).
			WithPassword(req.Password)).
		Build(), nil
}

func validateSignUpRequest(req *rf.SignUpAuthRequest) error {
	errs := map[string]string{}

	if req.Email == "" {
		errs["email"] = rf.EMEmailRequired
	}

	if req.Password == "" {
		errs["password"] = rf.EMPasswordRequired
	}

	if req.Name == "" {
		errs["name"] = rf.EMNameRequired
	}

	if len(errs) > 0 {
		return rf.NewAPIError(http.StatusBadRequest, errs)
	}

	return nil
}

func getSignInAuth(w http.ResponseWriter, r *http.Request) (*rf.Auth, error) {
	var req *rf.SignInAuthRequest

	if err := decodeJSON(w, r, &req); err != nil {
		return nil, rf.NewAPIError(http.StatusUnprocessableEntity, err)
	}

	if err := validateSigninRequest(req); err != nil {
		return nil, err
	}

	return builder.NewAuthBuilder().
		WithBasicAuth(builder.NewBasicAuthBuilder().
			WithEmail(req.Email).
			WithPassword(req.Password)).
		Build(), nil
}

func validateSigninRequest(req *rf.SignInAuthRequest) error {
	errs := map[string]string{}

	if req.Email == "" {
		errs["email"] = rf.EMEmailRequired
	}

	if req.Password == "" {
		errs["password"] = rf.EMPasswordRequired
	}

	if len(errs) > 0 {
		return rf.NewAPIError(http.StatusBadRequest, errs)
	}

	return nil
}
