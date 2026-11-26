// Package car は、車に関する usecase。
package car

import (
	"context"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/membership"
)

type ListInput struct {
	CircleID string
	MemberID string
}

type ListUseCase struct {
	cars    repository.CarRepository
	members membership.Directory
}

func NewList(cars repository.CarRepository, members membership.Directory) *ListUseCase {
	return &ListUseCase{cars: cars, members: members}
}

// Execute は、サークルの車の一覧を返す。見られるのは、そのサークルのメンバーのみ。
func (uc *ListUseCase) Execute(ctx context.Context, in ListInput) ([]entity.Car, error) {
	if err := membership.Require(ctx, uc.members, in.CircleID, in.MemberID); err != nil {
		return nil, err
	}
	return uc.cars.ListByCircle(ctx, in.CircleID)
}
