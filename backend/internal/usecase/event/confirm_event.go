package event

import (
	"context"
	"fmt"
	"time"

	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type ConfirmEventInput struct {
	EventID string
}

// ConfirmEventUseCase は参加者を確定し、confirmed メンバー全員に LINE 通知を送る。
// confirmed_at を記録し、以降はユーザー側で参加メンバー一覧・LINE グループ URL が表示される。
type ConfirmEventUseCase struct {
	eventRepo       domainrepo.EventRepository
	participantRepo domainrepo.EventParticipantRepository
	// TODO: LINE Messaging API クライアントを追加する
}

func NewConfirmEventUseCase(
	eventRepo domainrepo.EventRepository,
	participantRepo domainrepo.EventParticipantRepository,
) *ConfirmEventUseCase {
	return &ConfirmEventUseCase{eventRepo: eventRepo, participantRepo: participantRepo}
}

// Execute は confirmed_at を記録し、confirmed ユーザーに LINE グループ招待を通知する。
func (uc *ConfirmEventUseCase) Execute(ctx context.Context, input ConfirmEventInput) error {
	ev, err := uc.eventRepo.FindByID(ctx, input.EventID)
	if err != nil {
		return err
	}
	if ev.IsConfirmed() {
		return apperrors.NewConflict("event already confirmed", nil)
	}

	// confirmed_at を記録
	now := time.Now()
	ev.ConfirmedAt = &now
	if err := uc.eventRepo.Update(ctx, ev); err != nil {
		return fmt.Errorf("update event: %w", err)
	}

	// confirmed メンバーに LINE 通知
	// TODO: participantRepo.FindConfirmedByEvent で confirmed 一覧を取得する
	// TODO: LINE Messaging API で line_group_url を含む通知を送信する
	return nil
}
