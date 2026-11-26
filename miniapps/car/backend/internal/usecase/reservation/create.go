// Package reservation は、予約に関する usecase。
// ルール(重複禁止・本人のみ操作など)は domain が持ち、ここでは「取得 → ルールの呼び出し → 保存」の流れと、
// 外部の情報が要る判断(サークルのメンバーか)を担う。
package reservation

import (
	"context"
	"time"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/valueobject"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/membership"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/port"
)

type CreateInput struct {
	CarID    string
	MemberID string
	Start    time.Time
	End      time.Time
}

type CreateUseCase struct {
	cars         repository.CarRepository
	reservations repository.ReservationRepository
	members      membership.Directory
	clock        port.Clock
	ids          port.IDGenerator
}

func NewCreate(
	cars repository.CarRepository,
	reservations repository.ReservationRepository,
	members membership.Directory,
	clock port.Clock,
	ids port.IDGenerator,
) *CreateUseCase {
	return &CreateUseCase{cars: cars, reservations: reservations, members: members, clock: clock, ids: ids}
}

// Execute は車を予約する。事前チェックをすり抜けた同時予約は、Save が ErrReservationConflict で弾く。
func (uc *CreateUseCase) Execute(ctx context.Context, in CreateInput) (entity.Reservation, error) {
	car, err := uc.cars.FindByID(ctx, in.CarID)
	if err != nil {
		return entity.Reservation{}, err
	}
	if err := membership.Require(ctx, uc.members, car.CircleID, in.MemberID); err != nil {
		return entity.Reservation{}, err
	}
	period, err := valueobject.NewTimeRange(in.Start, in.End)
	if err != nil {
		return entity.Reservation{}, err
	}
	existing, err := uc.reservations.ListByCarInRange(ctx, car.ID, period)
	if err != nil {
		return entity.Reservation{}, err
	}
	res, err := entity.NewReservation(uc.ids.NewID(), car, in.MemberID, period, existing, uc.clock.Now())
	if err != nil {
		return entity.Reservation{}, err
	}
	if err := uc.reservations.Save(ctx, res); err != nil {
		return entity.Reservation{}, err
	}
	return res, nil
}
