package entity_test

import (
	"errors"
	"testing"
	"time"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/valueobject"
)

// テスト全体で使う「現在時刻」と、そこからの時刻ヘルパー。
// 現在時刻はドメインが時計を持たず、引数で受け取る設計にしている。
var now = time.Date(2026, 12, 1, 9, 0, 0, 0, time.UTC)

func after(hours int) time.Time { return now.Add(time.Duration(hours) * time.Hour) }

func period(t *testing.T, startH, endH int) valueobject.TimeRange {
	t.Helper()
	p, err := valueobject.NewTimeRange(after(startH), after(endH))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func availableCar() entity.Car {
	return entity.Car{ID: "car-1", CircleID: "circle-1", Name: "サークルカー", PlateNumber: "品川 300 あ 12-34", Available: true}
}

// newReserved は、car-1 を member-1 が 1〜3 時間後に予約した状態を作る。
func newReserved(t *testing.T) entity.Reservation {
	t.Helper()
	r, err := entity.NewReservation("res-1", availableCar(), "member-1", period(t, 1, 3), nil, now)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestNewReservation(t *testing.T) {
	t.Run("条件を満たせば予約でき、状態は reserved になる", func(t *testing.T) {
		r, err := entity.NewReservation("res-1", availableCar(), "member-1", period(t, 1, 3), nil, now)
		if err != nil {
			t.Fatalf("予約できるはず: %v", err)
		}
		if r.Status != entity.StatusReserved {
			t.Errorf("状態は reserved のはず: %v", r.Status)
		}
		if r.CarID != "car-1" || r.CircleID != "circle-1" || r.MemberID != "member-1" {
			t.Errorf("車・サークル・借りる人が保持されていない: %+v", r)
		}
	})

	t.Run("過去から始まる予約は作れない", func(t *testing.T) {
		_, err := entity.NewReservation("res-1", availableCar(), "member-1", period(t, -1, 2), nil, now)
		if !errors.Is(err, domain.ErrPastPeriod) {
			t.Errorf("ErrPastPeriod を期待: %v", err)
		}
	})

	t.Run("ちょうど今から始まる予約は作れる", func(t *testing.T) {
		if _, err := entity.NewReservation("res-1", availableCar(), "member-1", period(t, 0, 2), nil, now); err != nil {
			t.Errorf("作れるはず: %v", err)
		}
	})

	t.Run("利用不可の車は予約できない", func(t *testing.T) {
		car := availableCar()
		car.Available = false
		_, err := entity.NewReservation("res-1", car, "member-1", period(t, 1, 3), nil, now)
		if !errors.Is(err, domain.ErrCarUnavailable) {
			t.Errorf("ErrCarUnavailable を期待: %v", err)
		}
	})
}

func TestNewReservation_時間帯の重複(t *testing.T) {
	existing := func(status entity.ReservationStatus, carID string) []entity.Reservation {
		return []entity.Reservation{{ID: "res-0", CarID: carID, Period: period(t, 2, 4), Status: status}}
	}
	reserve := func(others []entity.Reservation) error {
		_, err := entity.NewReservation("res-1", availableCar(), "member-1", period(t, 3, 5), others, now)
		return err
	}

	t.Run("同じ車の予約済みの時間帯と重なると予約できない", func(t *testing.T) {
		if err := reserve(existing(entity.StatusReserved, "car-1")); !errors.Is(err, domain.ErrReservationConflict) {
			t.Errorf("ErrReservationConflict を期待: %v", err)
		}
	})

	t.Run("同じ車の利用中の時間帯と重なると予約できない", func(t *testing.T) {
		if err := reserve(existing(entity.StatusInUse, "car-1")); !errors.Is(err, domain.ErrReservationConflict) {
			t.Errorf("ErrReservationConflict を期待: %v", err)
		}
	})

	t.Run("キャンセル済みの予約とは重なっても予約できる", func(t *testing.T) {
		if err := reserve(existing(entity.StatusCancelled, "car-1")); err != nil {
			t.Errorf("予約できるはず: %v", err)
		}
	})

	t.Run("返却済みの予約とは重なっても予約できる", func(t *testing.T) {
		if err := reserve(existing(entity.StatusReturned, "car-1")); err != nil {
			t.Errorf("予約できるはず: %v", err)
		}
	})

	t.Run("別の車の予約とは重なっても予約できる", func(t *testing.T) {
		if err := reserve(existing(entity.StatusReserved, "car-2")); err != nil {
			t.Errorf("予約できるはず: %v", err)
		}
	})

	t.Run("時間帯が接しているだけなら予約できる", func(t *testing.T) {
		others := []entity.Reservation{{ID: "res-0", CarID: "car-1", Period: period(t, 1, 3), Status: entity.StatusReserved}}
		if err := reserve(others); err != nil {
			t.Errorf("予約できるはず: %v", err)
		}
	})
}

func TestReservation_Cancel(t *testing.T) {
	t.Run("予約した本人は、予約済みの予約をキャンセルできる", func(t *testing.T) {
		r := newReserved(t)
		if err := r.Cancel("member-1"); err != nil {
			t.Fatalf("キャンセルできるはず: %v", err)
		}
		if r.Status != entity.StatusCancelled {
			t.Errorf("状態は cancelled のはず: %v", r.Status)
		}
	})

	t.Run("本人以外はキャンセルできない", func(t *testing.T) {
		r := newReserved(t)
		if err := r.Cancel("member-2"); !errors.Is(err, domain.ErrNotOwner) {
			t.Errorf("ErrNotOwner を期待: %v", err)
		}
		if r.Status != entity.StatusReserved {
			t.Errorf("状態は変わらないはず: %v", r.Status)
		}
	})

	t.Run("利用中の予約はキャンセルできない", func(t *testing.T) {
		r := newReserved(t)
		_ = r.Start("member-1")
		if err := r.Cancel("member-1"); !errors.Is(err, domain.ErrInvalidStatus) {
			t.Errorf("ErrInvalidStatus を期待: %v", err)
		}
	})
}

func TestReservation_Start(t *testing.T) {
	t.Run("予約した本人は利用を開始でき、状態は in_use になる", func(t *testing.T) {
		r := newReserved(t)
		if err := r.Start("member-1"); err != nil {
			t.Fatalf("開始できるはず: %v", err)
		}
		if r.Status != entity.StatusInUse {
			t.Errorf("状態は in_use のはず: %v", r.Status)
		}
	})

	t.Run("本人以外は利用を開始できない", func(t *testing.T) {
		r := newReserved(t)
		if err := r.Start("member-2"); !errors.Is(err, domain.ErrNotOwner) {
			t.Errorf("ErrNotOwner を期待: %v", err)
		}
	})

	t.Run("キャンセル済みの予約は利用を開始できない", func(t *testing.T) {
		r := newReserved(t)
		_ = r.Cancel("member-1")
		if err := r.Start("member-1"); !errors.Is(err, domain.ErrInvalidStatus) {
			t.Errorf("ErrInvalidStatus を期待: %v", err)
		}
	})
}

func TestReservation_Return(t *testing.T) {
	inUse := func(t *testing.T) entity.Reservation {
		r := newReserved(t)
		if err := r.Start("member-1"); err != nil {
			t.Fatal(err)
		}
		return r
	}

	t.Run("返却すると走行距離とメモが記録され、状態は returned になる", func(t *testing.T) {
		r := inUse(t)
		if err := r.Return("member-1", 42, "給油済み", after(3)); err != nil {
			t.Fatalf("返却できるはず: %v", err)
		}
		if r.Status != entity.StatusReturned {
			t.Errorf("状態は returned のはず: %v", r.Status)
		}
		if r.DistanceKm != 42 || r.ReturnNote != "給油済み" || !r.ReturnedAt.Equal(after(3)) {
			t.Errorf("返却内容が記録されていない: %+v", r)
		}
	})

	t.Run("走行距離が 0 でも返却できる", func(t *testing.T) {
		r := inUse(t)
		if err := r.Return("member-1", 0, "", after(3)); err != nil {
			t.Errorf("返却できるはず: %v", err)
		}
	})

	t.Run("走行距離が負だと返却できない", func(t *testing.T) {
		r := inUse(t)
		if err := r.Return("member-1", -1, "", after(3)); !errors.Is(err, domain.ErrInvalidDistance) {
			t.Errorf("ErrInvalidDistance を期待: %v", err)
		}
		if r.Status != entity.StatusInUse {
			t.Errorf("状態は変わらないはず: %v", r.Status)
		}
	})

	t.Run("本人以外は返却できない", func(t *testing.T) {
		r := inUse(t)
		if err := r.Return("member-2", 10, "", after(3)); !errors.Is(err, domain.ErrNotOwner) {
			t.Errorf("ErrNotOwner を期待: %v", err)
		}
	})

	t.Run("利用を開始していない予約は返却できない", func(t *testing.T) {
		r := newReserved(t)
		if err := r.Return("member-1", 10, "", after(3)); !errors.Is(err, domain.ErrInvalidStatus) {
			t.Errorf("ErrInvalidStatus を期待: %v", err)
		}
	})
}
