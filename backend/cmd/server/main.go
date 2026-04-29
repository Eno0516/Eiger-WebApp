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
	memrepo "github.com/enoro/eiger-webapp/backend/internal/repository/memory"
	"github.com/enoro/eiger-webapp/backend/internal/repository/postgres"
	authuc "github.com/enoro/eiger-webapp/backend/internal/usecase/auth"
	circleuc "github.com/enoro/eiger-webapp/backend/internal/usecase/circle"
	eventuc "github.com/enoro/eiger-webapp/backend/internal/usecase/event"
	useruс "github.com/enoro/eiger-webapp/backend/internal/usecase/user"
	jwtpkg "github.com/enoro/eiger-webapp/backend/pkg/jwt"
	"github.com/enoro/eiger-webapp/backend/pkg/lineauth"
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
	userProfileRepo := postgres.NewUserProfileRepository(db)

	// Circle / Event リポジトリはまだ postgres 実装がないため memory を使用
	// TODO: postgres 実装が揃い次第差し替える
	circleRepo := memrepo.NewCircleRepository()
	circleMemberRepo := memrepo.NewCircleMemberRepository()
	eventRepo := memrepo.NewEventRepository()
	eventParticipantRepo := memrepo.NewEventParticipantRepository()

	// ── Infrastructure ───────────────────────────────────────────────────────
	// TODO: cfg に LINE_CHANNEL_ID / JWT_SECRET_KEY / GOOGLE_SERVICE_ACCOUNT / SIGNUP_SPREADSHEET_ID を追加する
	lineVerifier := lineauth.NewVerifier(cfg.LineChannelID)
	jwtGen := jwtpkg.NewGenerator(cfg.JWTSecretKey, 0)
	// TODO: Copilot - sheets.NewSignUpWriter(credentialsJSON, cfg.SignupSpreadsheetID, "Sheet1") で初期化する
	// TODO: Copilot - credentialsJSON は cfg.GoogleServiceAccount を base64 デコードして使う
	// TODO: Copilot - エラー時は return fmt.Errorf("init sheets writer: %w", err)

	// ── Use Cases ────────────────────────────────────────────────────────────
	lineAuthUC := authuc.NewLineAuthUseCase(userRepo, lineVerifier, jwtGen, log)
	registerUC := authuc.NewRegisterUseCase(userRepo, userProfileRepo, log)
	createUserUC := useruс.NewCreateUserUseCase(userRepo, userProfileRepo)
	getUserUC := useruс.NewGetUserUseCase(userRepo)

	// Circle use cases
	createCircleUC := circleuc.NewCreateCircleUseCase(circleRepo)
	listCirclesUC := circleuc.NewListCirclesUseCase(circleRepo, circleMemberRepo)
	deleteCircleUC := circleuc.NewDeleteCircleUseCase(circleRepo)
	getCircleMembersUC := circleuc.NewGetCircleMembersUseCase(circleMemberRepo)
	applyCircleUC := circleuc.NewApplyCircleUseCase(circleRepo, circleMemberRepo)
	updateMemberRoleUC := authuc.NewUpdateCircleMemberRoleUseCase(circleMemberRepo)

	// Event use cases
	createEventUC := eventuc.NewCreateEventUseCase(eventRepo)
	listEventsUC := eventuc.NewListEventsUseCase(eventRepo)
	applyEventUC := eventuc.NewApplyEventUseCase(eventParticipantRepo)
	updateParticipantStatusUC := eventuc.NewUpdateParticipantStatusUseCase(eventParticipantRepo)
	bulkUpdateParticipantsUC := eventuc.NewBulkUpdateParticipantsUseCase(eventParticipantRepo)
	confirmEventUC := eventuc.NewConfirmEventUseCase(eventRepo, eventParticipantRepo)

	// ── Handlers ─────────────────────────────────────────────────────────────
	authHandler := httphandler.NewAuthHandler(lineAuthUC, registerUC, log)
	userHandler := httphandler.NewUserHandler(createUserUC, getUserUC, log)
	adminHandler := httphandler.NewAdminHandler(
		createCircleUC, listCirclesUC, deleteCircleUC,
		getCircleMembersUC, updateMemberRoleUC, circleRepo, log,
	)
	circleHandler := httphandler.NewCircleHandler(
		applyCircleUC, updateMemberRoleUC,
		listCirclesUC, getCircleMembersUC, log,
	)
	eventHandler := httphandler.NewEventHandler(
		createEventUC, listEventsUC, applyEventUC,
		updateParticipantStatusUC, bulkUpdateParticipantsUC, confirmEventUC,
		eventParticipantRepo, log,
	)

	// ── Router ───────────────────────────────────────────────────────────────
	router := httphandler.NewRouter(
		userHandler, authHandler, adminHandler, circleHandler, eventHandler,
		circleMemberRepo, jwtGen, log,
	)

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
