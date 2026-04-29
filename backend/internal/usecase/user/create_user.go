package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

// CreateUserUseCase はシステム管理者が手動でユーザーを作成するユースケース。
// 通常のユーザー作成は LineAuthUseCase（初回 LINE ログイン時）が担う。
type CreateUserUseCase struct {
	repo        domainrepo.UserRepository
	profileRepo domainrepo.UserProfileRepository
}

type CreateUserInput struct {
	LineUserID  string
	DisplayName string
	University  string
	Faculty     string
	Grade       int
}

type CreateUserOutput struct {
	User *entity.User
}

func NewCreateUserUseCase(repo domainrepo.UserRepository, profileRepo domainrepo.UserProfileRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo, profileRepo: profileRepo}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	u := &entity.User{
		LineUserID: input.LineUserID,
	}
	if err := u.Validate(); err != nil {
		return nil, apperrors.NewBadRequest("validation failed", err)
	}
	created, err := uc.repo.Create(ctx, u)
	if err != nil {
		if errors.Is(err, apperrors.ErrConflict) {
			return nil, apperrors.NewConflict("line_user_id already in use", err)
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	profile := &entity.UserProfile{
		UserID:      created.ID,
		DisplayName: input.DisplayName,
		University:  input.University,
		Faculty:     input.Faculty,
		Grade:       input.Grade,
	}
	if _, err := uc.profileRepo.Create(ctx, profile); err != nil {
		return nil, fmt.Errorf("create user profile: %w", err)
	}
	return &CreateUserOutput{User: created}, nil
}
