package auth

import (
	"context"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
)

// EmailSender はメール送信の抽象。infrastructure/email で実装する。
type EmailSender interface {
	SendVerificationEmail(ctx context.Context, toEmail, token string) error
	SendPasswordResetEmail(ctx context.Context, toEmail, token string) error
}

// JWTGenerator は JWT 生成・検証の抽象。pkg/jwt で実装する。
type JWTGenerator interface {
	Generate(userID, email string, role entity.UserRole) (string, error)
}
