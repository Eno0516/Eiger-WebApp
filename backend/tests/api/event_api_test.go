//go:build api

package api_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/internal/handler/http/dto"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// モック定義
// =============================================================================

type mockEventRepoAPI struct {
	createErr  error
	findByIDFn func(ctx context.Context, id string) (*entity.Event, error)
	findAllFn  func(ctx context.Context, circleID string) ([]*entity.Event, error)
	updateErr  error
}

func (m *mockEventRepoAPI) Create(_ context.Context, _ *entity.Event) error { return m.createErr }
func (m *mockEventRepoAPI) FindByID(ctx context.Context, id string) (*entity.Event, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, apperrors.ErrNotFound
}
func (m *mockEventRepoAPI) FindByCircle(ctx context.Context, circleID string) ([]*entity.Event, error) {
	if m.findAllFn != nil {
		return m.findAllFn(ctx, circleID)
	}
	return nil, nil
}
func (m *mockEventRepoAPI) Update(_ context.Context, _ *entity.Event) error { return m.updateErr }
func (m *mockEventRepoAPI) Delete(_ context.Context, _ string) error        { return nil }

var _ domainrepo.EventRepository = (*mockEventRepoAPI)(nil)

type mockEventParticipantRepoAPI struct {
	upsertErr              error
	findByEventFn          func(ctx context.Context, eventID string) ([]*entity.EventParticipant, error)
	findByEventAndUserFn   func(ctx context.Context, eventID, userID string) (*entity.EventParticipant, error)
	findConfirmedByEventFn func(ctx context.Context, eventID string) ([]*entity.EventParticipant, error)
}

func (m *mockEventParticipantRepoAPI) Upsert(_ context.Context, _ *entity.EventParticipant) error {
	return m.upsertErr
}
func (m *mockEventParticipantRepoAPI) BulkUpsert(_ context.Context, _ []*entity.EventParticipant) error {
	return nil
}
func (m *mockEventParticipantRepoAPI) FindByEvent(ctx context.Context, eventID string) ([]*entity.EventParticipant, error) {
	if m.findByEventFn != nil {
		return m.findByEventFn(ctx, eventID)
	}
	return nil, nil
}
func (m *mockEventParticipantRepoAPI) FindByEventAndUser(ctx context.Context, eventID, userID string) (*entity.EventParticipant, error) {
	if m.findByEventAndUserFn != nil {
		return m.findByEventAndUserFn(ctx, eventID, userID)
	}
	return nil, apperrors.ErrNotFound
}
func (m *mockEventParticipantRepoAPI) FindConfirmedByEvent(ctx context.Context, eventID string) ([]*entity.EventParticipant, error) {
	if m.findConfirmedByEventFn != nil {
		return m.findConfirmedByEventFn(ctx, eventID)
	}
	return nil, nil
}

var _ domainrepo.EventParticipantRepository = (*mockEventParticipantRepoAPI)(nil)

// futureDate は翌日の日付文字列を返す。
func futureDate() string {
	return time.Now().Add(24 * time.Hour).Format("2006-01-02")
}

// =============================================================================
// GET /admin/circles/{id}/events  （JWT + system_admin or circle_admin 必須）
// =============================================================================

func TestListEventsAdminAPI(t *testing.T) {
	t.Run("success - イベント一覧が返る", func(t *testing.T) {
		// Arrange
		eventRepo := &mockEventRepoAPI{
			findAllFn: func(_ context.Context, _ string) ([]*entity.Event, error) {
				return []*entity.Event{
					{ID: "e1", CircleID: "c1", Title: "イベントA", Date: time.Now().Add(24 * time.Hour)},
				}, nil
			},
		}
		srv := newTestServer(t, apiTestConfig{
			circleRepo:           adminCircleRepo(),
			circleMemberRepo:     &mockCircleMemberRepoAPI{},
			eventRepo:            eventRepo,
			eventParticipantRepo: &mockEventParticipantRepoAPI{},
		})
		defer srv.Close()
		token := makeToken(t, "admin-user", true)

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/admin/circles/c1/events", token, nil)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body []dto.EventResponse
		decodeJSON(t, resp, &body)
		assert.Len(t, body, 1)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:           adminCircleRepo(),
			circleMemberRepo:     &mockCircleMemberRepoAPI{},
			eventRepo:            &mockEventRepoAPI{},
			eventParticipantRepo: &mockEventParticipantRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/admin/circles/c1/events", "", nil)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

// =============================================================================
// POST /admin/circles/{id}/events  （JWT + system_admin or circle_admin 必須）
// =============================================================================

func TestCreateEventAdminAPI(t *testing.T) {
	t.Run("success - 201 + 作成されたイベントが返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:           adminCircleRepo(),
			circleMemberRepo:     &mockCircleMemberRepoAPI{},
			eventRepo:            &mockEventRepoAPI{},
			eventParticipantRepo: &mockEventParticipantRepoAPI{},
		})
		defer srv.Close()
		token := makeToken(t, "admin-user", true)

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/admin/circles/c1/events", token,
			dto.CreateEventRequest{Title: "新イベント", Date: futureDate()},
		)

		// Assert
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		var body dto.EventResponse
		decodeJSON(t, resp, &body)
		assert.Equal(t, "新イベント", body.Title)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:           adminCircleRepo(),
			circleMemberRepo:     &mockCircleMemberRepoAPI{},
			eventRepo:            &mockEventRepoAPI{},
			eventParticipantRepo: &mockEventParticipantRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/admin/circles/c1/events", "",
			dto.CreateEventRequest{Title: "新イベント", Date: futureDate()},
		)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("invalid date format - 400 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:           adminCircleRepo(),
			circleMemberRepo:     &mockCircleMemberRepoAPI{},
			eventRepo:            &mockEventRepoAPI{},
			eventParticipantRepo: &mockEventParticipantRepoAPI{},
		})
		defer srv.Close()
		token := makeToken(t, "admin-user", true)

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/admin/circles/c1/events", token,
			dto.CreateEventRequest{Title: "新イベント", Date: "not-a-date"},
		)

		// Assert
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// =============================================================================
// GET /admin/circles/{id}/events/{eventId}/participants  （JWT + system_admin or circle_admin 必須）
// =============================================================================

func TestGetParticipantsAdminAPI(t *testing.T) {
	t.Run("success - 参加者一覧が返る", func(t *testing.T) {
		// Arrange
		participantRepo := &mockEventParticipantRepoAPI{
			findByEventFn: func(_ context.Context, _ string) ([]*entity.EventParticipant, error) {
				return []*entity.EventParticipant{
					{EventID: "e1", UserID: "u1", Status: entity.ParticipantStatusApplied},
				}, nil
			},
		}
		srv := newTestServer(t, apiTestConfig{
			circleRepo:           adminCircleRepo(),
			circleMemberRepo:     &mockCircleMemberRepoAPI{},
			eventRepo:            &mockEventRepoAPI{},
			eventParticipantRepo: participantRepo,
		})
		defer srv.Close()
		token := makeToken(t, "admin-user", true)

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/admin/circles/c1/events/e1/participants", token, nil)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body []dto.EventParticipantResponse
		decodeJSON(t, resp, &body)
		assert.Len(t, body, 1)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:           adminCircleRepo(),
			circleMemberRepo:     &mockCircleMemberRepoAPI{},
			eventRepo:            &mockEventRepoAPI{},
			eventParticipantRepo: &mockEventParticipantRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/admin/circles/c1/events/e1/participants", "", nil)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

// =============================================================================
// PATCH /admin/circles/{id}/events/{eventId}/participants  （JWT + system_admin or circle_admin 必須）
// =============================================================================

func TestBulkUpdateParticipantsAdminAPI(t *testing.T) {
	t.Run("success - 200 + participants updated が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:           adminCircleRepo(),
			circleMemberRepo:     &mockCircleMemberRepoAPI{},
			eventRepo:            &mockEventRepoAPI{},
			eventParticipantRepo: &mockEventParticipantRepoAPI{},
		})
		defer srv.Close()
		token := makeToken(t, "admin-user", true)

		// Act
		resp := doRequest(t, http.MethodPatch, srv.URL+"/api/v1/admin/circles/c1/events/e1/participants", token,
			dto.BulkUpdateParticipantsRequest{
				Participants: []dto.ParticipantStatusItem{
					{UserID: "u1", Status: "confirmed"},
				},
			},
		)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:           adminCircleRepo(),
			circleMemberRepo:     &mockCircleMemberRepoAPI{},
			eventRepo:            &mockEventRepoAPI{},
			eventParticipantRepo: &mockEventParticipantRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodPatch, srv.URL+"/api/v1/admin/circles/c1/events/e1/participants", "", nil)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

// =============================================================================
// PATCH /admin/circles/{id}/events/{eventId}/participants/{userId}  （JWT + system_admin or circle_admin 必須）
// =============================================================================

func TestUpdateParticipantStatusAdminAPI(t *testing.T) {
	t.Run("success - 200 + status updated が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:           adminCircleRepo(),
			circleMemberRepo:     &mockCircleMemberRepoAPI{},
			eventRepo:            &mockEventRepoAPI{},
			eventParticipantRepo: &mockEventParticipantRepoAPI{},
		})
		defer srv.Close()
		token := makeToken(t, "admin-user", true)

		// Act
		resp := doRequest(t, http.MethodPatch, srv.URL+"/api/v1/admin/circles/c1/events/e1/participants/u1", token,
			dto.UpdateParticipantStatusRequest{Status: "confirmed"},
		)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:           adminCircleRepo(),
			circleMemberRepo:     &mockCircleMemberRepoAPI{},
			eventRepo:            &mockEventRepoAPI{},
			eventParticipantRepo: &mockEventParticipantRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodPatch, srv.URL+"/api/v1/admin/circles/c1/events/e1/participants/u1", "", nil)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

// =============================================================================
// POST /admin/circles/{id}/events/{eventId}/confirm  （JWT + system_admin or circle_admin 必須）
// =============================================================================

func TestConfirmEventAdminAPI(t *testing.T) {
	t.Run("success - 200 + event confirmed が返る", func(t *testing.T) {
		// Arrange
		eventRepo := &mockEventRepoAPI{
			findByIDFn: func(_ context.Context, id string) (*entity.Event, error) {
				return &entity.Event{
					ID: id, CircleID: "c1", Title: "イベントA",
					Date: time.Now().Add(24 * time.Hour),
				}, nil
			},
		}
		srv := newTestServer(t, apiTestConfig{
			circleRepo:           adminCircleRepo(),
			circleMemberRepo:     &mockCircleMemberRepoAPI{},
			eventRepo:            eventRepo,
			eventParticipantRepo: &mockEventParticipantRepoAPI{},
		})
		defer srv.Close()
		token := makeToken(t, "admin-user", true)

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/admin/circles/c1/events/e1/confirm", token, nil)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body dto.MessageResponse
		decodeJSON(t, resp, &body)
		assert.Equal(t, "event confirmed", body.Message)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:           adminCircleRepo(),
			circleMemberRepo:     &mockCircleMemberRepoAPI{},
			eventRepo:            &mockEventRepoAPI{},
			eventParticipantRepo: &mockEventParticipantRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/admin/circles/c1/events/e1/confirm", "", nil)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("already confirmed - 409 が返る", func(t *testing.T) {
		// Arrange
		confirmedAt := time.Now()
		eventRepo := &mockEventRepoAPI{
			findByIDFn: func(_ context.Context, id string) (*entity.Event, error) {
				return &entity.Event{
					ID: id, CircleID: "c1", Title: "イベントA",
					Date:        time.Now().Add(24 * time.Hour),
					ConfirmedAt: &confirmedAt,
				}, nil
			},
		}
		srv := newTestServer(t, apiTestConfig{
			circleRepo:           adminCircleRepo(),
			circleMemberRepo:     &mockCircleMemberRepoAPI{},
			eventRepo:            eventRepo,
			eventParticipantRepo: &mockEventParticipantRepoAPI{},
		})
		defer srv.Close()
		token := makeToken(t, "admin-user", true)

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/admin/circles/c1/events/e1/confirm", token, nil)

		// Assert
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})
}

// =============================================================================
// POST /circles/{circleId}/events/{eventId}/apply  （JWT 必須）
// =============================================================================

func TestApplyEventAPI(t *testing.T) {
	t.Run("success - 200 + applied が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:           adminCircleRepo(),
			circleMemberRepo:     &mockCircleMemberRepoAPI{},
			eventRepo:            &mockEventRepoAPI{},
			eventParticipantRepo: &mockEventParticipantRepoAPI{},
		})
		defer srv.Close()
		token := makeToken(t, "user-1", false)

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/circles/c1/events/e1/apply", token, nil)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body dto.MessageResponse
		decodeJSON(t, resp, &body)
		assert.Equal(t, "applied", body.Message)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:           adminCircleRepo(),
			circleMemberRepo:     &mockCircleMemberRepoAPI{},
			eventRepo:            &mockEventRepoAPI{},
			eventParticipantRepo: &mockEventParticipantRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/circles/c1/events/e1/apply", "", nil)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
