package server

import (
	"errors"
	"net/http"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internals/data"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internals/util"
)

func (s *Server) routes() {
	router := http.NewServeMux()

	router.Handle("GET /healthz", s.handleHealthCheck())
	router.Handle("POST /users", s.handleUserCreate())

	s.Handler = router
}

// HealthCheckResponseMsg constant
const HealthCheckResponseMsg = "Health check ok!"

type inMemoryUserDB map[string]data.User

var userDB = make(inMemoryUserDB)

func (s *Server) handleHealthCheck() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.respondWithText(w, r, HealthCheckResponseMsg, http.StatusOK)
	})
}

func (s *Server) handleUserCreate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestData, err := data.GetCreateUserRequestFromBody(w, r)
		if err != nil {
			var mr util.MalformedRequest
			if errors.As(err, &mr) {
				http.Error(w, mr.Msg, mr.Status)
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}

		id := len(userDB) + 1
		userData := data.GetUserFromCreateUserRequestWithID(requestData, id)

		_, ok := userDB[userData.Email]
		if ok {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		userDB[userData.Email] = userData

		responseData := data.GetCreateUserResponseFromUser(userData)
		s.respondWithJSON(w, r, responseData, http.StatusCreated)
	})
}
