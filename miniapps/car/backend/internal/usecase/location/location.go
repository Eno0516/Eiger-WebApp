// Package location は、車の位置の自己申告に関する usecase。
package location

import (
	"context"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/valueobject"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/membership"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/port"
)

type ReportInput struct {
	ReservationID string
	MemberID      string
	Link          string
	Note          string
}

type ReportUseCase struct {
	reservations repository.ReservationRepository
	locations    repository.LocationReportRepository
	members      membership.Directory
	clock        port.Clock
}

func NewReport(
	reservations repository.ReservationRepository,
	locations repository.LocationReportRepository,
	members membership.Directory,
	clock port.Clock,
) *ReportUseCase {
	return &ReportUseCase{reservations: reservations, locations: locations, members: members, clock: clock}
}

// Execute は、利用中の本人による位置の申告を保存する。
func (uc *ReportUseCase) Execute(ctx context.Context, in ReportInput) (entity.LocationReport, error) {
	res, err := uc.reservations.FindByID(ctx, in.ReservationID)
	if err != nil {
		return entity.LocationReport{}, err
	}
	if err := membership.Require(ctx, uc.members, res.CircleID, in.MemberID); err != nil {
		return entity.LocationReport{}, err
	}
	link, err := valueobject.NewLocationLink(in.Link)
	if err != nil {
		return entity.LocationReport{}, err
	}
	report, err := entity.NewLocationReport(res, in.MemberID, link, in.Note, uc.clock.Now())
	if err != nil {
		return entity.LocationReport{}, err
	}
	if err := uc.locations.Save(ctx, report); err != nil {
		return entity.LocationReport{}, err
	}
	return report, nil
}

type GetInput struct {
	CarID    string
	MemberID string
}

type GetUseCase struct {
	cars      repository.CarRepository
	locations repository.LocationReportRepository
	members   membership.Directory
}

func NewGet(cars repository.CarRepository, locations repository.LocationReportRepository, members membership.Directory) *GetUseCase {
	return &GetUseCase{cars: cars, locations: locations, members: members}
}

// Execute は、車の最新の申告を返す。見られるのは、その車のサークルのメンバーのみ。
func (uc *GetUseCase) Execute(ctx context.Context, in GetInput) (entity.LocationReport, error) {
	car, err := uc.cars.FindByID(ctx, in.CarID)
	if err != nil {
		return entity.LocationReport{}, err
	}
	if err := membership.Require(ctx, uc.members, car.CircleID, in.MemberID); err != nil {
		return entity.LocationReport{}, err
	}
	return uc.locations.LatestByCar(ctx, car.ID)
}
