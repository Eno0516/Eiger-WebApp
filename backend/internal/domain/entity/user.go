package entity

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
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
