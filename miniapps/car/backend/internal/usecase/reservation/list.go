package reservation

import (
	"context"
	"time"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/valueobject"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/membership"
)

type ListByRangeInput struct {
	CircleID string
	MemberID string
	From     time.Time
	To       time.Time
}

// ListByRangeUseCase は、カレンダー表示用に、期間と重なる予約を返す。キャンセル済みは含めない。
type ListByRangeUseCase struct {
	reservations repository.ReservationRepository
	members      membership.Directory
}

func NewListByRange(reservations repository.ReservationRepository, members membership.Directory) *ListByRangeUseCase {
	return &ListByRangeUseCase{reservations: reservations, members: members}
}

func (uc *ListByRangeUseCase) Execute(ctx context.Context, in ListByRangeInput) ([]entity.Reservation, error) {
	if err := membership.Require(ctx, uc.members, in.CircleID, in.MemberID); err != nil {
		return nil, err
	}
	period, err := valueobject.NewTimeRange(in.From, in.To)
	if err != nil {
		return nil, err
	}
	all, err := uc.reservations.ListByCircleInRange(ctx, in.CircleID, period)
	if err != nil {
		return nil, err
	}
	out := make([]entity.Reservation, 0, len(all))
	for _, r := range all {
		if r.Status != entity.StatusCancelled {
			out = append(out, r)
		}
	}
	return out, nil
}

type ListMineInput struct {
	CircleID string
	MemberID string
}

type ListMineUseCase struct {
	reservations repository.ReservationRepository
	members      membership.Directory
}

func NewListMine(reservations repository.ReservationRepository, members membership.Directory) *ListMineUseCase {
	return &ListMineUseCase{reservations: reservations, members: members}
}

func (uc *ListMineUseCase) Execute(ctx context.Context, in ListMineInput) ([]entity.Reservation, error) {
	if err := membership.Require(ctx, uc.members, in.CircleID, in.MemberID); err != nil {
		return nil, err
	}
	return uc.reservations.ListByMember(ctx, in.CircleID, in.MemberID)
}
