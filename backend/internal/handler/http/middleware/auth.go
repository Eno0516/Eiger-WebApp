package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/jwt"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"}) //nolint:errcheck
}

func writeForbidden(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]string{"error": "forbidden"}) //nolint:errcheck
}

// JWTAuth は Authorization: Bearer <token> を検証するミドルウェア。
// トークンが無効の場合は 401 を返す。
func JWTAuth(jwtGen *jwt.Generator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == authHeader || tokenStr == "" {
				writeUnauthorized(w)
				return
			}
			claims, err := jwtGen.Parse(tokenStr)
			if err != nil {
				writeUnauthorized(w)
				return
			}
			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireSystemAdmin は is_system_admin = true のユーザーのみ通すミドルウェア。
// JWTAuth の後に使用すること。
func RequireSystemAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok || !claims.IsSystemAdmin {
				writeForbidden(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireCircleAdmin はサークル内ロールを検証するミドルウェア。
// JWTAuth の後に使用し、URL パラメータ "id" を circleID として使用する。
func RequireCircleAdmin(repo domainrepo.CircleMemberRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				writeUnauthorized(w)
				return
			}
			circleID := chi.URLParam(r, "id")
			member, err := repo.FindByCircleAndUser(r.Context(), circleID, claims.UserID)
			if err != nil || member.Role != entity.CircleRoleAdmin {
				writeForbidden(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireSystemAdminOrCircleAdmin は system_admin、または指定サークルの circle_admin のみ通すミドルウェア。
func RequireSystemAdminOrCircleAdmin(repo domainrepo.CircleMemberRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				writeUnauthorized(w)
				return
			}
			if claims.IsSystemAdmin {
				next.ServeHTTP(w, r)
				return
			}
			circleID := chi.URLParam(r, "id")
			member, err := repo.FindByCircleAndUser(r.Context(), circleID, claims.UserID)
			if err != nil || member.Role != entity.CircleRoleAdmin {
				writeForbidden(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ClaimsFromContext は context から JWT Claims を取り出すヘルパー。
func ClaimsFromContext(ctx context.Context) (*jwt.Claims, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*jwt.Claims)
	return claims, ok
}
