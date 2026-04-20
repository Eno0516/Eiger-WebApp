package entity

import (
	"errors"
	"time"
)

type CircleRole string

const (
	CircleRoleAdmin   CircleRole = "admin"
	CircleRoleMember  CircleRole = "member"
	CircleRoleVisitor CircleRole = "visitor"
)

type CircleMember struct {
	CircleID string
	UserID   string
	Role     CircleRole
	JoinedAt time.Time
}

var ErrInvalidCircleRole = errors.New("invalid circle role")

func ValidateCircleRole(role CircleRole) error {
	switch role {
	case CircleRoleAdmin, CircleRoleMember, CircleRoleVisitor:
		return nil
	default:
		return ErrInvalidCircleRole
	}
}
