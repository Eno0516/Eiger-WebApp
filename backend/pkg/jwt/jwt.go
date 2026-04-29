package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims は JWT の payload。
// サークル内ロールはリクエストのたびに DB から取得するため含めない。
type Claims struct {
	UserID        string `json:"user_id"`
	IsSystemAdmin bool   `json:"is_system_admin"`
	jwt.RegisteredClaims
}

// Generator は JWT の生成・検証を行う。
type Generator struct {
	secretKey string
	expiry    time.Duration
}

// NewGenerator は Generator を生成する。expiry に 0 を渡すとデフォルト 1h になる。
func NewGenerator(secretKey string, expiry time.Duration) *Generator {
	if expiry == 0 {
		expiry = time.Hour
	}
	return &Generator{secretKey: secretKey, expiry: expiry}
}

// Generate は Claims を元に署名済み JWT 文字列を生成する。
// アルゴリズム: HS256、有効期限: g.expiry。
func (g *Generator) Generate(userID string, isSystemAdmin bool) (string, error) {
	claims := Claims{
		UserID:        userID,
		IsSystemAdmin: isSystemAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.expiry)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(g.secretKey))
}

// Parse は JWT 文字列を検証して Claims を返す。
// 署名不正・期限切れの場合はエラーを返す。
func (g *Generator) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(g.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
