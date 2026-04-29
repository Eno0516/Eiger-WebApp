package circle

import (
	"context"
	"errors"
	"fmt"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type ApplyCircleInput struct {
	CircleID string
	UserID   string
}

// ApplyCircleUseCase はユーザーがサークルに参加申請する。
// visitor レコードを作成し、サークルの admin 全員にアプリ内通知を送る。
type ApplyCircleUseCase struct {
	circleRepo       domainrepo.CircleRepository
	circleMemberRepo domainrepo.CircleMemberRepository
}

func NewApplyCircleUseCase(
	circleRepo domainrepo.CircleRepository,
	circleMemberRepo domainrepo.CircleMemberRepository,
) *ApplyCircleUseCase {
	return &ApplyCircleUseCase{circleRepo: circleRepo, circleMemberRepo: circleMemberRepo}
}

// Execute は指定サークルに visitor としてメンバーを追加し、admin に通知する。
func (uc *ApplyCircleUseCase) Execute(ctx context.Context, input ApplyCircleInput) error {
	if _, err := uc.circleRepo.FindByID(ctx, input.CircleID); err != nil {
		return err
	}
	existing, err := uc.circleMemberRepo.FindByCircleAndUser(ctx, input.CircleID, input.UserID)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return err
	}
	if existing != nil {
		return apperrors.NewConflict("already applied or member", nil)
	}
	member := &entity.CircleMember{
		CircleID: input.CircleID,
		UserID:   input.UserID,
		Role:     entity.CircleRoleVisitor,
	}
	if err := uc.circleMemberRepo.Upsert(ctx, member); err != nil {
		return fmt.Errorf("upsert member: %w", err)
	}
	// NOTE: admin への通知は別途実装予定
	return nil
}

// import guard
var _ = entity.CircleRoleVisitor
