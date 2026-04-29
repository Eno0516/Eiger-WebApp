package event_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	. "github.com/enoro/eiger-webapp/backend/internal/usecase/event"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

// mockEventParticipantRepo は EventParticipantRepository の最小モック実装。
type mockEventParticipantRepo struct{}

func (m *mockEventParticipantRepo) Upsert(_ context.Context, _ *entity.EventParticipant) error {
	return nil
}
func (m *mockEventParticipantRepo) BulkUpsert(_ context.Context, _ []*entity.EventParticipant) error {
	return nil
}
func (m *mockEventParticipantRepo) FindByEvent(_ context.Context, _ string) ([]*entity.EventParticipant, error) {
	return nil, nil
}
func (m *mockEventParticipantRepo) FindByEventAndUser(_ context.Context, _, _ string) (*entity.EventParticipant, error) {
	return nil, apperrors.ErrNotFound
}
func (m *mockEventParticipantRepo) FindConfirmedByEvent(_ context.Context, _ string) ([]*entity.EventParticipant, error) {
	return nil, nil
}

func TestConfirmEventUseCase_Execute(t *testing.T) {
	t.Run("success - イベントが確定される", func(t *testing.T) {
		// Arrange
		eventDate := time.Now().Add(24 * time.Hour)
		eventRepo := &mockEventRepo{
			findByIDFn: func(_ context.Context, id string) (*entity.Event, error) {
				return &entity.Event{
					ID:       id,
					CircleID: "circle-1",
					Title:    "テストイベント",
					Date:     eventDate,
				}, nil
			},
		}
		participantRepo := &mockEventParticipantRepo{}
		uc := NewConfirmEventUseCase(eventRepo, participantRepo)

		// Act
		err := uc.Execute(context.Background(), ConfirmEventInput{EventID: "event-1"})

		// Assert
		require.NoError(t, err)
	})

	t.Run("already confirmed - ErrConflict が返る", func(t *testing.T) {
		// Arrange
		confirmedAt := time.Now()
		eventDate := time.Now().Add(24 * time.Hour)
		eventRepo := &mockEventRepo{
			findByIDFn: func(_ context.Context, id string) (*entity.Event, error) {
				return &entity.Event{
					ID:          id,
					CircleID:    "circle-1",
					Title:       "テストイベント",
					Date:        eventDate,
					ConfirmedAt: &confirmedAt,
				}, nil
			},
		}
		participantRepo := &mockEventParticipantRepo{}
		uc := NewConfirmEventUseCase(eventRepo, participantRepo)

		// Act
		err := uc.Execute(context.Background(), ConfirmEventInput{EventID: "event-1"})

		// Assert
		require.Error(t, err)
		assert.Equal(t, 409, apperrors.HTTPStatusCode(err))
	})
}
