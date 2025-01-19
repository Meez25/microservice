package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/meez25/microservice/handlers"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	ph := handlers.NewProducts(logger)

	sm := http.NewServeMux()
	sm.Handle("/", ph)

	s := http.Server{
		Addr:         ":3000",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error("Could not create the webserver", "error", err, "port", 3000)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	logger.Info("Received signal to shutdown", "signal", sig)
	tc, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))

	defer cancel()

	s.Shutdown(tc)
	if err := s.Shutdown(tc); err != nil {
		logger.Error("Error during shutdown", "error", err)
	}
}
