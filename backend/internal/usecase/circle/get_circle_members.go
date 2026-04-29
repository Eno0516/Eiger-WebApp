package circle

import (
	"context"
	"fmt"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
)

type GetCircleMembersInput struct {
	CircleID string
}

// GetCircleMembersUseCase はサークルのメンバー一覧を取得する。
// system_admin・circle_admin どちらからも呼び出される。
type GetCircleMembersUseCase struct {
	circleMemberRepo domainrepo.CircleMemberRepository
}

func NewGetCircleMembersUseCase(repo domainrepo.CircleMemberRepository) *GetCircleMembersUseCase {
	return &GetCircleMembersUseCase{circleMemberRepo: repo}
}

// Execute は指定サークルの全メンバーを返す。
func (uc *GetCircleMembersUseCase) Execute(ctx context.Context, input GetCircleMembersInput) ([]*entity.CircleMember, error) {
	members, err := uc.circleMemberRepo.FindByCircle(ctx, input.CircleID)
	if err != nil {
		return nil, fmt.Errorf("find members: %w", err)
	}
	return members, nil
}
