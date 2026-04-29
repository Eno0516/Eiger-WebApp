package event

import (
	"context"
	"fmt"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type UpdateParticipantStatusInput struct {
	EventID  string
	UserID   string
	Status   entity.ParticipantStatus
}

// UpdateParticipantStatusUseCase は参加者の個別ステータスを更新する。
// 確定前・確定後どちらでも変更可能。
type UpdateParticipantStatusUseCase struct {
	participantRepo domainrepo.EventParticipantRepository
}

func NewUpdateParticipantStatusUseCase(repo domainrepo.EventParticipantRepository) *UpdateParticipantStatusUseCase {
	return &UpdateParticipantStatusUseCase{participantRepo: repo}
}

// Execute は対象ユーザーのステータスを更新する。
func (uc *UpdateParticipantStatusUseCase) Execute(ctx context.Context, input UpdateParticipantStatusInput) error {
	if err := entity.ValidateParticipantStatus(input.Status); err != nil {
		return apperrors.NewBadRequest("invalid status", err)
	}
	p := &entity.EventParticipant{
		EventID: input.EventID,
		UserID:  input.UserID,
		Status:  input.Status,
	}
	if err := uc.participantRepo.Upsert(ctx, p); err != nil {
		return fmt.Errorf("update participant status: %w", err)
	}
	return nil
}
