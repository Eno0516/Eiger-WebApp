package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/backend/pkg/jwt"
)

type contextKey string

const ClaimsKey contextKey = "claims"

// JWTAuth は Authorization: Bearer <token> を検証するミドルウェア。
// トークンが無効の場合は 401 を返す。
func JWTAuth(jwtGen *jwt.Generator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Authorization ヘッダから "Bearer " プレフィックスを除いてトークンを取り出す
			// TODO: jwtGen.Parse でトークンを検証する（失敗なら 401）
			// TODO: claims を context.WithValue(r.Context(), ClaimsKey, claims) に詰める
			// TODO: next.ServeHTTP(w, r.WithContext(ctx)) で次のハンドラに渡す
			_ = strings.TrimPrefix // import guard
			panic("not implemented")
		})
	}
}

// RequireRole は指定ロール以上のアクセスを要求するミドルウェア。
// JWTAuth の後に使用すること。
func RequireRole(required entity.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: context から claims を取り出す
			// TODO: claims.Role が required 以上かチェックする（不足なら 403）
			//       ロール順: pending < member < admin
			// TODO: next.ServeHTTP に渡す
			panic("not implemented")
		})
	}
}

// ClaimsFromContext は context から JWT Claims を取り出すヘルパー。
func ClaimsFromContext(ctx context.Context) (*jwt.Claims, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*jwt.Claims)
	return claims, ok
}
