package repository

import (
	"context"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByVerificationToken(ctx context.Context, token string) (*entity.User, error)
	FindByResetToken(ctx context.Context, token string) (*entity.User, error)
	FindAll(ctx context.Context) ([]*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	UpdateRole(ctx context.Context, id string, role entity.UserRole) error
	Delete(ctx context.Context, id string) error
}
