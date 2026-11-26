package entity

import (
	"time"
	"unicode/utf8"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/valueobject"
)

const maxLocationNoteRunes = 200

// LocationReport は、借りている人による「いまここ」の自己申告(GPS の自動追跡ではない)。
// 返却時の駐車場所も、返却の直前(利用中)に申告する。
type LocationReport struct {
	CarID         string
	CircleID      string
	ReservationID string
	MemberID      string
	Link          valueobject.LocationLink
	Note          string
	ReportedAt    time.Time
}

// NewLocationReport は位置の申告を作る。申告できるのは、利用中の予約の本人のみ。
// 閲覧できるのはそのサークルのメンバーのみだが、その判定は usecase の責務。
func NewLocationReport(
	res Reservation,
	by string,
	link valueobject.LocationLink,
	note string,
	now time.Time,
) (LocationReport, error) {
	if err := res.requireOwner(by); err != nil {
		return LocationReport{}, err
	}
	if res.Status != StatusInUse {
		return LocationReport{}, domain.ErrInvalidStatus
	}
	if utf8.RuneCountInString(note) > maxLocationNoteRunes {
		return LocationReport{}, domain.ErrNoteTooLong
	}
	return LocationReport{
		CarID:         res.CarID,
		CircleID:      res.CircleID,
		ReservationID: res.ID,
		MemberID:      by,
		Link:          link,
		Note:          note,
		ReportedAt:    now,
	}, nil
}
