package reservation

import (
	"context"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/valueobject"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/membership"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/port"
)

type ReturnInput struct {
	ReservationID string
	MemberID      string
	DistanceKm    int
	Note          string
	// ParkingLink が空でなければ、返却の直前に、駐車場所として位置を申告する。
	ParkingLink string
	ParkingNote string
}

type ReturnUseCase struct {
	reservations repository.ReservationRepository
	locations    repository.LocationReportRepository
	members      membership.Directory
	clock        port.Clock
}

func NewReturn(
	reservations repository.ReservationRepository,
	locations repository.LocationReportRepository,
	members membership.Directory,
	clock port.Clock,
) *ReturnUseCase {
	return &ReturnUseCase{reservations: reservations, locations: locations, members: members, clock: clock}
}

// Execute は返却する。何かを保存する前に、すべての入力を検証する
// (リンク不正・走行距離不正のときに、位置だけが申告されて返却されない、という状態を作らない)。
// TODO: 位置の保存と返却の保存は、現状は別々。DB 実装(#9)でトランザクションにまとめる。
func (uc *ReturnUseCase) Execute(ctx context.Context, in ReturnInput) (entity.Reservation, error) {
	res, err := uc.reservations.FindByID(ctx, in.ReservationID)
	if err != nil {
		return entity.Reservation{}, err
	}
	if err := membership.Require(ctx, uc.members, res.CircleID, in.MemberID); err != nil {
		return entity.Reservation{}, err
	}

	var parking *entity.LocationReport
	if in.ParkingLink != "" {
		link, err := valueobject.NewLocationLink(in.ParkingLink)
		if err != nil {
			return entity.Reservation{}, err
		}
		// 位置の申告は利用中にしかできないため、返却の前の予約(res)で作る。
		report, err := entity.NewLocationReport(res, in.MemberID, link, in.ParkingNote, uc.clock.Now())
		if err != nil {
			return entity.Reservation{}, err
		}
		parking = &report
	}

	returned := res
	if err := returned.Return(in.MemberID, in.DistanceKm, in.Note, uc.clock.Now()); err != nil {
		return entity.Reservation{}, err
	}

	if parking != nil {
		if err := uc.locations.Save(ctx, *parking); err != nil {
			return entity.Reservation{}, err
		}
	}
	if err := uc.reservations.Save(ctx, returned); err != nil {
		return entity.Reservation{}, err
	}
	return returned, nil
}
