// Package repotest は、repository 実装が満たすべき契約テストを提供する。
// memory 実装と postgres 実装(#9)の両方に、同じテストを流して振る舞いを揃える。
package repotest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/valueobject"
)

var base = time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

func rng(t *testing.T, startH, endH int) valueobject.TimeRange {
	t.Helper()
	r, err := valueobject.NewTimeRange(base.Add(time.Duration(startH)*time.Hour), base.Add(time.Duration(endH)*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

type spec struct {
	id, carID, circleID, memberID string
	startH, endH                  int
	status                        entity.ReservationStatus
}

func mk(t *testing.T, s spec) entity.Reservation {
	t.Helper()
	return entity.Reservation{
		ID: s.id, CarID: s.carID, CircleID: s.circleID, MemberID: s.memberID,
		Period: rng(t, s.startH, s.endH), Status: s.status,
	}
}

func ids(rs []entity.Reservation) []string {
	out := make([]string, 0, len(rs))
	for _, r := range rs {
		out = append(out, r.ID)
	}
	return out
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// RunReservationRepositoryContract は ReservationRepository の契約テスト。newRepo は毎回、空の repository を返すこと。
func RunReservationRepositoryContract(t *testing.T, newRepo func() repository.ReservationRepository) {
	ctx := context.Background()
	save := func(t *testing.T, repo repository.ReservationRepository, s spec) error {
		t.Helper()
		return repo.Save(ctx, mk(t, s))
	}
	mustSave := func(t *testing.T, repo repository.ReservationRepository, s spec) {
		t.Helper()
		if err := save(t, repo, s); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("保存した予約を ID で取得できる", func(t *testing.T) {
		repo := newRepo()
		mustSave(t, repo, spec{"r1", "car-1", "circle-1", "m1", 10, 12, entity.StatusReserved})
		got, err := repo.FindByID(ctx, "r1")
		if err != nil {
			t.Fatal(err)
		}
		if got.CarID != "car-1" || got.MemberID != "m1" || got.Status != entity.StatusReserved || !got.Period.Start().Equal(base.Add(10*time.Hour)) {
			t.Errorf("保存した内容と一致しない: %+v", got)
		}
	})

	t.Run("存在しない ID は ErrNotFound になる", func(t *testing.T) {
		if _, err := newRepo().FindByID(ctx, "none"); !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("ErrNotFound を期待: %v", err)
		}
	})

	t.Run("同じ車の予約中の時間帯と重なる予約は保存できない", func(t *testing.T) {
		repo := newRepo()
		mustSave(t, repo, spec{"r1", "car-1", "circle-1", "m1", 10, 12, entity.StatusReserved})
		err := save(t, repo, spec{"r2", "car-1", "circle-1", "m2", 11, 13, entity.StatusReserved})
		if !errors.Is(err, domain.ErrReservationConflict) {
			t.Errorf("ErrReservationConflict を期待: %v", err)
		}
		if _, err := repo.FindByID(ctx, "r2"); !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("競合した予約は保存されないはず: %v", err)
		}
	})

	t.Run("キャンセル済みの予約と重なっていても保存できる", func(t *testing.T) {
		repo := newRepo()
		mustSave(t, repo, spec{"r1", "car-1", "circle-1", "m1", 10, 12, entity.StatusCancelled})
		if err := save(t, repo, spec{"r2", "car-1", "circle-1", "m2", 11, 13, entity.StatusReserved}); err != nil {
			t.Errorf("保存できるはず: %v", err)
		}
	})

	t.Run("別の車なら、時間帯が重なっていても保存できる", func(t *testing.T) {
		repo := newRepo()
		mustSave(t, repo, spec{"r1", "car-1", "circle-1", "m1", 10, 12, entity.StatusReserved})
		if err := save(t, repo, spec{"r2", "car-2", "circle-1", "m2", 11, 13, entity.StatusReserved}); err != nil {
			t.Errorf("保存できるはず: %v", err)
		}
	})

	t.Run("同じ予約を更新して保存しても、自分自身とは衝突しない", func(t *testing.T) {
		repo := newRepo()
		mustSave(t, repo, spec{"r1", "car-1", "circle-1", "m1", 10, 12, entity.StatusReserved})
		if err := save(t, repo, spec{"r1", "car-1", "circle-1", "m1", 10, 12, entity.StatusInUse}); err != nil {
			t.Fatalf("更新できるはず: %v", err)
		}
		got, _ := repo.FindByID(ctx, "r1")
		if got.Status != entity.StatusInUse {
			t.Errorf("更新が反映されていない: %v", got.Status)
		}
	})

	t.Run("車の予約を、期間と重なるものだけで取得できる(接しているだけは含まない)", func(t *testing.T) {
		repo := newRepo()
		mustSave(t, repo, spec{"in", "car-1", "circle-1", "m1", 10, 12, entity.StatusReserved})
		mustSave(t, repo, spec{"touch", "car-1", "circle-1", "m1", 12, 14, entity.StatusReserved})
		mustSave(t, repo, spec{"other-car", "car-2", "circle-1", "m1", 10, 12, entity.StatusReserved})
		got, err := repo.ListByCarInRange(ctx, "car-1", rng(t, 9, 12))
		if err != nil {
			t.Fatal(err)
		}
		if !equal(ids(got), []string{"in"}) {
			t.Errorf("in だけのはず: %v", ids(got))
		}
	})

	t.Run("サークルの予約を、期間で絞り、開始の早い順に取得できる", func(t *testing.T) {
		repo := newRepo()
		mustSave(t, repo, spec{"late", "car-1", "circle-1", "m1", 14, 16, entity.StatusReserved})
		mustSave(t, repo, spec{"early", "car-2", "circle-1", "m1", 9, 11, entity.StatusReserved})
		mustSave(t, repo, spec{"outside", "car-1", "circle-1", "m1", 20, 22, entity.StatusReserved})
		mustSave(t, repo, spec{"other-circle", "car-3", "circle-2", "m1", 10, 12, entity.StatusReserved})
		got, err := repo.ListByCircleInRange(ctx, "circle-1", rng(t, 8, 18))
		if err != nil {
			t.Fatal(err)
		}
		if !equal(ids(got), []string{"early", "late"}) {
			t.Errorf("early, late の順のはず: %v", ids(got))
		}
	})

	t.Run("メンバーの予約だけを、開始の早い順に取得できる", func(t *testing.T) {
		repo := newRepo()
		mustSave(t, repo, spec{"b", "car-1", "circle-1", "m1", 14, 16, entity.StatusReserved})
		mustSave(t, repo, spec{"a", "car-2", "circle-1", "m1", 9, 11, entity.StatusReturned})
		mustSave(t, repo, spec{"others", "car-1", "circle-1", "m2", 20, 22, entity.StatusReserved})
		got, err := repo.ListByMember(ctx, "circle-1", "m1")
		if err != nil {
			t.Fatal(err)
		}
		if !equal(ids(got), []string{"a", "b"}) {
			t.Errorf("a, b の順のはず: %v", ids(got))
		}
	})
}
