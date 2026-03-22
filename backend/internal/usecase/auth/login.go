package auth

import (
	"context"

	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
)

type LoginInput struct {
	Email    string
	Password string // プレーンテキスト
}

type LoginOutput struct {
	Token string // JWT アクセストークン
}

type LoginUseCase struct {
	repo domainrepo.UserRepository
	jwt  JWTGenerator
}

func NewLoginUseCase(repo domainrepo.UserRepository, jwt JWTGenerator) *LoginUseCase {
	return &LoginUseCase{repo: repo, jwt: jwt}
}

// Execute はメール・パスワードを検証し JWT を返す。
// メール未認証の場合は 403、認証情報が不正の場合は 401 を返す。
func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	// TODO: repo.FindByEmail でユーザーを取得する（存在しない場合は 401）
	// TODO: user.IsEmailVerified() をチェックする（false なら 403）
	// TODO: bcrypt.CompareHashAndPassword でパスワードを検証する（不一致なら 401）
	// TODO: jwt.Generate でトークンを生成して返す
	panic("not implemented")
}
