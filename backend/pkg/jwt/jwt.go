package jwt

import (
	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
)

// Claims は JWT の payload。
type Claims struct {
	UserID string          `json:"user_id"`
	Email  string          `json:"email"`
	Role   entity.UserRole `json:"role"`
}

// Generator は JWT の生成・検証を行う。
// usecase/auth.JWTGenerator interface を実装する。
type Generator struct {
	secretKey string
	// TODO: 有効期限を設定できるようにする（デフォルト 1h）
}

func NewGenerator(secretKey string) *Generator {
	return &Generator{secretKey: secretKey}
}

// Generate は Claims を元に署名済み JWT 文字列を生成する。
// アルゴリズム: HS256、有効期限: 1時間。
func (g *Generator) Generate(userID, email string, role entity.UserRole) (string, error) {
	// TODO: github.com/golang-jwt/jwt/v5 を使って Claims を組み立てる
	// TODO: jwt.NewWithClaims(jwt.SigningMethodHS256, claims) で署名する
	// TODO: token.SignedString([]byte(g.secretKey)) で文字列を返す
	panic("not implemented")
}

// Parse は JWT 文字列を検証して Claims を返す。
// 署名不正・期限切れの場合はエラーを返す。
func (g *Generator) Parse(tokenStr string) (*Claims, error) {
	// TODO: jwt.ParseWithClaims で検証する
	// TODO: claims を *Claims にキャストして返す
	panic("not implemented")
}
