package event

import (
	"context"
	"fmt"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type BulkParticipantUpdate struct {
	UserID string
	Status entity.ParticipantStatus
}

type BulkUpdateParticipantsInput struct {
	EventID     string
	Participant []BulkParticipantUpdate
}

// BulkUpdateParticipantsUseCase はランダム選出結果など一括でステータスを更新する。
// フロントがランダム選出した結果をそのまま受け取って保存する。
type BulkUpdateParticipantsUseCase struct {
	participantRepo domainrepo.EventParticipantRepository
}

func NewBulkUpdateParticipantsUseCase(repo domainrepo.EventParticipantRepository) *BulkUpdateParticipantsUseCase {
	return &BulkUpdateParticipantsUseCase{participantRepo: repo}
}

// Execute は渡されたユーザー全員のステータスを一括更新する。
func (uc *BulkUpdateParticipantsUseCase) Execute(ctx context.Context, input BulkUpdateParticipantsInput) error {
	participants := make([]*entity.EventParticipant, 0, len(input.Participant))
	for _, u := range input.Participant {
		if err := entity.ValidateParticipantStatus(u.Status); err != nil {
			return apperrors.NewBadRequest("invalid status", err)
		}
		participants = append(participants, &entity.EventParticipant{
			EventID: input.EventID,
			UserID:  u.UserID,
			Status:  u.Status,
		})
	}
	if err := uc.participantRepo.BulkUpsert(ctx, participants); err != nil {
		return fmt.Errorf("bulk update participants: %w", err)
	}
	return nil
}
