package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/rs/zerolog"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/internal/repository/memory"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

// --- モック定義 ---

// mockUserProfileRepo は domainrepo.UserProfileRepository のテスト用モック。
// memory 実装が完成次第、そちらに切り替える。
type mockUserProfileRepo struct {
	createdProfile *entity.UserProfile
	createErr      error
}

func (m *mockUserProfileRepo) FindByUserID(_ context.Context, _ string) (*entity.UserProfile, error) {
	panic("not called in these tests")
}

func (m *mockUserProfileRepo) Create(_ context.Context, profile *entity.UserProfile) (*entity.UserProfile, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.createdProfile != nil {
		return m.createdProfile, nil
	}
	return profile, nil
}

func (m *mockUserProfileRepo) Update(_ context.Context, profile *entity.UserProfile) (*entity.UserProfile, error) {
	panic("not called in these tests")
}

var _ domainrepo.UserProfileRepository = (*mockUserProfileRepo)(nil)

// --- テストヘルパー ---

func newTestRegisterUseCase(
	userRepo domainrepo.UserRepository,
	profileRepo domainrepo.UserProfileRepository,
) *RegisterUseCase {
	return NewRegisterUseCase(userRepo, profileRepo, zerolog.Nop())
}

// seedUser は UserRepository にテスト用ユーザーを追加して返す。
func seedUser(t *testing.T, repo domainrepo.UserRepository, u *entity.User) *entity.User {
	t.Helper()
	created, err := repo.Create(context.Background(), u)
	require.NoError(t, err, "seedUser")
	return created
}

// --- テストケース ---

func TestRegisterUseCase_Execute(t *testing.T) {
	t.Run("success - UserProfile が作成され is_registered が true になる", func(t *testing.T) {
		// Arrange
		userRepo := memory.NewUserRepository()
		u := seedUser(t, userRepo, &entity.User{LineUserID: "line-user-1"})
		uc := newTestRegisterUseCase(userRepo, &mockUserProfileRepo{})

		input := RegisterInput{
			UserID:      u.ID,
			DisplayName: "テスト太郎",
			University:  "○○大学",
			Faculty:     "工学部",
			Grade:       2,
		}

		// Act
		out, err := uc.Execute(context.Background(), input)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "テスト太郎", out.Profile.DisplayName)

		updated, err := userRepo.FindByID(context.Background(), u.ID)
		require.NoError(t, err)
		assert.True(t, updated.IsRegistered)
	})

	t.Run("already registered - ErrConflict が返る", func(t *testing.T) {
		// Arrange
		userRepo := memory.NewUserRepository()
		u := seedUser(t, userRepo, &entity.User{LineUserID: "line-user-2"})
		u.IsRegistered = true
		_, err := userRepo.Update(context.Background(), u)
		require.NoError(t, err)

		uc := newTestRegisterUseCase(userRepo, &mockUserProfileRepo{})
		input := RegisterInput{
			UserID: u.ID, DisplayName: "テスト花子", University: "○○大学", Faculty: "文学部", Grade: 1,
		}

		// Act
		_, err = uc.Execute(context.Background(), input)

		// Assert
		require.Error(t, err)
		assert.Equal(t, 409, apperrors.HTTPStatusCode(err))
	})

	t.Run("user not found - ErrNotFound が返る", func(t *testing.T) {
		// Arrange
		userRepo := memory.NewUserRepository()
		uc := newTestRegisterUseCase(userRepo, &mockUserProfileRepo{})
		input := RegisterInput{
			UserID: "non-existent-id", DisplayName: "テスト", University: "○○大学", Faculty: "理学部", Grade: 3,
		}

		// Act
		_, err := uc.Execute(context.Background(), input)

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
	})

	t.Run("validation error - DisplayName が空のとき ErrBadRequest が返る", func(t *testing.T) {
		// Arrange
		userRepo := memory.NewUserRepository()
		u := seedUser(t, userRepo, &entity.User{LineUserID: "line-user-3"})
		uc := newTestRegisterUseCase(userRepo, &mockUserProfileRepo{})
		input := RegisterInput{
			UserID: u.ID, DisplayName: "", University: "○○大学", Faculty: "法学部", Grade: 2,
		}

		// Act
		_, err := uc.Execute(context.Background(), input)

		// Assert
		require.Error(t, err)
		assert.Equal(t, 400, apperrors.HTTPStatusCode(err))
	})

	t.Run("profile create error - ErrInternal が返る", func(t *testing.T) {
		// Arrange
		userRepo := memory.NewUserRepository()
		u := seedUser(t, userRepo, &entity.User{LineUserID: "line-user-4"})
		profileRepo := &mockUserProfileRepo{createErr: apperrors.NewInternal("db error", nil)}
		uc := newTestRegisterUseCase(userRepo, profileRepo)
		input := RegisterInput{
			UserID: u.ID, DisplayName: "テスト次郎", University: "○○大学", Faculty: "経済学部", Grade: 4,
		}

		// Act
		_, err := uc.Execute(context.Background(), input)

		// Assert
		require.Error(t, err)
		assert.Equal(t, 500, apperrors.HTTPStatusCode(err))
	})
}

// import guard
var _ = apperrors.ErrNotFound
