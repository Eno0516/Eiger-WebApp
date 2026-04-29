package circle_test

import (
	"context"
	"testing"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	. "github.com/enoro/eiger-webapp/backend/internal/usecase/circle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockCircleMemberRepo は CircleMemberRepository の最小モック実装。
type mockCircleMemberRepo struct {
	findByCircleFn func(ctx context.Context, circleID string) ([]*entity.CircleMember, error)
	findByUserFn   func(ctx context.Context, userID string) ([]*entity.CircleMember, error)
}

func (m *mockCircleMemberRepo) Upsert(_ context.Context, _ *entity.CircleMember) error {
	return nil
}
func (m *mockCircleMemberRepo) FindByCircleAndUser(_ context.Context, _, _ string) (*entity.CircleMember, error) {
	return nil, nil
}
func (m *mockCircleMemberRepo) FindByCircle(ctx context.Context, circleID string) ([]*entity.CircleMember, error) {
	if m.findByCircleFn != nil {
		return m.findByCircleFn(ctx, circleID)
	}
	return nil, nil
}
func (m *mockCircleMemberRepo) FindByUser(ctx context.Context, userID string) ([]*entity.CircleMember, error) {
	if m.findByUserFn != nil {
		return m.findByUserFn(ctx, userID)
	}
	return nil, nil
}
func (m *mockCircleMemberRepo) FindAdminsByCircle(_ context.Context, _ string) ([]*entity.CircleMember, error) {
	return nil, nil
}
func (m *mockCircleMemberRepo) Delete(_ context.Context, _, _ string) error { return nil }

func TestGetCircleMembersUseCase_Execute(t *testing.T) {
	t.Run("success - メンバー一覧が返る", func(t *testing.T) {
		// Arrange
		expectedMembers := []*entity.CircleMember{
			{CircleID: "c1", UserID: "u1", Role: entity.CircleRoleMember},
			{CircleID: "c1", UserID: "u2", Role: entity.CircleRoleMember},
			{CircleID: "c1", UserID: "u3", Role: entity.CircleRoleAdmin},
		}
		repo := &mockCircleMemberRepo{
			findByCircleFn: func(_ context.Context, _ string) ([]*entity.CircleMember, error) {
				return expectedMembers, nil
			},
		}
		uc := NewGetCircleMembersUseCase(repo)

		// Act
		result, err := uc.Execute(context.Background(), GetCircleMembersInput{CircleID: "c1"})

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 3)
	})
}
