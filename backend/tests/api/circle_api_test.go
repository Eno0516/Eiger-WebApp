//go:build api

package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/internal/handler/http/dto"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// モック定義
// =============================================================================

type mockCircleRepoAPI struct {
	createErr  error
	findByIDFn func(ctx context.Context, id string) (*entity.Circle, error)
	findAllFn  func(ctx context.Context) ([]*entity.Circle, error)
	deleteErr  error
}

func (m *mockCircleRepoAPI) Create(_ context.Context, _ *entity.Circle) error { return m.createErr }
func (m *mockCircleRepoAPI) FindByID(ctx context.Context, id string) (*entity.Circle, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, apperrors.ErrNotFound
}
func (m *mockCircleRepoAPI) FindAll(ctx context.Context) ([]*entity.Circle, error) {
	if m.findAllFn != nil {
		return m.findAllFn(ctx)
	}
	return nil, nil
}
func (m *mockCircleRepoAPI) Update(_ context.Context, _ *entity.Circle) error { return nil }
func (m *mockCircleRepoAPI) Delete(_ context.Context, _ string) error         { return m.deleteErr }

var _ domainrepo.CircleRepository = (*mockCircleRepoAPI)(nil)

type mockCircleMemberRepoAPI struct {
	findByCircleAndUserFn func(ctx context.Context, circleID, userID string) (*entity.CircleMember, error)
	findByCircleFn        func(ctx context.Context, circleID string) ([]*entity.CircleMember, error)
	findByUserFn          func(ctx context.Context, userID string) ([]*entity.CircleMember, error)
	upsertErr             error
}

func (m *mockCircleMemberRepoAPI) Upsert(_ context.Context, _ *entity.CircleMember) error {
	return m.upsertErr
}
func (m *mockCircleMemberRepoAPI) FindByCircleAndUser(ctx context.Context, circleID, userID string) (*entity.CircleMember, error) {
	if m.findByCircleAndUserFn != nil {
		return m.findByCircleAndUserFn(ctx, circleID, userID)
	}
	return nil, apperrors.ErrNotFound
}
func (m *mockCircleMemberRepoAPI) FindByCircle(ctx context.Context, circleID string) ([]*entity.CircleMember, error) {
	if m.findByCircleFn != nil {
		return m.findByCircleFn(ctx, circleID)
	}
	return nil, nil
}
func (m *mockCircleMemberRepoAPI) FindByUser(ctx context.Context, userID string) ([]*entity.CircleMember, error) {
	if m.findByUserFn != nil {
		return m.findByUserFn(ctx, userID)
	}
	return nil, nil
}
func (m *mockCircleMemberRepoAPI) FindAdminsByCircle(_ context.Context, _ string) ([]*entity.CircleMember, error) {
	return nil, nil
}
func (m *mockCircleMemberRepoAPI) Delete(_ context.Context, _, _ string) error { return nil }

var _ domainrepo.CircleMemberRepository = (*mockCircleMemberRepoAPI)(nil)

// adminCircleRepo は system_admin テスト用の最小 circle repo。
func adminCircleRepo() *mockCircleRepoAPI {
	return &mockCircleRepoAPI{
		findAllFn: func(_ context.Context) ([]*entity.Circle, error) {
			return []*entity.Circle{{ID: "c1", Name: "テストサークル"}}, nil
		},
		findByIDFn: func(_ context.Context, id string) (*entity.Circle, error) {
			return &entity.Circle{ID: id, Name: "テストサークル"}, nil
		},
	}
}

// =============================================================================
// GET /circles  （JWT 必須）
// =============================================================================

func TestGetCirclesAPI(t *testing.T) {
	t.Run("success - サークル一覧が返る", func(t *testing.T) {
		// Arrange
		circleRepo := &mockCircleRepoAPI{
			findAllFn: func(_ context.Context) ([]*entity.Circle, error) {
				return []*entity.Circle{{ID: "c1", Name: "サークルA"}}, nil
			},
		}
		memberRepo := &mockCircleMemberRepoAPI{}
		srv := newTestServer(t, apiTestConfig{circleRepo: circleRepo, circleMemberRepo: memberRepo})
		defer srv.Close()
		token := makeToken(t, "user-1", true) // system_admin で全件取得

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/circles", token, nil)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body []dto.CircleResponse
		decodeJSON(t, resp, &body)
		assert.Len(t, body, 1)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       &mockCircleRepoAPI{},
			circleMemberRepo: &mockCircleMemberRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/circles", "", nil)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

// =============================================================================
// POST /circles/{id}/apply  （JWT 必須）
// =============================================================================

func TestApplyCircleAPI(t *testing.T) {
	t.Run("success - 200 + applied が返る", func(t *testing.T) {
		// Arrange
		circleRepo := &mockCircleRepoAPI{
			findByIDFn: func(_ context.Context, id string) (*entity.Circle, error) {
				return &entity.Circle{ID: id, Name: "サークルA"}, nil
			},
		}
		memberRepo := &mockCircleMemberRepoAPI{}
		srv := newTestServer(t, apiTestConfig{circleRepo: circleRepo, circleMemberRepo: memberRepo})
		defer srv.Close()
		token := makeToken(t, "user-1", false)

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/circles/c1/apply", token, nil)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body dto.MessageResponse
		decodeJSON(t, resp, &body)
		assert.Equal(t, "applied", body.Message)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       &mockCircleRepoAPI{},
			circleMemberRepo: &mockCircleMemberRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/circles/c1/apply", "", nil)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

// =============================================================================
// GET /circles/{id}/members  （JWT + circle_admin 必須）
// =============================================================================

func TestGetCircleMembersAPI(t *testing.T) {
	t.Run("success - メンバー一覧が返る", func(t *testing.T) {
		// Arrange
		memberRepo := &mockCircleMemberRepoAPI{
			findByCircleAndUserFn: func(_ context.Context, _, _ string) (*entity.CircleMember, error) {
				return &entity.CircleMember{Role: entity.CircleRoleAdmin}, nil
			},
			findByCircleFn: func(_ context.Context, _ string) ([]*entity.CircleMember, error) {
				return []*entity.CircleMember{
					{CircleID: "c1", UserID: "u1", Role: entity.CircleRoleMember},
				}, nil
			},
		}
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       &mockCircleRepoAPI{},
			circleMemberRepo: memberRepo,
		})
		defer srv.Close()
		token := makeToken(t, "admin-user", false)

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/circles/c1/members", token, nil)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body []dto.CircleMemberResponse
		decodeJSON(t, resp, &body)
		assert.Len(t, body, 1)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       &mockCircleRepoAPI{},
			circleMemberRepo: &mockCircleMemberRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/circles/c1/members", "", nil)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("circle_admin でない - 403 が返る", func(t *testing.T) {
		// Arrange
		memberRepo := &mockCircleMemberRepoAPI{
			findByCircleAndUserFn: func(_ context.Context, _, _ string) (*entity.CircleMember, error) {
				return nil, apperrors.ErrNotFound
			},
		}
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       &mockCircleRepoAPI{},
			circleMemberRepo: memberRepo,
		})
		defer srv.Close()
		token := makeToken(t, "regular-user", false)

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/circles/c1/members", token, nil)

		// Assert
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

// =============================================================================
// PATCH /circles/{id}/members/{userId}/role  （JWT + circle_admin 必須）
// =============================================================================

func TestUpdateCircleMemberRoleAPI(t *testing.T) {
	t.Run("success - 200 + role updated が返る", func(t *testing.T) {
		// Arrange
		memberRepo := &mockCircleMemberRepoAPI{
			findByCircleAndUserFn: func(_ context.Context, _, _ string) (*entity.CircleMember, error) {
				return &entity.CircleMember{Role: entity.CircleRoleAdmin}, nil
			},
		}
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       &mockCircleRepoAPI{},
			circleMemberRepo: memberRepo,
		})
		defer srv.Close()
		token := makeToken(t, "admin-user", false)

		// Act
		resp := doRequest(t, http.MethodPatch, srv.URL+"/api/v1/circles/c1/members/target-user/role", token,
			dto.UpdateCircleMemberRoleRequest{Role: "member"},
		)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       &mockCircleRepoAPI{},
			circleMemberRepo: &mockCircleMemberRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodPatch, srv.URL+"/api/v1/circles/c1/members/u1/role", "",
			dto.UpdateCircleMemberRoleRequest{Role: "member"},
		)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

// =============================================================================
// POST /admin/circles  （JWT + system_admin 必須）
// =============================================================================

func TestCreateAdminCircleAPI(t *testing.T) {
	t.Run("success - 201 + 作成されたサークルが返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       &mockCircleRepoAPI{},
			circleMemberRepo: &mockCircleMemberRepoAPI{},
		})
		defer srv.Close()
		token := makeToken(t, "admin-user", true)

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/admin/circles", token,
			dto.CreateCircleRequest{Name: "新サークル"},
		)

		// Assert
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		var body dto.CircleResponse
		decodeJSON(t, resp, &body)
		assert.Equal(t, "新サークル", body.Name)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       &mockCircleRepoAPI{},
			circleMemberRepo: &mockCircleMemberRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/admin/circles", "",
			dto.CreateCircleRequest{Name: "新サークル"},
		)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("system_admin でない - 403 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       &mockCircleRepoAPI{},
			circleMemberRepo: &mockCircleMemberRepoAPI{},
		})
		defer srv.Close()
		token := makeToken(t, "regular-user", false)

		// Act
		resp := doRequest(t, http.MethodPost, srv.URL+"/api/v1/admin/circles", token,
			dto.CreateCircleRequest{Name: "新サークル"},
		)

		// Assert
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

// =============================================================================
// DELETE /admin/circles/{id}  （JWT + system_admin 必須）
// =============================================================================

func TestDeleteAdminCircleAPI(t *testing.T) {
	t.Run("success - 204 が返る", func(t *testing.T) {
		// Arrange
		circleRepo := &mockCircleRepoAPI{
			findByIDFn: func(_ context.Context, id string) (*entity.Circle, error) {
				return &entity.Circle{ID: id, Name: "サークルA"}, nil
			},
		}
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       circleRepo,
			circleMemberRepo: &mockCircleMemberRepoAPI{},
		})
		defer srv.Close()
		token := makeToken(t, "admin-user", true)

		// Act
		resp := doRequest(t, http.MethodDelete, srv.URL+"/api/v1/admin/circles/c1", token, nil)

		// Assert
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("system_admin でない - 403 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       &mockCircleRepoAPI{},
			circleMemberRepo: &mockCircleMemberRepoAPI{},
		})
		defer srv.Close()
		token := makeToken(t, "regular-user", false)

		// Act
		resp := doRequest(t, http.MethodDelete, srv.URL+"/api/v1/admin/circles/c1", token, nil)

		// Assert
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

// =============================================================================
// GET /admin/circles  （JWT + system_admin or circle_admin 必須）
// =============================================================================

func TestGetAdminCirclesAPI(t *testing.T) {
	t.Run("success - system_admin で全件返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       adminCircleRepo(),
			circleMemberRepo: &mockCircleMemberRepoAPI{},
		})
		defer srv.Close()
		token := makeToken(t, "admin-user", true)

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/admin/circles", token, nil)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body []dto.CircleResponse
		decodeJSON(t, resp, &body)
		assert.NotEmpty(t, body)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       adminCircleRepo(),
			circleMemberRepo: &mockCircleMemberRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/admin/circles", "", nil)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

// =============================================================================
// GET /admin/circles/{id}  （JWT + system_admin or circle_admin 必須）
// =============================================================================

func TestGetAdminCircleAPI(t *testing.T) {
	t.Run("success - サークル詳細が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       adminCircleRepo(),
			circleMemberRepo: &mockCircleMemberRepoAPI{},
		})
		defer srv.Close()
		token := makeToken(t, "admin-user", true)

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/admin/circles/c1", token, nil)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body dto.CircleResponse
		decodeJSON(t, resp, &body)
		assert.Equal(t, "c1", body.ID)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       adminCircleRepo(),
			circleMemberRepo: &mockCircleMemberRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/admin/circles/c1", "", nil)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

// =============================================================================
// GET /admin/circles/{id}/members  （JWT + system_admin or circle_admin 必須）
// =============================================================================

func TestGetAdminCircleMembersAPI(t *testing.T) {
	t.Run("success - system_admin でメンバー一覧が返る", func(t *testing.T) {
		// Arrange
		memberRepo := &mockCircleMemberRepoAPI{
			findByCircleFn: func(_ context.Context, _ string) ([]*entity.CircleMember, error) {
				return []*entity.CircleMember{
					{CircleID: "c1", UserID: "u1", Role: entity.CircleRoleMember},
				}, nil
			},
		}
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       adminCircleRepo(),
			circleMemberRepo: memberRepo,
		})
		defer srv.Close()
		token := makeToken(t, "admin-user", true)

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/admin/circles/c1/members", token, nil)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body []dto.CircleMemberResponse
		decodeJSON(t, resp, &body)
		assert.Len(t, body, 1)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       adminCircleRepo(),
			circleMemberRepo: &mockCircleMemberRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodGet, srv.URL+"/api/v1/admin/circles/c1/members", "", nil)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

// =============================================================================
// PATCH /admin/circles/{id}/members/{userId}/role  （JWT + system_admin or circle_admin 必須）
// =============================================================================

func TestUpdateAdminMemberRoleAPI(t *testing.T) {
	t.Run("success - 200 + role updated が返る", func(t *testing.T) {
		// Arrange
		memberRepo := &mockCircleMemberRepoAPI{
			findByCircleAndUserFn: func(_ context.Context, _, _ string) (*entity.CircleMember, error) {
				return &entity.CircleMember{CircleID: "c1", UserID: "u1", Role: entity.CircleRoleMember}, nil
			},
		}
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       adminCircleRepo(),
			circleMemberRepo: memberRepo,
		})
		defer srv.Close()
		token := makeToken(t, "admin-user", true)

		// Act
		resp := doRequest(t, http.MethodPatch, srv.URL+"/api/v1/admin/circles/c1/members/u1/role", token,
			dto.UpdateCircleMemberRoleRequest{Role: "member"},
		)

		// Assert
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("no JWT - 401 が返る", func(t *testing.T) {
		// Arrange
		srv := newTestServer(t, apiTestConfig{
			circleRepo:       adminCircleRepo(),
			circleMemberRepo: &mockCircleMemberRepoAPI{},
		})
		defer srv.Close()

		// Act
		resp := doRequest(t, http.MethodPatch, srv.URL+"/api/v1/admin/circles/c1/members/u1/role", "",
			dto.UpdateCircleMemberRoleRequest{Role: "member"},
		)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
