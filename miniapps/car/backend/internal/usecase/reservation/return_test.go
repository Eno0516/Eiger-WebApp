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

func newReturn(e *testutil.Env) *reservation.ReturnUseCase {
	return reservation.NewReturn(e.Reservations, e.Locations, e.Members, e.Clock)
}

func TestReturn(t *testing.T) {
	ctx := context.Background()

	t.Run("返却すると、走行距離・メモ・返却時刻が保存される", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedInUse(t, testutil.Member1, 1, 3)

		_, err := newReturn(e).Execute(ctx, reservation.ReturnInput{
			ReservationID: r.ID, MemberID: testutil.Member1, DistanceKm: 42, Note: "給油済み",
		})
		if err != nil {
			t.Fatalf("返却できるはず: %v", err)
		}
		saved, _ := e.Reservations.FindByID(ctx, r.ID)
		if saved.Status != entity.StatusReturned || saved.DistanceKm != 42 || saved.ReturnNote != "給油済み" || !saved.ReturnedAt.Equal(e.Clock.T) {
			t.Errorf("返却内容が保存されていない: %+v", saved)
		}
	})

	t.Run("返却時に駐車場所を渡すと、車の位置として申告される", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedInUse(t, testutil.Member1, 1, 3)

		_, err := newReturn(e).Execute(ctx, reservation.ReturnInput{
			ReservationID: r.ID, MemberID: testutil.Member1, DistanceKm: 10,
			ParkingLink: "https://maps.app.goo.gl/AbCdEf123", ParkingNote: "駅前の駐車場",
		})
		if err != nil {
			t.Fatalf("返却できるはず: %v", err)
		}
		latest, err := e.Locations.LatestByCar(ctx, testutil.CarID)
		if err != nil {
			t.Fatalf("駐車場所が申告されているはず: %v", err)
		}
		if latest.Link.String() != "https://maps.app.goo.gl/AbCdEf123" || latest.Note != "駅前の駐車場" {
			t.Errorf("申告内容が想定と違う: %+v", latest)
		}
	})

	t.Run("駐車場所のリンクが不正なら、返却されず状態も変わらない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedInUse(t, testutil.Member1, 1, 3)

		_, err := newReturn(e).Execute(ctx, reservation.ReturnInput{
			ReservationID: r.ID, MemberID: testutil.Member1, DistanceKm: 10,
			ParkingLink: "https://evil.example.com/maps",
		})
		if !errors.Is(err, domain.ErrInvalidLocationLink) {
			t.Fatalf("ErrInvalidLocationLink を期待: %v", err)
		}
		saved, _ := e.Reservations.FindByID(ctx, r.ID)
		if saved.Status != entity.StatusInUse {
			t.Errorf("状態は変わらないはず: %v", saved.Status)
		}
	})

	t.Run("走行距離が不正なら、駐車場所も申告されない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedInUse(t, testutil.Member1, 1, 3)

		_, err := newReturn(e).Execute(ctx, reservation.ReturnInput{
			ReservationID: r.ID, MemberID: testutil.Member1, DistanceKm: -1,
			ParkingLink: "https://maps.app.goo.gl/AbCdEf123",
		})
		if !errors.Is(err, domain.ErrInvalidDistance) {
			t.Fatalf("ErrInvalidDistance を期待: %v", err)
		}
		if _, err := e.Locations.LatestByCar(ctx, testutil.CarID); !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("位置は申告されないはず: %v", err)
		}
	})

	t.Run("本人以外は返却できない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedInUse(t, testutil.Member1, 1, 3)

		_, err := newReturn(e).Execute(ctx, reservation.ReturnInput{ReservationID: r.ID, MemberID: testutil.Member2, DistanceKm: 10})
		if !errors.Is(err, domain.ErrNotOwner) {
			t.Errorf("ErrNotOwner を期待: %v", err)
		}
	})
}
