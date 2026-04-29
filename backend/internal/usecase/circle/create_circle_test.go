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

// mockCircleRepo は CircleRepository の最小モック実装。
type mockCircleRepo struct {
	createErr  error
	findByIDFn func(ctx context.Context, id string) (*entity.Circle, error)
	findAllFn  func(ctx context.Context) ([]*entity.Circle, error)
	deleteErr  error
}

func (m *mockCircleRepo) Create(_ context.Context, _ *entity.Circle) error { return m.createErr }
func (m *mockCircleRepo) FindByID(ctx context.Context, id string) (*entity.Circle, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *mockCircleRepo) FindAll(ctx context.Context) ([]*entity.Circle, error) {
	if m.findAllFn != nil {
		return m.findAllFn(ctx)
	}
	return nil, nil
}
func (m *mockCircleRepo) Update(_ context.Context, _ *entity.Circle) error { return nil }
func (m *mockCircleRepo) Delete(_ context.Context, _ string) error         { return m.deleteErr }

func TestCreateCircleUseCase_Execute(t *testing.T) {
	t.Run("success - サークルが作成される", func(t *testing.T) {
		// Arrange
		repo := &mockCircleRepo{createErr: nil}
		uc := NewCreateCircleUseCase(repo)

		// Act
		result, err := uc.Execute(context.Background(), CreateCircleInput{Name: "テストサークル"})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "テストサークル", result.Name)
	})

	t.Run("empty name - ErrBadRequest が返る", func(t *testing.T) {
		// Arrange
		repo := &mockCircleRepo{}
		uc := NewCreateCircleUseCase(repo)

		// Act
		_, err := uc.Execute(context.Background(), CreateCircleInput{Name: ""})

		// Assert
		require.Error(t, err)
		assert.Equal(t, 400, apperrors.HTTPStatusCode(err))
	})
}
