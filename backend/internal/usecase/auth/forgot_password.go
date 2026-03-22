package auth

import (
	"context"

	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
)

type ForgotPasswordInput struct {
	Email string
}

// ForgotPasswordUseCase はパスワードリセットメールを送信する。
type ForgotPasswordUseCase struct {
	repo   domainrepo.UserRepository
	mailer EmailSender
}

func NewForgotPasswordUseCase(repo domainrepo.UserRepository, mailer EmailSender) *ForgotPasswordUseCase {
	return &ForgotPasswordUseCase{repo: repo, mailer: mailer}
}

// Execute はリセットトークンを生成してメールを送信する。
// ユーザーが存在しない場合も 200 を返す（ユーザー存在確認を防ぐ）。
func (uc *ForgotPasswordUseCase) Execute(ctx context.Context, input ForgotPasswordInput) error {
	// TODO: repo.FindByEmail でユーザーを取得する
	// TODO: ユーザーが存在しない場合（ErrNotFound）は何もせず return nil（セキュリティ上）
	// TODO: crypto/rand でリセットトークンを生成する（128bit hex）
	// TODO: repo で password_reset_token / expires_at（Now+24h）を更新する
	// TODO: mailer.SendPasswordResetEmail でメールを送信する
	panic("not implemented")
}
