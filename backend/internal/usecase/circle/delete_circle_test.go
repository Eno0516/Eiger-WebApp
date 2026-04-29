package circle_test

import (
	"context"
	"testing"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	. "github.com/enoro/eiger-webapp/backend/internal/usecase/circle"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteCircleUseCase_Execute(t *testing.T) {
	t.Run("success - サークルが削除される", func(t *testing.T) {
		// Arrange
		repo := &mockCircleRepo{
			findByIDFn: func(_ context.Context, id string) (*entity.Circle, error) {
				return &entity.Circle{ID: id, Name: "サークルA"}, nil
			},
			deleteErr: nil,
		}
		uc := NewDeleteCircleUseCase(repo)

		// Act
		err := uc.Execute(context.Background(), DeleteCircleInput{CircleID: "circle-1"})

		// Assert
		require.NoError(t, err)
	})

	t.Run("not found - ErrNotFound が返る", func(t *testing.T) {
		// Arrange
		repo := &mockCircleRepo{
			findByIDFn: func(_ context.Context, _ string) (*entity.Circle, error) {
				return nil, apperrors.ErrNotFound
			},
		}
		uc := NewDeleteCircleUseCase(repo)

		// Act
		err := uc.Execute(context.Background(), DeleteCircleInput{CircleID: "no-such-id"})

		// Assert
		require.Error(t, err)
		assert.Equal(t, 404, apperrors.HTTPStatusCode(err))
	})
}
