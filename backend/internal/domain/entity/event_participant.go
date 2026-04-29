package entity

import (
	"errors"
	"time"
)

type ParticipantStatus string

const (
	ParticipantStatusApplied   ParticipantStatus = "applied"
	ParticipantStatusConfirmed ParticipantStatus = "confirmed"
	ParticipantStatusRejected  ParticipantStatus = "rejected"
)

type EventParticipant struct {
	EventID   string
	UserID    string
	Status    ParticipantStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

var ErrInvalidParticipantStatus = errors.New("invalid participant status")

func ValidateParticipantStatus(s ParticipantStatus) error {
	switch s {
	case ParticipantStatusApplied, ParticipantStatusConfirmed, ParticipantStatusRejected:
		return nil
	default:
		return ErrInvalidParticipantStatus
	}
}
