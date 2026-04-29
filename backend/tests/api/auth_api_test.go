//go:build api

package api_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/backend/internal/handler/http/dto"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

// =============================================================================
// テストケース
// =============================================================================

func TestLineAuthAPI(t *testing.T) {
	t.Run("existing user - 200 + is_new_user=false が返る", func(t *testing.T) {
		// Arrange
		existingUser := &entity.User{ID: "user-1", LineUserID: "line-123", DisplayName: "テストユーザー"}
		srv := newTestServer(t, apiTestConfig{
			userRepo:     &mockLineAuthUserRepo{findByLineUserIDUser: existingUser},
			lineVerifier: &mockVerifier{lineUserID: "line-123"},
			jwtGenUC:     &mockJWTGen{token: "jwt-token"},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/auth/line", "", dto.LineAuthRequest{IDToken: "valid-token"})

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body dto.LoginResponse
		decodeJSON(t, resp, &body)
		assert.Equal(t, "jwt-token", body.Token)
		assert.False(t, body.IsNewUser)
	})

	t.Run("new user - 200 + is_new_user=true が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			userRepo:     &mockLineAuthUserRepo{findByLineUserIDErr: apperrors.ErrNotFound},
			lineVerifier: &mockVerifier{lineUserID: "line-456"},
			jwtGenUC:     &mockJWTGen{token: "new-jwt-token"},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/auth/line", "", dto.LineAuthRequest{IDToken: "valid-token"})

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body dto.LoginResponse
		decodeJSON(t, resp, &body)
		assert.Equal(t, "new-jwt-token", body.Token)
		assert.True(t, body.IsNewUser)
	})

	t.Run("invalid id token - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			userRepo:     &mockLineAuthUserRepo{},
			lineVerifier: &mockVerifier{err: apperrors.NewUnauthorized("invalid token", nil)},
			jwtGenUC:     &mockJWTGen{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/auth/line", "", dto.LineAuthRequest{IDToken: "bad-token"})

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("invalid request body - 400 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			userRepo:     &mockLineAuthUserRepo{},
			lineVerifier: &mockVerifier{},
			jwtGenUC:     &mockJWTGen{},
		})
		defer srv.Close()

		// Act
		req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/auth/line", strings.NewReader("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestLogoutAPI(t *testing.T) {
	t.Run("success - 200 + logged out メッセージが返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/auth/logout", "", nil)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body dto.MessageResponse
		decodeJSON(t, resp, &body)
		assert.Equal(t, "logged out", body.Message)
	})
}

func TestRegisterAPI(t *testing.T) {
	t.Run("success - 201 + display_name が返る", func(t *testing.T) {
		// Arrange
		user := &entity.User{ID: "user-1", LineUserID: "line-1", IsRegistered: false}
		srv := newTestServer(t, apiTestConfig{
			userRepo:        &mockFullUserRepo{findByIDUser: user, updateUser: user},
			userProfileRepo: &mockUserProfileRepoAPI{},
		})
		defer srv.Close()
		token := makeToken(t, "user-1", false)

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/auth/signup", token, dto.RegisterRequest{
			DisplayName: "テスト太郎",
			University:  "○○大学",
			Faculty:     "工学部",
			Grade:       2,
		})

		// Assert
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		var body dto.RegisterResponse
		decodeJSON(t, resp, &body)
		assert.Equal(t, "テスト太郎", body.DisplayName)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			userRepo:        &mockFullUserRepo{},
			userProfileRepo: &mockUserProfileRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/auth/signup", "", dto.RegisterRequest{DisplayName: "テスト"})

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("already registered - 409 が返る", func(t *testing.T) {
		// Arrange
		registeredUser := &entity.User{ID: "user-2", LineUserID: "line-2", IsRegistered: true}
		srv := newTestServer(t, apiTestConfig{
			userRepo:        &mockFullUserRepo{findByIDUser: registeredUser},
			userProfileRepo: &mockUserProfileRepoAPI{},
		})
		defer srv.Close()
		token := makeToken(t, "user-2", false)

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/auth/signup", token, dto.RegisterRequest{
			DisplayName: "テスト花子", University: "○○大学", Faculty: "文学部", Grade: 1,
		})

		// Assert
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})
}
