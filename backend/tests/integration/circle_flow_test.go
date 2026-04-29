//go:build integration

package integration_test

import (
	"context"
	"testing"

	postgresrepo "github.com/enoro/eiger-webapp/backend/internal/repository/postgres"
	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCircleRepo() *postgresrepo.CircleRepository {
	return postgresrepo.NewCircleRepository(testDB)
}

func newCircleMemberRepo() *postgresrepo.CircleMemberRepository {
	return postgresrepo.NewCircleMemberRepository(testDB)
}

// seedTestCircle は統合テスト用サークルを DB に作成して返すヘルパー。
func seedTestCircle(t *testing.T, name string) *entity.Circle {
	t.Helper()
	repo := newCircleRepo()
	c := &entity.Circle{Name: name}
	err := repo.Create(context.Background(), c)
	require.NoError(t, err)
	// Create は ID を返さないため FindAll で取得する
	all, err := repo.FindAll(context.Background())
	require.NoError(t, err)
	for _, circle := range all {
		if circle.Name == name {
			return circle
		}
	}
	t.Fatalf("seedTestCircle: created circle %q not found", name)
	return nil
}

// =============================================================================
// CircleRepository
// =============================================================================

func TestCircleRepository(t *testing.T) {
	t.Run("create - サークルが作成され FindAll で取得できる", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		repo := newCircleRepo()

		// Act
		err := repo.Create(ctx, &entity.Circle{Name: "テストサークルA", Description: "説明"})

		// Assert
		require.NoError(t, err)
		all, err := repo.FindAll(ctx)
		require.NoError(t, err)
		names := make([]string, 0, len(all))
		for _, c := range all {
			names = append(names, c.Name)
		}
		assert.Contains(t, names, "テストサークルA")
	})

	t.Run("find by id - 作成後に取得できる", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		c := seedTestCircle(t, "テストサークルB")

		// Act
		found, err := newCircleRepo().FindByID(ctx, c.ID)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "テストサークルB", found.Name)
	})

	t.Run("find by id - 存在しない ID で ErrNotFound が返る", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		repo := newCircleRepo()

		// Act
		_, err := repo.FindByID(ctx, "00000000-0000-0000-0000-000000000000")

		// Assert
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
	})

	t.Run("update - Description の変更が反映される", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		c := seedTestCircle(t, "テストサークルC")
		repo := newCircleRepo()

		// Act
		c.Description = "更新後の説明"
		err := repo.Update(ctx, c)

		// Assert
		require.NoError(t, err)
		found, err := repo.FindByID(ctx, c.ID)
		require.NoError(t, err)
		assert.Equal(t, "更新後の説明", found.Description)
	})

	t.Run("delete - 削除後に FindByID で ErrNotFound が返る", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		c := seedTestCircle(t, "テストサークルD")
		repo := newCircleRepo()

		// Act
		err := repo.Delete(ctx, c.ID)

		// Assert
		require.NoError(t, err)
		_, err = repo.FindByID(ctx, c.ID)
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
	})
}

// =============================================================================
// CircleMemberRepository
// =============================================================================

func TestCircleMemberRepository(t *testing.T) {
	t.Run("upsert - メンバーが追加され FindByCircleAndUser で取得できる", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u := seedTestUser(t, "line-member-001")
		c := seedTestCircle(t, "メンバーテストサークルA")
		repo := newCircleMemberRepo()

		// Act
		err := repo.Upsert(ctx, &entity.CircleMember{
			CircleID: c.ID,
			UserID:   u.ID,
			Role:     entity.CircleRoleMember,
		})

		// Assert
		require.NoError(t, err)
		found, err := repo.FindByCircleAndUser(ctx, c.ID, u.ID)
		require.NoError(t, err)
		assert.Equal(t, entity.CircleRoleMember, found.Role)
	})

	t.Run("upsert - ロール変更が反映される", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u := seedTestUser(t, "line-member-002")
		c := seedTestCircle(t, "メンバーテストサークルB")
		repo := newCircleMemberRepo()
		err := repo.Upsert(ctx, &entity.CircleMember{CircleID: c.ID, UserID: u.ID, Role: entity.CircleRoleMember})
		require.NoError(t, err)

		// Act
		err = repo.Upsert(ctx, &entity.CircleMember{CircleID: c.ID, UserID: u.ID, Role: entity.CircleRoleAdmin})

		// Assert
		require.NoError(t, err)
		found, err := repo.FindByCircleAndUser(ctx, c.ID, u.ID)
		require.NoError(t, err)
		assert.Equal(t, entity.CircleRoleAdmin, found.Role)
	})

	t.Run("find by circle and user - 存在しない組み合わせで ErrNotFound が返る", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u := seedTestUser(t, "line-member-003")
		c := seedTestCircle(t, "メンバーテストサークルC")
		repo := newCircleMemberRepo()

		// Act
		_, err := repo.FindByCircleAndUser(ctx, c.ID, u.ID)

		// Assert
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
	})

	t.Run("find by circle - サークルの全メンバーが返る", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u1 := seedTestUser(t, "line-member-004")
		u2 := seedTestUser(t, "line-member-005")
		c := seedTestCircle(t, "メンバーテストサークルD")
		repo := newCircleMemberRepo()
		require.NoError(t, repo.Upsert(ctx, &entity.CircleMember{CircleID: c.ID, UserID: u1.ID, Role: entity.CircleRoleMember}))
		require.NoError(t, repo.Upsert(ctx, &entity.CircleMember{CircleID: c.ID, UserID: u2.ID, Role: entity.CircleRoleAdmin}))

		// Act
		members, err := repo.FindByCircle(ctx, c.ID)

		// Assert
		require.NoError(t, err)
		assert.Len(t, members, 2)
	})

	t.Run("find admins by circle - admin ロールのみ返る", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u1 := seedTestUser(t, "line-member-006")
		u2 := seedTestUser(t, "line-member-007")
		c := seedTestCircle(t, "メンバーテストサークルE")
		repo := newCircleMemberRepo()
		require.NoError(t, repo.Upsert(ctx, &entity.CircleMember{CircleID: c.ID, UserID: u1.ID, Role: entity.CircleRoleMember}))
		require.NoError(t, repo.Upsert(ctx, &entity.CircleMember{CircleID: c.ID, UserID: u2.ID, Role: entity.CircleRoleAdmin}))

		// Act
		admins, err := repo.FindAdminsByCircle(ctx, c.ID)

		// Assert
		require.NoError(t, err)
		assert.Len(t, admins, 1)
		assert.Equal(t, entity.CircleRoleAdmin, admins[0].Role)
	})

	t.Run("delete - 削除後に FindByCircleAndUser で ErrNotFound が返る", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u := seedTestUser(t, "line-member-008")
		c := seedTestCircle(t, "メンバーテストサークルF")
		repo := newCircleMemberRepo()
		require.NoError(t, repo.Upsert(ctx, &entity.CircleMember{CircleID: c.ID, UserID: u.ID, Role: entity.CircleRoleMember}))

		// Act
		err := repo.Delete(ctx, c.ID, u.ID)

		// Assert
		require.NoError(t, err)
		_, err = repo.FindByCircleAndUser(ctx, c.ID, u.ID)
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
	})
}
