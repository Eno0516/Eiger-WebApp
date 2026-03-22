package entity

import (
	"errors"
	"strings"
	"time"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleMember  UserRole = "member"
	RolePending UserRole = "pending"
)

type User struct {
	ID                         string
	Name                       string
	Email                      string
	PasswordHash               string     // bcrypt ハッシュ
	Role                       UserRole   // pending / member / admin
	EmailVerifiedAt            *time.Time // nil = 未認証
	EmailVerificationToken     *string
	EmailVerificationExpiresAt *time.Time
	PasswordResetToken         *string
	PasswordResetExpiresAt     *time.Time
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

var (
	ErrUserNameRequired  = errors.New("user name is required")
	ErrUserEmailRequired = errors.New("user email is required")
	ErrUserEmailInvalid  = errors.New("user email is invalid")
)

func (u *User) Validate() error {
	if strings.TrimSpace(u.Name) == "" {
		return ErrUserNameRequired
	}
	if strings.TrimSpace(u.Email) == "" {
		return ErrUserEmailRequired
	}
	if !strings.Contains(u.Email, "@") {
		return ErrUserEmailInvalid
	}
	return nil
}

// IsEmailVerified はメール認証が完了しているかを返す。
func (u *User) IsEmailVerified() bool {
	return u.EmailVerifiedAt != nil
}

// IsPending は admin 承認待ち状態かを返す。
// pending ≠ 未ログインユーザー（登録済みだが未承認のユーザー）。
func (u *User) IsPending() bool {
	return u.Role == RolePending
}
