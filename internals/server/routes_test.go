package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internals/server"
	"github.com/google/go-cmp/cmp"
)

func TestServer(t *testing.T) {
	t.Parallel()

	svr := server.NewServer()

	t.Run("GET /healthz returns 200", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		wantText := server.HealthCheckResponseMsg
		gotText := response.Body.String()

		assertResponseStatusCode(t, response.Code, http.StatusOK)
		assertResponseText(t, gotText, wantText)
	})
}

func assertResponseStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if !cmp.Equal(got, want) {
		t.Errorf("Did not get correct status code, got %d, want %d", got, want)
	}
}

func assertResponseText(t testing.TB, got, want string) {
	t.Helper()

	if !cmp.Equal(got, want) {
		t.Errorf("Did not get correct status code, got %s, want %s", got, want)
	}
}
