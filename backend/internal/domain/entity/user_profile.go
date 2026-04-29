package entity

import (
	"errors"
	"strings"
	"time"
)

// UserProfile はアプリ内プロフィール。users テーブルと 1:1 対応。
// 初期登録（LIFFフォーム）完了時に作成される。
type UserProfile struct {
	UserID      string
	DisplayName string    // 部内あだ名（個人を特定できる本名はNG）
	University  string    // 広告ターゲティング用
	Faculty     string    // 広告ターゲティング用
	Grade       int       // 学年（1〜4）。留年時は本人が更新
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var (
	ErrUserProfileDisplayNameRequired = errors.New("display name is required")
	ErrUserProfileUserIDRequired      = errors.New("user id is required")
)

func (p *UserProfile) Validate() error {
	if strings.TrimSpace(p.UserID) == "" {
		return ErrUserProfileUserIDRequired
	}
	if strings.TrimSpace(p.DisplayName) == "" {
		return ErrUserProfileDisplayNameRequired
	}
	return nil
}
