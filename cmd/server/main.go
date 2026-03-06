package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/thrgamon/retro-triage/internal/analyser"
	"github.com/thrgamon/retro-triage/internal/api"
	"github.com/thrgamon/retro-triage/internal/config"
	"github.com/thrgamon/retro-triage/internal/db"
	"github.com/thrgamon/retro-triage/internal/server"
)

// @title Retro Triage API
// @version 1.0
// @host localhost:8080
// @BasePath /api

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.LoadConfig()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer pool.Close()

	pingCtx, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	if err := pool.Ping(pingCtx); err != nil {
		cancelPing()
		log.Fatalf("ping database: %v", err)
	}
	cancelPing()

	queries := db.New(pool)
	openaiClient := analyser.NewOpenAIClient(cfg.OpenAIKey)
	anal := analyser.New(openaiClient)

	handler := api.NewHandler(api.HandlerConfig{
		Queries:  queries,
		Analyser: anal,
		Cfg:      cfg,
	})

	srv := server.New(server.Options{
		Config:  cfg,
		Handler: handler,
	})

	addr := fmt.Sprintf(":%d", cfg.Port)

	go func() {
		if err := srv.Run(addr); err != nil && !errors.Is(err, server.ErrServerClosed) {
			log.Fatalf("server stopped: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}
}
