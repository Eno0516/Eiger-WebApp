package repository

import (
	"context"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
)

type UserProfileRepository interface {
	// FindByUserID はユーザーIDからプロフィールを取得する。存在しない場合は apperrors.ErrNotFound を返す。
	FindByUserID(ctx context.Context, userID string) (*entity.UserProfile, error)
	// Create はプロフィールを新規作成する。同一 userID が存在する場合は apperrors.ErrConflict を返す。
	Create(ctx context.Context, profile *entity.UserProfile) (*entity.UserProfile, error)
	// Update はプロフィールを更新する。存在しない場合は apperrors.ErrNotFound を返す。
	Update(ctx context.Context, profile *entity.UserProfile) (*entity.UserProfile, error)
}
