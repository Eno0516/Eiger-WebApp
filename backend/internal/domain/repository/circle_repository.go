package repository

import (
	"context"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
)

type CircleRepository interface {
	Create(ctx context.Context, circle *entity.Circle) error
	FindByID(ctx context.Context, id string) (*entity.Circle, error)
	FindAll(ctx context.Context) ([]*entity.Circle, error)
	Update(ctx context.Context, circle *entity.Circle) error
	Delete(ctx context.Context, id string) error
}
