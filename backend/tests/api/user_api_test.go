//go:build api

package api_test

import (
	"net/http"
	"testing"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/internal/handler/http/dto"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"context"
)

// mockUserGetRepo は FindByID のみ使う UserRepository モック（user endpoint テスト用）。
type mockUserGetRepo struct {
	findByIDUser *entity.User
	findByIDErr  error
}

func (m *mockUserGetRepo) FindByID(_ context.Context, _ string) (*entity.User, error) {
	return m.findByIDUser, m.findByIDErr
}
func (m *mockUserGetRepo) FindByLineUserID(_ context.Context, _ string) (*entity.User, error) {
	panic("not called in these tests")
}
func (m *mockUserGetRepo) FindAll(_ context.Context) ([]*entity.User, error) {
	panic("not called in these tests")
}
func (m *mockUserGetRepo) Create(_ context.Context, u *entity.User) (*entity.User, error) {
	panic("not called in these tests")
}
func (m *mockUserGetRepo) Update(_ context.Context, u *entity.User) (*entity.User, error) {
	panic("not called in these tests")
}
func (m *mockUserGetRepo) Delete(_ context.Context, _ string) error {
	panic("not called in these tests")
}

var _ domainrepo.UserRepository = (*mockUserGetRepo)(nil)

// =============================================================================
// GET /users/me
// =============================================================================

func TestGetMeAPI(t *testing.T) {
	t.Run("success - ユーザー情報が返る", func(t *testing.T) {
		// Arrange
		user := &entity.User{ID: "user-1", IsRegistered: true, IsSystemAdmin: false}
		srv := newTestServer(t, apiTestConfig{
			userRepo: &mockUserGetRepo{findByIDUser: user},
		})
		defer srv.Close()
		token := makeToken(t, "user-1", false)

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/users/me", token, nil)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body dto.MeResponse
		decodeJSON(t, resp, &body)
		assert.Equal(t, "user-1", body.ID)
		assert.True(t, body.IsRegistered)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{userRepo: &mockUserGetRepo{}})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/users/me", "", nil)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("user not found - 404 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			userRepo: &mockUserGetRepo{findByIDErr: apperrors.ErrNotFound},
		})
		defer srv.Close()
		token := makeToken(t, "ghost-user", false)

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/users/me", token, nil)

		// Assert
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

// =============================================================================
// GET /users/{id}
// =============================================================================

func TestGetUserByIDAPI(t *testing.T) {
	t.Run("success - ユーザー情報が返る", func(t *testing.T) {
		// Arrange
		user := &entity.User{ID: "user-1", IsRegistered: true}
		srv := newTestServer(t, apiTestConfig{
			userRepo: &mockUserGetRepo{findByIDUser: user},
		})
		defer srv.Close()
		token := makeToken(t, "user-1", false)

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/users/user-1", token, nil)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body dto.UserResponse
		decodeJSON(t, resp, &body)
		assert.Equal(t, "user-1", body.ID)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{userRepo: &mockUserGetRepo{}})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/users/user-1", "", nil)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("not found - 404 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			userRepo: &mockUserGetRepo{findByIDErr: apperrors.ErrNotFound},
		})
		defer srv.Close()
		token := makeToken(t, "caller", false)

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/users/ghost-id", token, nil)

		// Assert
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

// =============================================================================
// PATCH /users/{id}/display-name
// =============================================================================

func TestUpdateDisplayNameAPI(t *testing.T) {
	t.Run("not implemented - 501 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{userRepo: &mockUserGetRepo{}})
		defer srv.Close()
		token := makeToken(t, "user-1", false)

		// Act
		resp := doRequest(t, http.MethodPatch, srv.URL+"/api/v1/users/user-1/display-name", token,
			dto.UpdateDisplayNameRequest{DisplayName: "新しい名前"},
		)

		// Assert
		assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{userRepo: &mockUserGetRepo{}})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodPatch, srv.URL+"/api/v1/users/user-1/display-name", "",
			dto.UpdateDisplayNameRequest{DisplayName: "新しい名前"},
		)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("forbidden - 他ユーザーの ID を指定すると 403 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{userRepo: &mockUserGetRepo{}})
		defer srv.Close()
		token := makeToken(t, "other-user", false)

		// Act
		resp := doRequest(t, http.MethodPatch, srv.URL+"/api/v1/users/user-1/display-name", token,
			dto.UpdateDisplayNameRequest{DisplayName: "新しい名前"},
		)

		// Assert
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}
