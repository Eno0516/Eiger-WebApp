package auth

import "context"

// LINEIDTokenVerifier は LINE ID トークンの検証抽象。
// pkg/lineauth で実装する。
type LINEIDTokenVerifier interface {
	// Verify は ID トークンを検証して LineUserID（sub クレーム）を返す。
	// トークンが無効の場合はエラーを返す。
	Verify(ctx context.Context, idToken string) (lineUserID string, err error)
}

// JWTGenerator は JWT 生成の抽象。
// pkg/jwt で実装する。
type JWTGenerator interface {
	Generate(userID string, isSystemAdmin bool) (string, error)
}

