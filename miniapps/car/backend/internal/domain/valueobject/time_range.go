package valueobject

import (
	"time"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
)

// TimeRange は [開始, 終了) の半開区間。終了の瞬間は含まない。
// このため、前の予約の終了と次の予約の開始が同時刻でも重複しない。
type TimeRange struct {
	start time.Time
	end   time.Time
}

// NewTimeRange は期間を作る。終了が開始以前なら domain.ErrInvalidPeriod を返す。
func NewTimeRange(start, end time.Time) (TimeRange, error) {
	if !end.After(start) {
		return TimeRange{}, domain.ErrInvalidPeriod
	}
	return TimeRange{start: start, end: end}, nil
}

func (r TimeRange) Start() time.Time { return r.start }
func (r TimeRange) End() time.Time   { return r.end }

// Overlaps は 2 つの期間が重なるかを返す。接しているだけなら重ならない。
func (r TimeRange) Overlaps(other TimeRange) bool {
	return r.start.Before(other.end) && other.start.Before(r.end)
}
