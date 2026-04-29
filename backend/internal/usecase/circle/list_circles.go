package circle

import (
	"context"
	"fmt"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
)

type ListCirclesInput struct {
	RequesterUserID string
	IsSystemAdmin   bool
}

// ListCirclesUseCase はサークル一覧を取得する。
// system_admin は全件、circle_admin は自分が所属するサークルのみ返す。
type ListCirclesUseCase struct {
	circleRepo       domainrepo.CircleRepository
	circleMemberRepo domainrepo.CircleMemberRepository
}

func NewListCirclesUseCase(
	circleRepo domainrepo.CircleRepository,
	circleMemberRepo domainrepo.CircleMemberRepository,
) *ListCirclesUseCase {
	return &ListCirclesUseCase{circleRepo: circleRepo, circleMemberRepo: circleMemberRepo}
}

// Execute は呼び出し元のロールに応じてサークル一覧を返す。
func (uc *ListCirclesUseCase) Execute(ctx context.Context, input ListCirclesInput) ([]*entity.Circle, error) {
	if input.IsSystemAdmin {
		circles, err := uc.circleRepo.FindAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("find all circles: %w", err)
		}
		return circles, nil
	}

	// circle_admin: 自分が所属するサークルのみ返す
	members, err := uc.circleMemberRepo.FindByUser(ctx, input.RequesterUserID)
	if err != nil {
		return nil, fmt.Errorf("find circles by user: %w", err)
	}
	circles := make([]*entity.Circle, 0, len(members))
	for _, m := range members {
		if m.Role != entity.CircleRoleAdmin {
			continue
		}
		c, err := uc.circleRepo.FindByID(ctx, m.CircleID)
		if err != nil {
			return nil, fmt.Errorf("find circle %s: %w", m.CircleID, err)
		}
		circles = append(circles, c)
	}
	return circles, nil
}
