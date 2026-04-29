package entity

import (
	"errors"
	"strings"
	"time"
)

type Event struct {
	ID            string
	CircleID      string
	Title         string
	Date          time.Time
	Description   string
	FormID        *string  // nil = フォーム未設定
	LiffURL       string
	LineGroupURL  string
	ConfirmedAt   *time.Time // nil = 未確定
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

var (
	ErrEventTitleRequired = errors.New("event title is required")
	ErrEventDateRequired  = errors.New("event date is required")
)

func (e *Event) Validate() error {
	if strings.TrimSpace(e.Title) == "" {
		return ErrEventTitleRequired
	}
	if e.Date.IsZero() {
		return ErrEventDateRequired
	}
	return nil
}

func (e *Event) IsConfirmed() bool {
	return e.ConfirmedAt != nil
}
