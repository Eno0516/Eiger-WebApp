package repository

import (
	"context"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
)

type EventParticipantRepository interface {
	// Upsert は参加申請・ステータス更新に使用する。
	Upsert(ctx context.Context, p *entity.EventParticipant) error
	// BulkUpsert は一括ステータス更新（ランダム選出確定など）に使用する。
	BulkUpsert(ctx context.Context, participants []*entity.EventParticipant) error
	FindByEvent(ctx context.Context, eventID string) ([]*entity.EventParticipant, error)
	FindByEventAndUser(ctx context.Context, eventID, userID string) (*entity.EventParticipant, error)
	// FindConfirmedByEvent は confirmed のみ返す（LINE通知送信先の取得に使用）。
	FindConfirmedByEvent(ctx context.Context, eventID string) ([]*entity.EventParticipant, error)
}
