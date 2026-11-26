// Package repository は、永続化の interface を定義する。実装は internal/repository(memory / postgres)に置く。
// 見つからない場合は、すべて domain.ErrNotFound を返す。
package repository

import (
	"context"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/valueobject"
)

type CarRepository interface {
	FindByID(ctx context.Context, id string) (entity.Car, error)
	ListByCircle(ctx context.Context, circleID string) ([]entity.Car, error)
	Save(ctx context.Context, car entity.Car) error
}

type ReservationRepository interface {
	// Save は予約を保存する。同時に別の予約が入って時間帯が重なった場合は、
	// 事前チェックをすり抜けても、実装側(DB の排他制約など)で domain.ErrReservationConflict を返すこと。
	Save(ctx context.Context, r entity.Reservation) error
	FindByID(ctx context.Context, id string) (entity.Reservation, error)
	// ListByCarInRange は、車の予約のうち、期間と重なるものを返す(予約の重複チェック用)。
	ListByCarInRange(ctx context.Context, carID string, period valueobject.TimeRange) ([]entity.Reservation, error)
	// ListByCircleInRange は、サークルの全車の予約のうち、期間と重なるものを返す(カレンダー表示用)。
	ListByCircleInRange(ctx context.Context, circleID string, period valueobject.TimeRange) ([]entity.Reservation, error)
	ListByMember(ctx context.Context, circleID, memberID string) ([]entity.Reservation, error)
}

type LocationReportRepository interface {
	Save(ctx context.Context, r entity.LocationReport) error
	// LatestByCar は、車の最新の申告を返す。
	LatestByCar(ctx context.Context, carID string) (entity.LocationReport, error)
}
