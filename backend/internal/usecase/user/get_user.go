package user

import (
	"context"
	"fmt"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
)

type GetUserOutput struct {
	User *entity.User
}

type GetUsersOutput struct {
	Users []*entity.User
}

type GetUserUseCase struct {
	repo domainrepo.UserRepository
}

func NewGetUserUseCase(repo domainrepo.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{repo: repo}
}

func (uc *GetUserUseCase) ExecuteByID(ctx context.Context, id string) (*GetUserOutput, error) {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	return &GetUserOutput{User: user}, nil
}

func (uc *GetUserUseCase) ExecuteAll(ctx context.Context) (*GetUsersOutput, error) {
	users, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("find all users: %w", err)
	}
	return &GetUsersOutput{Users: users}, nil
}
