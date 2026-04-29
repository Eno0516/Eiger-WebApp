package auth

import (
	"context"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type UpdateCircleMemberRoleInput struct {
	CircleID     string
	TargetUserID string
	NewRole      entity.CircleRole
}

// UpdateCircleMemberRoleUseCase はサークル admin がメンバーのロールを変更する。
type UpdateCircleMemberRoleUseCase struct {
	circleMemberRepo domainrepo.CircleMemberRepository
}

func NewUpdateCircleMemberRoleUseCase(repo domainrepo.CircleMemberRepository) *UpdateCircleMemberRoleUseCase {
	return &UpdateCircleMemberRoleUseCase{circleMemberRepo: repo}
}

// Execute は対象ユーザーのサークル内ロールを更新する。
func (uc *UpdateCircleMemberRoleUseCase) Execute(ctx context.Context, input UpdateCircleMemberRoleInput) error {
	if err := entity.ValidateCircleRole(input.NewRole); err != nil {
		return apperrors.NewBadRequest("invalid circle role", err)
	}
	member, err := uc.circleMemberRepo.FindByCircleAndUser(ctx, input.CircleID, input.TargetUserID)
	if err != nil {
		return err
	}
	member.Role = input.NewRole
	return uc.circleMemberRepo.Upsert(ctx, member)
}
