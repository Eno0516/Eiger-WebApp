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

func TestCancel(t *testing.T) {
	ctx := context.Background()

	t.Run("本人はキャンセルでき、状態が保存される", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedReservation(t, testutil.Member1, 1, 3)
		uc := reservation.NewCancel(e.Reservations, e.Members)

		if _, err := uc.Execute(ctx, reservation.CancelInput{ReservationID: r.ID, MemberID: testutil.Member1}); err != nil {
			t.Fatalf("キャンセルできるはず: %v", err)
		}
		saved, _ := e.Reservations.FindByID(ctx, r.ID)
		if saved.Status != entity.StatusCancelled {
			t.Errorf("保存された状態は cancelled のはず: %v", saved.Status)
		}
	})

	t.Run("本人以外はキャンセルできず、状態は変わらない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedReservation(t, testutil.Member1, 1, 3)
		uc := reservation.NewCancel(e.Reservations, e.Members)

		_, err := uc.Execute(ctx, reservation.CancelInput{ReservationID: r.ID, MemberID: testutil.Member2})
		if !errors.Is(err, domain.ErrNotOwner) {
			t.Fatalf("ErrNotOwner を期待: %v", err)
		}
		saved, _ := e.Reservations.FindByID(ctx, r.ID)
		if saved.Status != entity.StatusReserved {
			t.Errorf("状態は変わらないはず: %v", saved.Status)
		}
	})

	t.Run("サークルのメンバーでなくなった人はキャンセルできない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedReservation(t, testutil.Outsider, 1, 3) // 予約後に脱退した想定
		uc := reservation.NewCancel(e.Reservations, e.Members)

		_, err := uc.Execute(ctx, reservation.CancelInput{ReservationID: r.ID, MemberID: testutil.Outsider})
		if !errors.Is(err, domain.ErrNotMember) {
			t.Errorf("ErrNotMember を期待: %v", err)
		}
	})

	t.Run("存在しない予約はキャンセルできない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		uc := reservation.NewCancel(e.Reservations, e.Members)
		_, err := uc.Execute(ctx, reservation.CancelInput{ReservationID: "none", MemberID: testutil.Member1})
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("ErrNotFound を期待: %v", err)
		}
	})
}

func TestStart(t *testing.T) {
	ctx := context.Background()

	t.Run("本人は利用を開始でき、状態が保存される", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedReservation(t, testutil.Member1, 1, 3)
		uc := reservation.NewStart(e.Reservations, e.Members)

		if _, err := uc.Execute(ctx, reservation.StartInput{ReservationID: r.ID, MemberID: testutil.Member1}); err != nil {
			t.Fatalf("開始できるはず: %v", err)
		}
		saved, _ := e.Reservations.FindByID(ctx, r.ID)
		if saved.Status != entity.StatusInUse {
			t.Errorf("保存された状態は in_use のはず: %v", saved.Status)
		}
	})

	t.Run("本人以外は利用を開始できない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedReservation(t, testutil.Member1, 1, 3)
		uc := reservation.NewStart(e.Reservations, e.Members)

		_, err := uc.Execute(ctx, reservation.StartInput{ReservationID: r.ID, MemberID: testutil.Member2})
		if !errors.Is(err, domain.ErrNotOwner) {
			t.Errorf("ErrNotOwner を期待: %v", err)
		}
	})
}
