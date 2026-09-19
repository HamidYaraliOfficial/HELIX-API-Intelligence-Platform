// Command server boots the HELIX backend: it loads configuration, connects
// to PostgreSQL, applies pending migrations, and serves the REST API that
// powers the API Catalog, Discovery Engine and Control Center frontend.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/helix-platform/helix-backend/internal/config"
	"github.com/helix-platform/helix-backend/internal/db"
	"github.com/helix-platform/helix-backend/internal/handlers"
	"github.com/helix-platform/helix-backend/internal/repository"
)

func main() {
	cfg := config.Load()
	log.Printf("[helix] starting backend (env=%s port=%s)", cfg.Environment, cfg.Port)

	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[helix] database connection failed: %v", err)
	}
	defer conn.Close()

	if err := db.Migrate(conn); err != nil {
		log.Fatalf("[helix] migration failed: %v", err)
	}

	if err := repository.EnsureDefaultOrganization(conn, cfg.DefaultOrgID, "Default Organization", "default"); err != nil {
		log.Fatalf("[helix] failed to ensure default organization: %v", err)
	}

	router := handlers.NewRouter(conn, cfg)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("[helix] listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[helix] server error: %v", err)
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM (Kubernetes-friendly).
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("[helix] shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[helix] graceful shutdown failed: %v", err)
	}
	log.Println("[helix] stopped")
}
