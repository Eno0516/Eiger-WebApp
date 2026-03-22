package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type CreateUserInput struct {
	Name  string
	Email string
}

type CreateUserOutput struct {
	User *entity.User
}

type CreateUserUseCase struct {
	repo domainrepo.UserRepository
}

func NewCreateUserUseCase(repo domainrepo.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	user := &entity.User{
		Name:  input.Name,
		Email: input.Email,
	}

	if err := user.Validate(); err != nil {
		return nil, apperrors.Wrap(400, "validation failed", err)
	}

	existing, err := uc.repo.FindByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, fmt.Errorf("check existing user: %w", err)
	}
	if existing != nil {
		return nil, apperrors.Wrap(409, "email already in use", nil)
	}

	if err := uc.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &CreateUserOutput{User: user}, nil
}
