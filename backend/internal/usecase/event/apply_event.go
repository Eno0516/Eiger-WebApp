package event

import (
	"context"
	"fmt"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type ApplyEventInput struct {
	EventID string
	UserID  string
}

// ApplyEventUseCase はフォーム回答後にユーザーが参加申請する（applied レコード作成）。
// フォーム送信とは別に呼び出す（案B: 2回叩く方式）。
type ApplyEventUseCase struct {
	participantRepo domainrepo.EventParticipantRepository
}

func NewApplyEventUseCase(repo domainrepo.EventParticipantRepository) *ApplyEventUseCase {
	return &ApplyEventUseCase{participantRepo: repo}
}

// Execute は参加申請レコードを作成する。既に申請済みの場合は 409 を返す。
func (uc *ApplyEventUseCase) Execute(ctx context.Context, input ApplyEventInput) error {
	existing, err := uc.participantRepo.FindByEventAndUser(ctx, input.EventID, input.UserID)
	if err != nil && !isNotFound(err) {
		return err
	}
	if existing != nil {
		return apperrors.NewConflict("already applied", nil)
	}
	p := &entity.EventParticipant{
		EventID: input.EventID,
		UserID:  input.UserID,
		Status:  entity.ParticipantStatusApplied,
	}
	if err := uc.participantRepo.Upsert(ctx, p); err != nil {
		return fmt.Errorf("apply event: %w", err)
	}
	return nil
}

func isNotFound(err error) bool {
	return apperrors.HTTPStatusCode(err) == 404
}
