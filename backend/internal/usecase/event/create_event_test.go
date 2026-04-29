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

// mockEventRepo は EventRepository の最小モック実装。
type mockEventRepo struct {
	createErr  error
	findByIDFn func(ctx context.Context, id string) (*entity.Event, error)
	updateErr  error
}

func (m *mockEventRepo) Create(_ context.Context, _ *entity.Event) error { return m.createErr }
func (m *mockEventRepo) FindByID(ctx context.Context, id string) (*entity.Event, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, apperrors.ErrNotFound
}
func (m *mockEventRepo) FindByCircle(_ context.Context, _ string) ([]*entity.Event, error) {
	return nil, nil
}
func (m *mockEventRepo) Update(_ context.Context, _ *entity.Event) error { return m.updateErr }
func (m *mockEventRepo) Delete(_ context.Context, _ string) error        { return nil }

func TestCreateEventUseCase_Execute(t *testing.T) {
	t.Run("success - イベントが作成される", func(t *testing.T) {
		// Arrange
		repo := &mockEventRepo{createErr: nil}
		uc := NewCreateEventUseCase(repo)

		// Act
		result, err := uc.Execute(context.Background(), CreateEventInput{
			CircleID: "circle-1",
			Title:    "テストイベント",
			Date:     time.Now().Add(24 * time.Hour),
		})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "テストイベント", result.Title)
		assert.Equal(t, "circle-1", result.CircleID)
	})

	t.Run("empty title - ErrBadRequest が返る", func(t *testing.T) {
		// Arrange
		repo := &mockEventRepo{}
		uc := NewCreateEventUseCase(repo)

		// Act
		_, err := uc.Execute(context.Background(), CreateEventInput{
			CircleID: "circle-1",
			Title:    "",
			Date:     time.Now().Add(24 * time.Hour),
		})

		// Assert
		require.Error(t, err)
		assert.Equal(t, 400, apperrors.HTTPStatusCode(err))
	})

	t.Run("zero date - ErrBadRequest が返る", func(t *testing.T) {
		// Arrange
		repo := &mockEventRepo{}
		uc := NewCreateEventUseCase(repo)

		// Act
		_, err := uc.Execute(context.Background(), CreateEventInput{
			CircleID: "circle-1",
			Title:    "テストイベント",
			Date:     time.Time{},
		})

		// Assert
		require.Error(t, err)
		assert.Equal(t, 400, apperrors.HTTPStatusCode(err))
	})
}
