package user

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type UpdateDisplayNameInput struct {
	UserID      string
	DisplayName string
}

type UpdateDisplayNameOutput struct {
	UserID      string
	DisplayName string
}

// UpdateDisplayNameUseCase は display_name を更新する。
// 初回ログイン時のニックネーム設定と、その後の変更に使用する。
type UpdateDisplayNameUseCase struct {
	userRepo domainrepo.UserRepository
	log      zerolog.Logger
}

func NewUpdateDisplayNameUseCase(userRepo domainrepo.UserRepository, log zerolog.Logger) *UpdateDisplayNameUseCase {
	return &UpdateDisplayNameUseCase{userRepo: userRepo, log: log}
}

// Execute は display_name を更新して更新後のユーザー情報を返す。
func (uc *UpdateDisplayNameUseCase) Execute(ctx context.Context, input UpdateDisplayNameInput) (*UpdateDisplayNameOutput, error) {
	u, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	u.DisplayName = input.DisplayName
	if err := u.Validate(); err != nil {
		return nil, apperrors.NewBadRequest("validation failed", err)
	}
	updated, err := uc.userRepo.Update(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return &UpdateDisplayNameOutput{UserID: updated.ID, DisplayName: updated.DisplayName}, nil
}
