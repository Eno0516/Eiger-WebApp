package repository

import (
	"context"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
)

type CircleMemberRepository interface {
	// Upsert はメンバーを追加または更新する（参加申請・ロール変更に使用）。
	Upsert(ctx context.Context, member *entity.CircleMember) error
	FindByCircleAndUser(ctx context.Context, circleID, userID string) (*entity.CircleMember, error)
	FindByCircle(ctx context.Context, circleID string) ([]*entity.CircleMember, error)
	FindByUser(ctx context.Context, userID string) ([]*entity.CircleMember, error)
	// FindAdminsByCircle はサークルの admin ユーザー一覧を返す（通知送信先の取得に使用）。
	FindAdminsByCircle(ctx context.Context, circleID string) ([]*entity.CircleMember, error)
	Delete(ctx context.Context, circleID, userID string) error
}
