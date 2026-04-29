package repository

import (
	"context"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
)

type UserRepository interface {
	// FindByID は UUID でユーザーを取得する。存在しない場合は apperrors.ErrNotFound を返す。
	FindByID(ctx context.Context, id string) (*entity.User, error)
	// FindByLineUserID は LINE Login の sub クレームでユーザーを取得する。存在しない場合は apperrors.ErrNotFound を返す。
	FindByLineUserID(ctx context.Context, lineUserID string) (*entity.User, error)
	// FindAll は全ユーザーを取得する。
	FindAll(ctx context.Context) ([]*entity.User, error)
	// Create はユーザーを新規作成する。LineUserID が重複する場合は apperrors.ErrConflict を返す。
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	// Update はユーザー情報を更新する。
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	// Delete はユーザーを削除する。存在しない場合は apperrors.ErrNotFound を返す。
	Delete(ctx context.Context, id string) error
}
