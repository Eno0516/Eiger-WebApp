package reservation

import (
	"context"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/membership"
)

// transition は、予約を取得し、メンバーであることを確認し、ルール(apply)を適用して保存する、共通の流れ。
func transition(
	ctx context.Context,
	reservations repository.ReservationRepository,
	members membership.Directory,
	reservationID, memberID string,
	apply func(*entity.Reservation) error,
) (entity.Reservation, error) {
	res, err := reservations.FindByID(ctx, reservationID)
	if err != nil {
		return entity.Reservation{}, err
	}
	if err := membership.Require(ctx, members, res.CircleID, memberID); err != nil {
		return entity.Reservation{}, err
	}
	if err := apply(&res); err != nil {
		return entity.Reservation{}, err
	}
	if err := reservations.Save(ctx, res); err != nil {
		return entity.Reservation{}, err
	}
	return res, nil
}

type CancelInput struct {
	ReservationID string
	MemberID      string
}

type CancelUseCase struct {
	reservations repository.ReservationRepository
	members      membership.Directory
}

func NewCancel(reservations repository.ReservationRepository, members membership.Directory) *CancelUseCase {
	return &CancelUseCase{reservations: reservations, members: members}
}

func (uc *CancelUseCase) Execute(ctx context.Context, in CancelInput) (entity.Reservation, error) {
	return transition(ctx, uc.reservations, uc.members, in.ReservationID, in.MemberID,
		func(r *entity.Reservation) error { return r.Cancel(in.MemberID) })
}

type StartInput struct {
	ReservationID string
	MemberID      string
}

type StartUseCase struct {
	reservations repository.ReservationRepository
	members      membership.Directory
}

func NewStart(reservations repository.ReservationRepository, members membership.Directory) *StartUseCase {
	return &StartUseCase{reservations: reservations, members: members}
}

func (uc *StartUseCase) Execute(ctx context.Context, in StartInput) (entity.Reservation, error) {
	return transition(ctx, uc.reservations, uc.members, in.ReservationID, in.MemberID,
		func(r *entity.Reservation) error { return r.Start(in.MemberID) })
}
