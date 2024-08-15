package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() { <-c; cancel() }()

	m := &Main{
		APIServer: http.NewPostgresAPIServer(),
	}

	if err := m.Run(ctx); err != nil {
		m.Close()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	<-ctx.Done()

	if err := m.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Main struct {
	APIServer *http.APIServer
}

func (m *Main) Run(ctx context.Context) error {
	if err := m.APIServer.Open(); err != nil {
		return err
	}

	log.Printf("running: url=%q dsn=%q", m.APIServer.URL(), rf.Config.DatabaseURL)

	return nil
}

func (m *Main) Close() error {
	if err := m.APIServer.Close(); err != nil {
		return err
	}

	return nil
}
