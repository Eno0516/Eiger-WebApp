package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/rs/zerolog"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

// --- モック定義 ---

type mockVerifier struct {
	lineUserID string
	err        error
}

func (m *mockVerifier) Verify(_ context.Context, _ string) (string, error) {
	return m.lineUserID, m.err
}

type mockJWTGen struct {
	token string
	err   error
}

func (m *mockJWTGen) Generate(_ string, _ bool) (string, error) {
	return m.token, m.err
}

// mockUserRepo は domainrepo.UserRepository を満たすテスト用モック。
type mockUserRepo struct {
	findByLineUserIDUser *entity.User
	findByLineUserIDErr  error
	createdUser          *entity.User
	createErr            error
}

func (m *mockUserRepo) FindByID(_ context.Context, _ string) (*entity.User, error) {
	panic("not called in these tests")
}

func (m *mockUserRepo) FindByLineUserID(_ context.Context, _ string) (*entity.User, error) {
	return m.findByLineUserIDUser, m.findByLineUserIDErr
}

func (m *mockUserRepo) FindAll(_ context.Context) ([]*entity.User, error) {
	panic("not called in these tests")
}

func (m *mockUserRepo) Create(_ context.Context, user *entity.User) (*entity.User, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.createdUser != nil {
		return m.createdUser, nil
	}
	return user, nil
}

func (m *mockUserRepo) Update(_ context.Context, user *entity.User) (*entity.User, error) {
	panic("not called in these tests")
}

func (m *mockUserRepo) Delete(_ context.Context, _ string) error {
	panic("not called in these tests")
}

var _ domainrepo.UserRepository = (*mockUserRepo)(nil)

// --- テストヘルパー ---

func newTestUseCase(repo domainrepo.UserRepository, verifier LINEIDTokenVerifier, jwtGen JWTGenerator) *LineAuthUseCase {
	return NewLineAuthUseCase(repo, verifier, jwtGen, zerolog.Nop())
}

// --- テストケース ---

func TestLineAuthUseCase_Execute(t *testing.T) {
	t.Run("existing user - JWT を返し IsNewUser=false", func(t *testing.T) {
		// Arrange
		existingUser := &entity.User{ID: "user-1", LineUserID: "line-user-123"}
		repo := &mockUserRepo{findByLineUserIDUser: existingUser}
		verifier := &mockVerifier{lineUserID: "line-user-123"}
		jwtGen := &mockJWTGen{token: "jwt-token"}
		uc := newTestUseCase(repo, verifier, jwtGen)

		// Act
		out, err := uc.Execute(context.Background(), LineAuthInput{IDToken: "valid-token"})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "jwt-token", out.Token)
		assert.False(t, out.IsNewUser)
	})

	t.Run("new user - ユーザーを作成し IsNewUser=true", func(t *testing.T) {
		// Arrange
		repo := &mockUserRepo{findByLineUserIDErr: apperrors.ErrNotFound}
		verifier := &mockVerifier{lineUserID: "new-line-user"}
		jwtGen := &mockJWTGen{token: "new-jwt"}
		uc := newTestUseCase(repo, verifier, jwtGen)

		// Act
		out, err := uc.Execute(context.Background(), LineAuthInput{IDToken: "valid-token"})

		// Assert
		require.NoError(t, err)
		assert.True(t, out.IsNewUser)
		assert.Equal(t, "new-jwt", out.Token)
	})

	t.Run("invalid id token - 401 エラー", func(t *testing.T) {
		// Arrange
		repo := &mockUserRepo{}
		verifier := &mockVerifier{err: apperrors.ErrUnauthorized}
		jwtGen := &mockJWTGen{}
		uc := newTestUseCase(repo, verifier, jwtGen)

		// Act
		_, err := uc.Execute(context.Background(), LineAuthInput{IDToken: "invalid-token"})

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, apperrors.ErrUnauthorized)
	})

	t.Run("repository error - 内部エラーとして伝播", func(t *testing.T) {
		// Arrange
		repo := &mockUserRepo{findByLineUserIDErr: apperrors.ErrInternal}
		verifier := &mockVerifier{lineUserID: "line-user-123"}
		jwtGen := &mockJWTGen{}
		uc := newTestUseCase(repo, verifier, jwtGen)

		// Act
		_, err := uc.Execute(context.Background(), LineAuthInput{IDToken: "valid-token"})

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, apperrors.ErrInternal)
	})

	t.Run("jwt generation error - エラーが返る", func(t *testing.T) {
		// Arrange
		existingUser := &entity.User{ID: "user-1", LineUserID: "line-user-123"}
		repo := &mockUserRepo{findByLineUserIDUser: existingUser}
		verifier := &mockVerifier{lineUserID: "line-user-123"}
		jwtGen := &mockJWTGen{err: apperrors.ErrInternal}
		uc := newTestUseCase(repo, verifier, jwtGen)

		// Act
		_, err := uc.Execute(context.Background(), LineAuthInput{IDToken: "valid-token"})

		// Assert
		require.Error(t, err)
	})
}
