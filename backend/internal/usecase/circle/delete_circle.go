package circle

import (
	"context"
	"fmt"

	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
)

type DeleteCircleInput struct {
	CircleID string
}

// DeleteCircleUseCase は system_admin がサークルを削除する。
// 関連する circle_members は DB の CASCADE で自動削除される。
type DeleteCircleUseCase struct {
	circleRepo domainrepo.CircleRepository
}

func NewDeleteCircleUseCase(circleRepo domainrepo.CircleRepository) *DeleteCircleUseCase {
	return &DeleteCircleUseCase{circleRepo: circleRepo}
}

// Execute は指定されたサークルを削除する。存在しない場合は 404 を返す。
func (uc *DeleteCircleUseCase) Execute(ctx context.Context, input DeleteCircleInput) error {
	if _, err := uc.circleRepo.FindByID(ctx, input.CircleID); err != nil {
		return err
	}
	if err := uc.circleRepo.Delete(ctx, input.CircleID); err != nil {
		return fmt.Errorf("delete circle: %w", err)
	}
	return nil
}
