package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internals/account"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internals/common"
)

func (s *server) routes() {
	router := http.NewServeMux()

	router.Handle("GET /healthz", makeHttpHandlerFunc(s.handleHealthCheck()))
	router.Handle("POST /accounts", makeHttpHandlerFunc(s.handleAccountCreate()))

	s.Handler = router
}

const HealthCheckResponseMsg = "Health check ok!"

func makeHttpHandlerFunc(a apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := a(w, r); err != nil {
			var mr *common.MalformedRequestError
			var a *account.AccountError

			if errors.As(err, &mr) {
				http.Error(w, mr.Error(), mr.Status)
			} else if errors.As(err, &a) {
				http.Error(w, a.Error(), a.Status)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func (s *server) handleHealthCheck() apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		respondWithText(w, http.StatusOK, HealthCheckResponseMsg)
		return nil
	}
}

func (s *server) handleAccountCreate() apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		defer io.Copy(io.Discard, r.Body)
		defer r.Body.Close()

		req, err := account.GetCreateAccountRequestFromBody(w, r)
		if err != nil {
			return err
		}

		res, err := s.accountService.CreateAccount(req)
		if err != nil {
			return err
		}

		return respondWithJSON(w, http.StatusCreated, res)
	}
}

func respondWithText(w http.ResponseWriter, status int, data string) {
	w.Header().Set("Content-Type", ContentTypePlainText.Value)
	w.WriteHeader(status)
	w.Write([]byte(data))
}

func respondWithJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", ContentTypeJSON.Value)
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
