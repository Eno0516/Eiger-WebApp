package circle_test

import (
	"context"
	"testing"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	. "github.com/enoro/eiger-webapp/backend/internal/usecase/circle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCirclesUseCase_Execute(t *testing.T) {
	t.Run("system admin - 全サークルが返る", func(t *testing.T) {
		// Arrange
		circles := []*entity.Circle{
			{ID: "c1", Name: "サークルA"},
			{ID: "c2", Name: "サークルB"},
		}
		circleRepo := &mockCircleRepo{
			findAllFn: func(_ context.Context) ([]*entity.Circle, error) {
				return circles, nil
			},
		}
		uc := NewListCirclesUseCase(circleRepo, &mockCircleMemberRepo{})

		// Act
		result, err := uc.Execute(context.Background(), ListCirclesInput{IsSystemAdmin: true})

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("circle admin - 自分が管理するサークルのみ返る", func(t *testing.T) {
		// Arrange
		circleA := &entity.Circle{ID: "c1", Name: "サークルA"}
		circleRepo := &mockCircleRepo{
			findByIDFn: func(_ context.Context, id string) (*entity.Circle, error) {
				if id == "c1" {
					return circleA, nil
				}
				return nil, nil
			},
		}
		memberRepo := &mockCircleMemberRepo{
			findByUserFn: func(_ context.Context, _ string) ([]*entity.CircleMember, error) {
				return []*entity.CircleMember{
					{CircleID: "c1", UserID: "u1", Role: entity.CircleRoleAdmin},
				}, nil
			},
		}
		uc := NewListCirclesUseCase(circleRepo, memberRepo)

		// Act
		result, err := uc.Execute(context.Background(), ListCirclesInput{IsSystemAdmin: false, RequesterUserID: "u1"})

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "c1", result[0].ID)
	})
}
