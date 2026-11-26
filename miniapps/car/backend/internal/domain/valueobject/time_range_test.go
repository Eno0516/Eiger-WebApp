package valueobject_test

import (
	"errors"
	"testing"
	"time"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/valueobject"
)

var base = time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

func at(hour int) time.Time { return base.Add(time.Duration(hour) * time.Hour) }

func TestNewTimeRange(t *testing.T) {
	t.Run("終了が開始より後なら期間を作れる", func(t *testing.T) {
		r, err := valueobject.NewTimeRange(at(9), at(12))
		if err != nil {
			t.Fatalf("エラーなく作れるはず: %v", err)
		}
		if !r.Start().Equal(at(9)) || !r.End().Equal(at(12)) {
			t.Errorf("開始と終了が保持されていない: %v - %v", r.Start(), r.End())
		}
	})

	t.Run("終了が開始と同じ時刻だと作れない", func(t *testing.T) {
		_, err := valueobject.NewTimeRange(at(9), at(9))
		if !errors.Is(err, domain.ErrInvalidPeriod) {
			t.Errorf("ErrInvalidPeriod を期待: %v", err)
		}
	})

	t.Run("終了が開始より前だと作れない", func(t *testing.T) {
		_, err := valueobject.NewTimeRange(at(12), at(9))
		if !errors.Is(err, domain.ErrInvalidPeriod) {
			t.Errorf("ErrInvalidPeriod を期待: %v", err)
		}
	})
}

func TestTimeRange_Overlaps(t *testing.T) {
	// 基準となる期間は 10:00-12:00
	ref, _ := valueobject.NewTimeRange(at(10), at(12))

	cases := []struct {
		name       string
		start, end int
		want       bool
	}{
		{"一部が重なっていると重複する(前側)", 9, 11, true},
		{"一部が重なっていると重複する(後ろ側)", 11, 13, true},
		{"基準を完全に含むと重複する", 9, 13, true},
		{"基準に完全に含まれると重複する", 10, 11, true},
		{"前の期間の終了と開始が接しているだけなら重複しない", 8, 10, false},
		{"後ろの期間の開始と終了が接しているだけなら重複しない", 12, 14, false},
		{"完全に離れていれば重複しない", 14, 16, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			other, err := valueobject.NewTimeRange(at(c.start), at(c.end))
			if err != nil {
				t.Fatal(err)
			}
			if got := ref.Overlaps(other); got != c.want {
				t.Errorf("Overlaps = %v, 期待 %v", got, c.want)
			}
			// 重複判定は対称であること
			if got := other.Overlaps(ref); got != c.want {
				t.Errorf("逆向きの Overlaps = %v, 期待 %v", got, c.want)
			}
		})
	}
}
