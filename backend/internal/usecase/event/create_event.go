package event

import (
	"context"
	"fmt"
	"time"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type CreateEventInput struct {
	CircleID     string
	Title        string
	Date         time.Time
	Description  string
	FormID       *string
	LiffURL      string
	LineGroupURL string
}

// CreateEventUseCase は circle_admin がイベントを作成する。
type CreateEventUseCase struct {
	eventRepo domainrepo.EventRepository
}

func NewCreateEventUseCase(eventRepo domainrepo.EventRepository) *CreateEventUseCase {
	return &CreateEventUseCase{eventRepo: eventRepo}
}

// Execute は入力を検証してイベントを作成し、作成済みエンティティを返す。
func (uc *CreateEventUseCase) Execute(ctx context.Context, input CreateEventInput) (*entity.Event, error) {
	ev := &entity.Event{
		CircleID:     input.CircleID,
		Title:        input.Title,
		Date:         input.Date,
		Description:  input.Description,
		FormID:       input.FormID,
		LiffURL:      input.LiffURL,
		LineGroupURL: input.LineGroupURL,
	}
	if err := ev.Validate(); err != nil {
		return nil, apperrors.NewBadRequest("validation failed", err)
	}
	if err := uc.eventRepo.Create(ctx, ev); err != nil {
		return nil, fmt.Errorf("create event: %w", err)
	}
	return ev, nil
}
