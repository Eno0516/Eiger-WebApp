package event

import (
	"context"
	"fmt"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
)

type ListEventsInput struct {
	CircleID string
}

// ListEventsUseCase は指定サークルのイベント一覧を取得する。
type ListEventsUseCase struct {
	eventRepo domainrepo.EventRepository
}

func NewListEventsUseCase(eventRepo domainrepo.EventRepository) *ListEventsUseCase {
	return &ListEventsUseCase{eventRepo: eventRepo}
}

// Execute は指定サークルの全イベントを返す。
func (uc *ListEventsUseCase) Execute(ctx context.Context, input ListEventsInput) ([]*entity.Event, error) {
	events, err := uc.eventRepo.FindByCircle(ctx, input.CircleID)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	return events, nil
}
