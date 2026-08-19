// Command server is the composition root of the calculator backend: it reads
// configuration, wires the domain service into the HTTP router, and runs the
// server with sensible timeouts and graceful shutdown. All construction lives
// here so the packages it depends on stay free of global state.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/portella94/sezzle-challange-calculator-go-react/backend/internal/api"
	"github.com/portella94/sezzle-challange-calculator-go-react/backend/internal/calculator"
	"github.com/portella94/sezzle-challange-calculator-go-react/backend/internal/config"
	"github.com/portella94/sezzle-challange-calculator-go-react/backend/internal/version"
)

func main() {
	cfg := config.Load()

	handler := api.NewHandler(calculator.New())
	router := api.NewRouter(handler, cfg.AllowedOrigins)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Run the listener in the background so main can wait for a shutdown signal.
	errCh := make(chan error, 1)
	go func() {
		log.Printf("calculator backend v%s listening on %s", version.Version, srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		log.Fatalf("server error: %v", err)
	case <-ctx.Done():
		log.Println("shutdown signal received, draining connections...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("graceful shutdown failed: %v", err)
		}
		log.Println("server stopped cleanly")
	}
	os.Exit(0)
}
