package dto

import (
	"time"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
)

type UserResponse struct {
	ID            string    `json:"id"`
	IsSystemAdmin bool      `json:"is_system_admin"`
	IsRegistered  bool      `json:"is_registered"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func UserFromEntity(u *entity.User) UserResponse {
	return UserResponse{
		ID:            u.ID,
		IsSystemAdmin: u.IsSystemAdmin,
		IsRegistered:  u.IsRegistered,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

// MeResponse は GET /users/me のレスポンス。
type MeResponse struct {
	ID            string `json:"id"`
	IsRegistered  bool   `json:"is_registered"`
	IsSystemAdmin bool   `json:"is_system_admin"`
}

// RegisterResponse は POST /auth/signup のレスポンス。
type RegisterResponse struct {
	DisplayName string `json:"display_name"`
}

type CircleResponse struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	ApplicationFormURL string    `json:"application_form_url"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func CircleFromEntity(c *entity.Circle) CircleResponse {
	return CircleResponse{
		ID:                 c.ID,
		Name:               c.Name,
		Description:        c.Description,
		ApplicationFormURL: c.ApplicationFormURL,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
	}
}

type CircleMemberResponse struct {
	CircleID string    `json:"circle_id"`
	UserID   string    `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

func CircleMemberFromEntity(m *entity.CircleMember) CircleMemberResponse {
	return CircleMemberResponse{
		CircleID: m.CircleID,
		UserID:   m.UserID,
		Role:     string(m.Role),
		JoinedAt: m.JoinedAt,
	}
}

type EventResponse struct {
	ID           string     `json:"id"`
	CircleID     string     `json:"circle_id"`
	Title        string     `json:"title"`
	Date         string     `json:"date"` // "YYYY-MM-DD"
	Description  string     `json:"description"`
	FormID       *string    `json:"form_id"`
	LiffURL      string     `json:"liff_url"`
	LineGroupURL string     `json:"line_group_url"`
	ConfirmedAt  *time.Time `json:"confirmed_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

func EventFromEntity(e *entity.Event) EventResponse {
	return EventResponse{
		ID:           e.ID,
		CircleID:     e.CircleID,
		Title:        e.Title,
		Date:         e.Date.Format("2006-01-02"),
		Description:  e.Description,
		FormID:       e.FormID,
		LiffURL:      e.LiffURL,
		LineGroupURL: e.LineGroupURL,
		ConfirmedAt:  e.ConfirmedAt,
		CreatedAt:    e.CreatedAt,
	}
}

type EventParticipantResponse struct {
	EventID  string `json:"event_id"`
	UserID   string `json:"user_id"`
	Status   string `json:"status"`
	JoinedAt string `json:"joined_at"`
}

func EventParticipantFromEntity(p *entity.EventParticipant) EventParticipantResponse {
	return EventParticipantResponse{
		EventID:  p.EventID,
		UserID:   p.UserID,
		Status:   string(p.Status),
		JoinedAt: p.CreatedAt.Format(time.RFC3339),
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	IsNewUser bool   `json:"is_new_user"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
