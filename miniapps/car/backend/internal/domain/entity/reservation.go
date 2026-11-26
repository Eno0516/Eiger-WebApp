package entity

import (
	"time"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/valueobject"
)

type ReservationStatus string

const (
	StatusReserved  ReservationStatus = "reserved"  // 予約済み(まだ借りていない)
	StatusInUse     ReservationStatus = "in_use"    // 利用中
	StatusReturned  ReservationStatus = "returned"  // 返却済み
	StatusCancelled ReservationStatus = "cancelled" // キャンセル済み
)

// occupiesCar は、その予約が車の時間帯を押さえているかを返す。
// キャンセル済み・返却済みの予約は、時間帯を押さえない。
func (s ReservationStatus) occupiesCar() bool {
	return s == StatusReserved || s == StatusInUse
}

// Reservation は、メンバーによる車の予約。
// 状態遷移: reserved → in_use → returned、または reserved → cancelled。
type Reservation struct {
	ID       string
	CarID    string
	CircleID string
	MemberID string
	Period   valueobject.TimeRange
	Status   ReservationStatus

	// 返却時に記録される。
	DistanceKm int
	ReturnNote string
	ReturnedAt time.Time
}

// NewReservation は予約を作る。existing には、同じ車の既存の予約を渡す。
// now は現在時刻(ドメインは時計を持たない)。
func NewReservation(
	id string,
	car Car,
	memberID string,
	period valueobject.TimeRange,
	existing []Reservation,
	now time.Time,
) (Reservation, error) {
	if !car.Available {
		return Reservation{}, domain.ErrCarUnavailable
	}
	if period.Start().Before(now) {
		return Reservation{}, domain.ErrPastPeriod
	}
	for _, e := range existing {
		if e.CarID == car.ID && e.Status.occupiesCar() && e.Period.Overlaps(period) {
			return Reservation{}, domain.ErrReservationConflict
		}
	}
	return Reservation{
		ID:       id,
		CarID:    car.ID,
		CircleID: car.CircleID,
		MemberID: memberID,
		Period:   period,
		Status:   StatusReserved,
	}, nil
}

// Cancel は予約をキャンセルする。予約した本人のみ、予約済みの状態でのみ可能。
func (r *Reservation) Cancel(by string) error {
	if err := r.requireOwner(by); err != nil {
		return err
	}
	if r.Status != StatusReserved {
		return domain.ErrInvalidStatus
	}
	r.Status = StatusCancelled
	return nil
}

// Start は利用を開始する。予約した本人のみ、予約済みの状態でのみ可能。
func (r *Reservation) Start(by string) error {
	if err := r.requireOwner(by); err != nil {
		return err
	}
	if r.Status != StatusReserved {
		return domain.ErrInvalidStatus
	}
	r.Status = StatusInUse
	return nil
}

// Return は返却する。走行距離とメモを記録する。予約した本人のみ、利用中の状態でのみ可能。
func (r *Reservation) Return(by string, distanceKm int, note string, now time.Time) error {
	if err := r.requireOwner(by); err != nil {
		return err
	}
	if r.Status != StatusInUse {
		return domain.ErrInvalidStatus
	}
	if distanceKm < 0 {
		return domain.ErrInvalidDistance
	}
	r.Status = StatusReturned
	r.DistanceKm = distanceKm
	r.ReturnNote = note
	r.ReturnedAt = now
	return nil
}

func (r *Reservation) requireOwner(by string) error {
	if r.MemberID != by {
		return domain.ErrNotOwner
	}
	return nil
}
