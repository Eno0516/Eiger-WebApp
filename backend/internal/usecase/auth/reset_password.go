package auth

import (
	"context"

	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
)

type ResetPasswordInput struct {
	Token       string
	NewPassword string // プレーンテキスト。usecase 内でハッシュ化する
}

// ResetPasswordUseCase はトークンを検証してパスワードを更新する。
type ResetPasswordUseCase struct {
	repo domainrepo.UserRepository
}

func NewResetPasswordUseCase(repo domainrepo.UserRepository) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{repo: repo}
}

// Execute はリセットトークンを検証し、パスワードを更新してトークンを無効化する。
// トークンが無効・期限切れの場合は 400 を返す。
func (uc *ResetPasswordUseCase) Execute(ctx context.Context, input ResetPasswordInput) error {
	// TODO: repo.FindByResetToken でユーザーを取得する（存在しない/期限切れなら 400）
	// TODO: bcrypt で新しいパスワードをハッシュ化する
	// TODO: repo で password_hash を更新し、reset_token / expires_at を NULL にする
	panic("not implemented")
}
