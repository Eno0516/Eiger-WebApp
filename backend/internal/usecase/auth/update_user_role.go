package auth

import (
	"context"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
)

type UpdateUserRoleInput struct {
	TargetUserID string
	NewRole      entity.UserRole
}

// UpdateUserRoleUseCase は admin がユーザーのロールを変更する。
type UpdateUserRoleUseCase struct {
	repo domainrepo.UserRepository
}

func NewUpdateUserRoleUseCase(repo domainrepo.UserRepository) *UpdateUserRoleUseCase {
	return &UpdateUserRoleUseCase{repo: repo}
}

// Execute は対象ユーザーのロールを更新する。
// 対象ユーザーが存在しない場合は 404 を返す。
// NewRole が不正な値の場合は 400 を返す。
func (uc *UpdateUserRoleUseCase) Execute(ctx context.Context, input UpdateUserRoleInput) error {
	// TODO: input.NewRole が RoleAdmin / RoleMember / RolePending のいずれかか検証する（不正なら 400）
	// TODO: repo.FindByID で対象ユーザーの存在確認（存在しなければ 404）
	// TODO: repo.UpdateRole でロールを更新する
	panic("not implemented")
}
