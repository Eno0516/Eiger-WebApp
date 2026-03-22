package auth

import (
	"context"

	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
)

type VerifyEmailInput struct {
	Token string
}

type ResendVerificationInput struct {
	Email string
}

// VerifyEmailUseCase はメール内リンクのトークンを検証して認証完了にする。
type VerifyEmailUseCase struct {
	repo domainrepo.UserRepository
}

func NewVerifyEmailUseCase(repo domainrepo.UserRepository) *VerifyEmailUseCase {
	return &VerifyEmailUseCase{repo: repo}
}

// Execute はトークンを検証し email_verified_at を記録する。
// トークンが無効・期限切れの場合は 400 を返す。
func (uc *VerifyEmailUseCase) Execute(ctx context.Context, input VerifyEmailInput) error {
	// TODO: repo.FindByVerificationToken でユーザーを取得する（存在しない/期限切れなら 400）
	// TODO: すでに認証済み（IsEmailVerified）なら何もせず return nil
	// TODO: repo の UpdateEmailVerified に相当する更新を行う
	//       （email_verified_at=NOW(), token/expires=NULL）
	panic("not implemented")
}

// ResendVerificationUseCase は認証メールを再送信する。
type ResendVerificationUseCase struct {
	repo   domainrepo.UserRepository
	mailer EmailSender
}

func NewResendVerificationUseCase(repo domainrepo.UserRepository, mailer EmailSender) *ResendVerificationUseCase {
	return &ResendVerificationUseCase{repo: repo, mailer: mailer}
}

// Execute はトークンを再生成してメールを再送する。
// 既に認証済みの場合は 400、1分以内の再送は 429 を返す。
func (uc *ResendVerificationUseCase) Execute(ctx context.Context, input ResendVerificationInput) error {
	// TODO: repo.FindByEmail でユーザーを取得する（存在しない場合は 404）
	// TODO: IsEmailVerified() == true なら 400（既に認証済み）
	// TODO: EmailVerificationExpiresAt が Now+23h より未来なら 429（1分レート制限）
	// TODO: 新しいトークンを生成し repo で token/expires を更新する
	// TODO: mailer.SendVerificationEmail でメールを再送する
	panic("not implemented")
}
