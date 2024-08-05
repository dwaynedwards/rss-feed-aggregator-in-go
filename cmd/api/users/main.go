package main

import (
	"log"
	"net/http"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/common"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/users"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/users/store"
)

func main() {
	portStr := common.GetEnvVar("PORT")

	usersStore, cleanup, err := store.NewPostgresUsersStore()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	usersService := users.NewService(usersStore)
	usersServer := users.NewServer(usersService)
	server := makeNewServer(usersServer)

	log.Printf("Starting listening at http://localhost:%s\n", portStr)
	if err := http.ListenAndServe(":"+portStr, server); err != nil {
		log.Fatalf("Could not listen at http://localhost:%s - %v", portStr, err)
	}
}

type server struct {
	http.Handler
}

func makeNewServer(usersServer users.UsersServer) *server {
	s := new(server)

	router := http.NewServeMux()
	usersServer.RegisterEndpoints(router)
	s.Handler = router

	return s
}
