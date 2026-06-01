package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	err := run(ctx)

	if err != nil && !errors.Is(err, context.Canceled) {
		log.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Get database url
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}

	// Connect to and verify database
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	// Create table if nonexistent
	_, err = db.ExecContext(ctx, `
      CREATE TABLE IF NOT EXISTS events (
              id         SERIAL PRIMARY KEY,
              path       TEXT NOT NULL,
              referrer   TEXT,
              user_agent TEXT,
              ip         TEXT,
              created_at TIMESTAMPTZ DEFAULT NOW()
      )
  `)
	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	// Register handlers
	mux := http.NewServeMux()
	mux.HandleFunc("POST /events", handleEvent(db))
	mux.HandleFunc("GET /stats", handleStats(db))
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Build server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: corsMiddleware(mux),
	}

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server: %v", err)
		}
	}()
	log.Println("listening on :8080")

	// Block until shutdown
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
