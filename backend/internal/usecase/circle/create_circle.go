package circle

import (
	"context"
	"fmt"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type CreateCircleInput struct {
	Name               string
	Description        string
	ApplicationFormURL string
	SpreadsheetID      string
}

// CreateCircleUseCase は system_admin がサークルを新規作成する。
type CreateCircleUseCase struct {
	circleRepo domainrepo.CircleRepository
}

func NewCreateCircleUseCase(circleRepo domainrepo.CircleRepository) *CreateCircleUseCase {
	return &CreateCircleUseCase{circleRepo: circleRepo}
}

// Execute は入力を検証してサークルを作成し、作成済みエンティティを返す。
func (uc *CreateCircleUseCase) Execute(ctx context.Context, input CreateCircleInput) (*entity.Circle, error) {
	circle := &entity.Circle{
		Name:               input.Name,
		Description:        input.Description,
		ApplicationFormURL: input.ApplicationFormURL,
		SpreadsheetID:      input.SpreadsheetID,
	}
	if err := circle.Validate(); err != nil {
		return nil, apperrors.NewBadRequest("validation failed", err)
	}
	if err := uc.circleRepo.Create(ctx, circle); err != nil {
		return nil, fmt.Errorf("create circle: %w", err)
	}
	return circle, nil
}
