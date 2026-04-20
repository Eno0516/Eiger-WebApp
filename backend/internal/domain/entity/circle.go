package entity

import (
	"errors"
	"strings"
	"time"
)

type Circle struct {
	ID                 string
	Name               string
	Description        string
	ApplicationFormURL string // LIFF の申し込みフォーム URL
	SpreadsheetID      string // Form 回答が蓄積される Spreadsheet の ID
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

var (
	ErrCircleNameRequired = errors.New("circle name is required")
)

func (c *Circle) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return ErrCircleNameRequired
	}
	return nil
}
