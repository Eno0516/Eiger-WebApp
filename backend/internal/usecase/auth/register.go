package auth

import (
	"context"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
)

type RegisterInput struct {
	Name     string
	Email    string
	Password string // プレーンテキスト。usecase 内でハッシュ化する
}

type RegisterOutput struct {
	User *entity.User
}

type RegisterUseCase struct {
	repo   domainrepo.UserRepository
	mailer EmailSender
}

func NewRegisterUseCase(repo domainrepo.UserRepository, mailer EmailSender) *RegisterUseCase {
	return &RegisterUseCase{repo: repo, mailer: mailer}
}

// Execute はユーザー登録を行い、メール認証メールを送信する。
// メールが既に使用中の場合は 409 エラーを返す。
func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	// TODO: input.Email の重複チェック（FindByEmail → 存在したら 409）
	// TODO: bcrypt でパスワードをハッシュ化する
	// TODO: crypto/rand で email_verification_token を生成する（128bit hex）
	// TODO: entity.User を組み立てる（Role=RolePending, EmailVerificationExpiresAt=Now+24h）
	// TODO: user.Validate() を呼ぶ
	// TODO: repo.Create でDBに保存する
	// TODO: mailer.SendVerificationEmail でメールを送信する
	// TODO: RegisterOutput を返す
	panic("not implemented")
}
