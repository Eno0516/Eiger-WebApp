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

func newUserProfileRepo() *postgresrepo.UserProfileRepository {
	return postgresrepo.NewUserProfileRepository(testDB)
}

// seedTestUser は統合テスト用ユーザーを DB に作成して返すヘルパー。
func seedTestUser(t *testing.T, lineUserID string) *entity.User {
	t.Helper()
	repo := newUserRepo()
	u, err := repo.Create(context.Background(), &entity.User{LineUserID: lineUserID})
	require.NoError(t, err)
	return u
}

func TestUserProfileRepository(t *testing.T) {
	t.Run("create - プロフィールが作成され FindByUserID で取得できる", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u := seedTestUser(t, "line-profile-001")
		repo := newUserProfileRepo()

		// Act
		created, err := repo.Create(ctx, &entity.UserProfile{
			UserID:      u.ID,
			DisplayName: "テスト太郎",
			University:  "○○大学",
			Faculty:     "工学部",
			Grade:       2,
		})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, u.ID, created.UserID)
		assert.Equal(t, "テスト太郎", created.DisplayName)
		assert.Equal(t, "○○大学", created.University)
		assert.Equal(t, "工学部", created.Faculty)
		assert.Equal(t, 2, created.Grade)
	})

	t.Run("find by user id - 作成後に取得できる", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u := seedTestUser(t, "line-profile-002")
		repo := newUserProfileRepo()
		_, err := repo.Create(ctx, &entity.UserProfile{
			UserID:      u.ID,
			DisplayName: "テスト花子",
		})
		require.NoError(t, err)

		// Act
		found, err := repo.FindByUserID(ctx, u.ID)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "テスト花子", found.DisplayName)
	})

	t.Run("find by user id - 存在しない ID で ErrNotFound が返る", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		repo := newUserProfileRepo()

		// Act
		_, err := repo.FindByUserID(ctx, "00000000-0000-0000-0000-000000000000")

		// Assert
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
	})

	t.Run("create duplicate - 同じ user_id で ErrConflict が返る", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u := seedTestUser(t, "line-profile-003")
		repo := newUserProfileRepo()
		_, err := repo.Create(ctx, &entity.UserProfile{UserID: u.ID, DisplayName: "初回"})
		require.NoError(t, err)

		// Act
		_, err = repo.Create(ctx, &entity.UserProfile{UserID: u.ID, DisplayName: "二回目"})

		// Assert
		assert.ErrorIs(t, err, apperrors.ErrConflict)
	})

	t.Run("update - DisplayName の変更が反映される", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u := seedTestUser(t, "line-profile-004")
		repo := newUserProfileRepo()
		created, err := repo.Create(ctx, &entity.UserProfile{UserID: u.ID, DisplayName: "変更前"})
		require.NoError(t, err)

		// Act
		created.DisplayName = "変更後"
		updated, err := repo.Update(ctx, created)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "変更後", updated.DisplayName)

		found, err := repo.FindByUserID(ctx, u.ID)
		require.NoError(t, err)
		assert.Equal(t, "変更後", found.DisplayName)
	})

	t.Run("update - 存在しない user_id で ErrNotFound が返る", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		repo := newUserProfileRepo()

		// Act
		_, err := repo.Update(ctx, &entity.UserProfile{
			UserID:      "00000000-0000-0000-0000-000000000000",
			DisplayName: "存在しない",
		})

		// Assert
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
	})

	t.Run("create with optional fields empty - university/faculty/grade が省略できる", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u := seedTestUser(t, "line-profile-005")
		repo := newUserProfileRepo()

		// Act
		created, err := repo.Create(ctx, &entity.UserProfile{
			UserID:      u.ID,
			DisplayName: "オプション省略",
			// University・Faculty・Grade はゼロ値
		})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "", created.University)
		assert.Equal(t, "", created.Faculty)
		assert.Equal(t, 0, created.Grade)
	})
}
