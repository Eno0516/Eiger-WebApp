package reservation_test

import (
	"context"
	"errors"
	"testing"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/testutil"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/reservation"
)

func TestListByRange(t *testing.T) {
	ctx := context.Background()

	t.Run("メンバーは、カレンダーの期間と重なる予約を、開始の早い順に見られる", func(t *testing.T) {
		e := testutil.NewEnv(t)
		late := e.SeedReservation(t, testutil.Member1, 6, 8)
		early := e.SeedReservation(t, testutil.Member2, 1, 3)
		e.SeedReservation(t, testutil.Member1, 40, 42) // 期間外
		uc := reservation.NewListByRange(e.Reservations, e.Members)

		got, err := uc.Execute(ctx, reservation.ListByRangeInput{
			CircleID: testutil.CircleID, MemberID: testutil.Member1, From: e.After(0), To: e.After(24),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 2 || got[0].ID != early.ID || got[1].ID != late.ID {
			t.Errorf("early, late の順に 2 件のはず: %+v", got)
		}
	})

	t.Run("キャンセル済みの予約はカレンダーに出さない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedReservation(t, testutil.Member1, 1, 3)
		if _, err := reservation.NewCancel(e.Reservations, e.Members).Execute(ctx, reservation.CancelInput{ReservationID: r.ID, MemberID: testutil.Member1}); err != nil {
			t.Fatal(err)
		}
		uc := reservation.NewListByRange(e.Reservations, e.Members)

		got, err := uc.Execute(ctx, reservation.ListByRangeInput{
			CircleID: testutil.CircleID, MemberID: testutil.Member1, From: e.After(0), To: e.After(24),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Errorf("0 件のはず: %+v", got)
		}
	})

	t.Run("サークルのメンバーでない人は見られない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		uc := reservation.NewListByRange(e.Reservations, e.Members)
		_, err := uc.Execute(ctx, reservation.ListByRangeInput{
			CircleID: testutil.CircleID, MemberID: testutil.Outsider, From: e.After(0), To: e.After(24),
		})
		if !errors.Is(err, domain.ErrNotMember) {
			t.Errorf("ErrNotMember を期待: %v", err)
		}
	})

	t.Run("期間が不正だと見られない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		uc := reservation.NewListByRange(e.Reservations, e.Members)
		_, err := uc.Execute(ctx, reservation.ListByRangeInput{
			CircleID: testutil.CircleID, MemberID: testutil.Member1, From: e.After(24), To: e.After(0),
		})
		if !errors.Is(err, domain.ErrInvalidPeriod) {
			t.Errorf("ErrInvalidPeriod を期待: %v", err)
		}
	})
}

func TestListMine(t *testing.T) {
	ctx := context.Background()

	t.Run("自分の予約だけを見られる", func(t *testing.T) {
		e := testutil.NewEnv(t)
		mine := e.SeedReservation(t, testutil.Member1, 1, 3)
		e.SeedReservation(t, testutil.Member2, 5, 7)
		uc := reservation.NewListMine(e.Reservations, e.Members)

		got, err := uc.Execute(ctx, reservation.ListMineInput{CircleID: testutil.CircleID, MemberID: testutil.Member1})
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0].ID != mine.ID {
			t.Errorf("自分の 1 件だけのはず: %+v", got)
		}
	})

	t.Run("サークルのメンバーでない人は見られない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		uc := reservation.NewListMine(e.Reservations, e.Members)
		_, err := uc.Execute(ctx, reservation.ListMineInput{CircleID: testutil.CircleID, MemberID: testutil.Outsider})
		if !errors.Is(err, domain.ErrNotMember) {
			t.Errorf("ErrNotMember を期待: %v", err)
		}
	})
}
