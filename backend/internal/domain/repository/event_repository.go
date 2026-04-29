package repository

import (
	"context"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
)

type EventRepository interface {
	Create(ctx context.Context, event *entity.Event) error
	FindByID(ctx context.Context, id string) (*entity.Event, error)
	FindByCircle(ctx context.Context, circleID string) ([]*entity.Event, error)
	Update(ctx context.Context, event *entity.Event) error
	Delete(ctx context.Context, id string) error
}
