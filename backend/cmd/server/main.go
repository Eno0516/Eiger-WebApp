package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/enoro/eiger-webapp/backend/internal/infrastructure/config"
	"github.com/enoro/eiger-webapp/backend/internal/infrastructure/database"
	"github.com/enoro/eiger-webapp/backend/internal/infrastructure/logger"

	httphandler "github.com/enoro/eiger-webapp/backend/internal/handler/http"
	"github.com/enoro/eiger-webapp/backend/internal/repository/postgres"
	"github.com/enoro/eiger-webapp/backend/internal/usecase/user"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// ── Config ──────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// ── Logger ──────────────────────────────────────────────────────────────
	log := logger.New(cfg.LogLevel, cfg.Env)

	// ── Database ─────────────────────────────────────────────────────────────
	ctx := context.Background()
	db, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer db.Close()
	log.Info().Msg("connected to database")

	// ── Repositories ─────────────────────────────────────────────────────────
	userRepo := postgres.NewUserRepository(db)

	// ── Use Cases ────────────────────────────────────────────────────────────
	createUserUC := user.NewCreateUserUseCase(userRepo)
	getUserUC := user.NewGetUserUseCase(userRepo)

	// ── Handlers ─────────────────────────────────────────────────────────────
	userHandler := httphandler.NewUserHandler(createUserUC, getUserUC, log)

	// ── Router ───────────────────────────────────────────────────────────────
	router := httphandler.NewRouter(userHandler, log)

	// ── Server ───────────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info().Str("port", cfg.Port).Msg("starting server")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	<-quit
	log.Info().Msg("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	log.Info().Msg("server stopped")
	return nil
}
