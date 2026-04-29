//go:build api

package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	httphandler "github.com/enoro/eiger-webapp/backend/internal/handler/http"
	authuc "github.com/enoro/eiger-webapp/backend/internal/usecase/auth"
	circleuc "github.com/enoro/eiger-webapp/backend/internal/usecase/circle"
	eventuc "github.com/enoro/eiger-webapp/backend/internal/usecase/event"
	userusecase "github.com/enoro/eiger-webapp/backend/internal/usecase/user"
	jwtpkg "github.com/enoro/eiger-webapp/backend/pkg/jwt"
)

const testJWTSecret = "test-secret"

// =============================================================================
// JWT / リクエストヘルパー
// =============================================================================

func makeToken(t *testing.T, userID string, isSystemAdmin bool) string {
	t.Helper()
	gen := jwtpkg.NewGenerator(testJWTSecret, time.Hour)
	token, err := gen.Generate(userID, isSystemAdmin)
	require.NoError(t, err)
	return token
}

func doRequest(t *testing.T, method, url, token string, body any) *http.Response {
	t.Helper()
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		require.NoError(t, err)
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	require.NoError(t, err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func decodeJSON(t *testing.T, resp *http.Response, dst any) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}

// =============================================================================
// テストサーバービルダー
// =============================================================================

func newTestServer(t *testing.T, cfg apiTestConfig) *httptest.Server {
	t.Helper()
	log := zerolog.Nop()
	jwtGen := jwtpkg.NewGenerator(testJWTSecret, time.Hour)

	// auth
	var lineAuthUC *authuc.LineAuthUseCase
	if cfg.userRepo != nil && cfg.lineVerifier != nil && cfg.jwtGenUC != nil {
		lineAuthUC = authuc.NewLineAuthUseCase(cfg.userRepo, cfg.lineVerifier, cfg.jwtGenUC, log)
	}
	var registerUC *authuc.RegisterUseCase
	if cfg.userRepo != nil && cfg.userProfileRepo != nil {
		registerUC = authuc.NewRegisterUseCase(cfg.userRepo, cfg.userProfileRepo, log)
	}
	authHandler := httphandler.NewAuthHandler(lineAuthUC, registerUC, log)

	// user
	var getUserUC *userusecase.GetUserUseCase
	if cfg.userRepo != nil {
		getUserUC = userusecase.NewGetUserUseCase(cfg.userRepo)
	}
	userHandler := httphandler.NewUserHandler(nil, getUserUC, log)

	// circle use cases（admin/circle 両ハンドラで共有）
	var listCirclesUC *circleuc.ListCirclesUseCase
	if cfg.circleRepo != nil && cfg.circleMemberRepo != nil {
		listCirclesUC = circleuc.NewListCirclesUseCase(cfg.circleRepo, cfg.circleMemberRepo)
	}
	var getCircleMembersUC *circleuc.GetCircleMembersUseCase
	if cfg.circleMemberRepo != nil {
		getCircleMembersUC = circleuc.NewGetCircleMembersUseCase(cfg.circleMemberRepo)
	}
	var updateMemberRoleUC *authuc.UpdateCircleMemberRoleUseCase
	if cfg.circleMemberRepo != nil {
		updateMemberRoleUC = authuc.NewUpdateCircleMemberRoleUseCase(cfg.circleMemberRepo)
	}

	// admin handler
	var createCircleUC *circleuc.CreateCircleUseCase
	if cfg.circleRepo != nil {
		createCircleUC = circleuc.NewCreateCircleUseCase(cfg.circleRepo)
	}
	var deleteCircleUC *circleuc.DeleteCircleUseCase
	if cfg.circleRepo != nil {
		deleteCircleUC = circleuc.NewDeleteCircleUseCase(cfg.circleRepo)
	}
	adminHandler := httphandler.NewAdminHandler(
		createCircleUC, listCirclesUC, deleteCircleUC,
		getCircleMembersUC, updateMemberRoleUC, cfg.circleRepo, log,
	)

	// circle handler
	var applyCircleUC *circleuc.ApplyCircleUseCase
	if cfg.circleRepo != nil && cfg.circleMemberRepo != nil {
		applyCircleUC = circleuc.NewApplyCircleUseCase(cfg.circleRepo, cfg.circleMemberRepo)
	}
	circleHandler := httphandler.NewCircleHandler(
		applyCircleUC, updateMemberRoleUC, listCirclesUC, getCircleMembersUC, log,
	)

	// event handler
	var createEventUC *eventuc.CreateEventUseCase
	if cfg.eventRepo != nil {
		createEventUC = eventuc.NewCreateEventUseCase(cfg.eventRepo)
	}
	var listEventsUC *eventuc.ListEventsUseCase
	if cfg.eventRepo != nil {
		listEventsUC = eventuc.NewListEventsUseCase(cfg.eventRepo)
	}
	var applyEventUC *eventuc.ApplyEventUseCase
	if cfg.eventParticipantRepo != nil {
		applyEventUC = eventuc.NewApplyEventUseCase(cfg.eventParticipantRepo)
	}
	var updateStatusUC *eventuc.UpdateParticipantStatusUseCase
	if cfg.eventParticipantRepo != nil {
		updateStatusUC = eventuc.NewUpdateParticipantStatusUseCase(cfg.eventParticipantRepo)
	}
	var bulkUpdateUC *eventuc.BulkUpdateParticipantsUseCase
	if cfg.eventParticipantRepo != nil {
		bulkUpdateUC = eventuc.NewBulkUpdateParticipantsUseCase(cfg.eventParticipantRepo)
	}
	var confirmEventUC *eventuc.ConfirmEventUseCase
	if cfg.eventRepo != nil && cfg.eventParticipantRepo != nil {
		confirmEventUC = eventuc.NewConfirmEventUseCase(cfg.eventRepo, cfg.eventParticipantRepo)
	}
	eventHandler := httphandler.NewEventHandler(
		createEventUC, listEventsUC, applyEventUC,
		updateStatusUC, bulkUpdateUC, confirmEventUC,
		cfg.eventParticipantRepo, log,
	)

	router := httphandler.NewRouter(
		userHandler, authHandler, adminHandler, circleHandler, eventHandler,
		cfg.circleMemberRepo, jwtGen, log,
	)
	return httptest.NewServer(router)
}
