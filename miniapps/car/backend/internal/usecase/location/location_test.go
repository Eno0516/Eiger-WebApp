package location_test

import (
	"context"
	"errors"
	"testing"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/testutil"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/location"
)

const validLink = "https://maps.app.goo.gl/AbCdEf123"

func newReport(e *testutil.Env) *location.ReportUseCase {
	return location.NewReport(e.Reservations, e.Locations, e.Members, e.Clock)
}

func TestReport(t *testing.T) {
	ctx := context.Background()

	t.Run("利用中の本人は位置を申告でき、最新の位置として取得できる", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedInUse(t, testutil.Member1, 1, 3)

		if _, err := newReport(e).Execute(ctx, location.ReportInput{
			ReservationID: r.ID, MemberID: testutil.Member1, Link: validLink, Note: "駅前",
		}); err != nil {
			t.Fatalf("申告できるはず: %v", err)
		}
		got, err := location.NewGet(e.Cars, e.Locations, e.Members).Execute(ctx, location.GetInput{CarID: testutil.CarID, MemberID: testutil.Member2})
		if err != nil {
			t.Fatalf("同じサークルのメンバーは見られるはず: %v", err)
		}
		if got.Link.String() != validLink || got.Note != "駅前" || !got.ReportedAt.Equal(e.Clock.T) {
			t.Errorf("申告内容が想定と違う: %+v", got)
		}
	})

	t.Run("2 回申告すると、新しい方が最新になる", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedInUse(t, testutil.Member1, 1, 3)
		uc := newReport(e)
		_, _ = uc.Execute(ctx, location.ReportInput{ReservationID: r.ID, MemberID: testutil.Member1, Link: validLink, Note: "1 回目"})
		e.Clock.T = e.After(1)
		uc = newReport(e)
		_, _ = uc.Execute(ctx, location.ReportInput{ReservationID: r.ID, MemberID: testutil.Member1, Link: "https://goo.gl/maps/XyZ", Note: "2 回目"})

		got, err := location.NewGet(e.Cars, e.Locations, e.Members).Execute(ctx, location.GetInput{CarID: testutil.CarID, MemberID: testutil.Member1})
		if err != nil || got.Note != "2 回目" {
			t.Errorf("2 回目が最新のはず: %+v, %v", got, err)
		}
	})

	t.Run("リンクが不正だと申告できない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedInUse(t, testutil.Member1, 1, 3)
		_, err := newReport(e).Execute(ctx, location.ReportInput{ReservationID: r.ID, MemberID: testutil.Member1, Link: "https://evil.example.com/maps"})
		if !errors.Is(err, domain.ErrInvalidLocationLink) {
			t.Errorf("ErrInvalidLocationLink を期待: %v", err)
		}
	})

	t.Run("本人以外は申告できない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedInUse(t, testutil.Member1, 1, 3)
		_, err := newReport(e).Execute(ctx, location.ReportInput{ReservationID: r.ID, MemberID: testutil.Member2, Link: validLink})
		if !errors.Is(err, domain.ErrNotOwner) {
			t.Errorf("ErrNotOwner を期待: %v", err)
		}
	})

	t.Run("利用を開始していない予約では申告できない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedReservation(t, testutil.Member1, 1, 3)
		_, err := newReport(e).Execute(ctx, location.ReportInput{ReservationID: r.ID, MemberID: testutil.Member1, Link: validLink})
		if !errors.Is(err, domain.ErrInvalidStatus) {
			t.Errorf("ErrInvalidStatus を期待: %v", err)
		}
	})
}

func TestGet(t *testing.T) {
	ctx := context.Background()

	t.Run("申告がまだない車は ErrNotFound になる", func(t *testing.T) {
		e := testutil.NewEnv(t)
		_, err := location.NewGet(e.Cars, e.Locations, e.Members).Execute(ctx, location.GetInput{CarID: testutil.CarID, MemberID: testutil.Member1})
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("ErrNotFound を期待: %v", err)
		}
	})

	t.Run("サークルのメンバーでない人は位置を見られない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedInUse(t, testutil.Member1, 1, 3)
		_, _ = newReport(e).Execute(ctx, location.ReportInput{ReservationID: r.ID, MemberID: testutil.Member1, Link: validLink})

		_, err := location.NewGet(e.Cars, e.Locations, e.Members).Execute(ctx, location.GetInput{CarID: testutil.CarID, MemberID: testutil.Outsider})
		if !errors.Is(err, domain.ErrNotMember) {
			t.Errorf("ErrNotMember を期待: %v", err)
		}
	})
}
