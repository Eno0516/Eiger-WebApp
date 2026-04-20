package entity

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID            string
	LineUserID    string // LINE Login の sub クレーム（必須）
	DisplayName   string // 表示名（ニックネーム）
	IsSystemAdmin bool   // サークル作成権限を持つスーパー管理者
	IsRegistered  bool   // 初期登録（LIFFフォーム）完了フラグ
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

var (
	ErrUserLineUserIDRequired = errors.New("line user id is required")
)

func (u *User) Validate() error {
	if strings.TrimSpace(u.LineUserID) == "" {
		return ErrUserLineUserIDRequired
	}
	return nil
}
