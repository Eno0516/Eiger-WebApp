package email

import "context"

// Config はメール送信に必要な設定。
type Config struct {
	Host       string
	Port       int
	User       string
	Password   string
	From       string
	AppBaseURL string // メール内リンクのベース URL
}

// Sender は gomail を使った SMTP メール送信実装。
// usecase/auth.EmailSender interface を実装する。
type Sender struct {
	cfg Config
}

func NewSender(cfg Config) *Sender {
	return &Sender{cfg: cfg}
}

// SendVerificationEmail はメール認証リンクを送信する。
// リンク形式: {AppBaseURL}/verify-email?token={token}
func (s *Sender) SendVerificationEmail(ctx context.Context, toEmail, token string) error {
	// TODO: gopkg.in/gomail.v2 で本文を組み立てる
	// TODO: 認証リンク: s.cfg.AppBaseURL + "/verify-email?token=" + token
	// TODO: dialer.DialAndSend で送信する
	panic("not implemented")
}

// SendPasswordResetEmail はパスワードリセットリンクを送信する。
// リンク形式: {AppBaseURL}/reset-password?token={token}
func (s *Sender) SendPasswordResetEmail(ctx context.Context, toEmail, token string) error {
	// TODO: gopkg.in/gomail.v2 で本文を組み立てる
	// TODO: リセットリンク: s.cfg.AppBaseURL + "/reset-password?token=" + token
	// TODO: dialer.DialAndSend で送信する
	panic("not implemented")
}
