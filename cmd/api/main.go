package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/account"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internals/server"
	"github.com/google/go-cmp/cmp"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	portStr := getEnvVar("PORT")

	accountStore := account.NewStore()
	accountService := account.NewService(accountStore)
	accountServer := account.NewServer(accountService)
	server := server.NewServer(accountServer)

	log.Printf("Starting listening at http://localhost:%s\n", portStr)
	if err := http.ListenAndServe(":"+portStr, server); err != nil {
		log.Fatalf("Could not listen at http://localhost:%s - %v", portStr, err)
	}
}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Failed to load env file: %s\n", err.Error())
	}
}

func getEnvVar(key string) string {
	variable := os.Getenv(key)
	if cmp.Equal(variable, "") {
		log.Fatalf("%s is not set in as an environment variable\n", key)
	}
	return variable
}
