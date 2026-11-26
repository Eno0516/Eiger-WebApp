package reservation_test

import (
	"context"
	"errors"
	"testing"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/testutil"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/reservation"
)

func newCreate(e *testutil.Env) *reservation.CreateUseCase {
	return reservation.NewCreate(e.Cars, e.Reservations, e.Members, e.Clock, e.IDs)
}

func TestCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("サークルのメンバーは車を予約でき、予約が保存される", func(t *testing.T) {
		e := testutil.NewEnv(t)
		got, err := newCreate(e).Execute(ctx, reservation.CreateInput{
			CarID: testutil.CarID, MemberID: testutil.Member1, Start: e.After(1), End: e.After(3),
		})
		if err != nil {
			t.Fatalf("予約できるはず: %v", err)
		}
		if got.ID == "" || got.Status != entity.StatusReserved || got.MemberID != testutil.Member1 || got.CircleID != testutil.CircleID {
			t.Errorf("予約の内容が想定と違う: %+v", got)
		}
		saved, err := e.Reservations.FindByID(ctx, got.ID)
		if err != nil || saved.ID != got.ID {
			t.Errorf("予約が保存されていない: %v", err)
		}
	})

	t.Run("サークルのメンバーでない人は予約できない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		_, err := newCreate(e).Execute(ctx, reservation.CreateInput{
			CarID: testutil.CarID, MemberID: testutil.Outsider, Start: e.After(1), End: e.After(3),
		})
		if !errors.Is(err, domain.ErrNotMember) {
			t.Errorf("ErrNotMember を期待: %v", err)
		}
	})

	t.Run("存在しない車は予約できない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		_, err := newCreate(e).Execute(ctx, reservation.CreateInput{
			CarID: "none", MemberID: testutil.Member1, Start: e.After(1), End: e.After(3),
		})
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("ErrNotFound を期待: %v", err)
		}
	})

	t.Run("期間が不正だと予約できない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		_, err := newCreate(e).Execute(ctx, reservation.CreateInput{
			CarID: testutil.CarID, MemberID: testutil.Member1, Start: e.After(3), End: e.After(1),
		})
		if !errors.Is(err, domain.ErrInvalidPeriod) {
			t.Errorf("ErrInvalidPeriod を期待: %v", err)
		}
	})

	t.Run("現在時刻より前から始まる予約はできない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		_, err := newCreate(e).Execute(ctx, reservation.CreateInput{
			CarID: testutil.CarID, MemberID: testutil.Member1, Start: e.After(-1), End: e.After(2),
		})
		if !errors.Is(err, domain.ErrPastPeriod) {
			t.Errorf("ErrPastPeriod を期待: %v", err)
		}
	})

	t.Run("同じ時間帯に別の予約があると予約できず、2 件目は保存されない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		e.SeedReservation(t, testutil.Member2, 2, 4)
		_, err := newCreate(e).Execute(ctx, reservation.CreateInput{
			CarID: testutil.CarID, MemberID: testutil.Member1, Start: e.After(3), End: e.After(5),
		})
		if !errors.Is(err, domain.ErrReservationConflict) {
			t.Fatalf("ErrReservationConflict を期待: %v", err)
		}
		list, _ := e.Reservations.ListByCarInRange(ctx, testutil.CarID, e.Period(t, 0, 10))
		if len(list) != 1 {
			t.Errorf("予約は最初の 1 件だけのはず: %d 件", len(list))
		}
	})
}
