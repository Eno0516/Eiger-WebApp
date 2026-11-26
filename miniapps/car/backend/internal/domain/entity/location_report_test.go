package entity_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/valueobject"
)

func mapLink(t *testing.T) valueobject.LocationLink {
	t.Helper()
	l, err := valueobject.NewLocationLink("https://maps.app.goo.gl/AbCdEf123")
	if err != nil {
		t.Fatal(err)
	}
	return l
}

func inUseReservation(t *testing.T) entity.Reservation {
	t.Helper()
	r := newReserved(t)
	if err := r.Start("member-1"); err != nil {
		t.Fatal(err)
	}
	return r
}

func TestNewLocationReport(t *testing.T) {
	t.Run("利用中の本人は位置を申告できる", func(t *testing.T) {
		rep, err := entity.NewLocationReport(inUseReservation(t), "member-1", mapLink(t), "駅前の駐車場", after(2))
		if err != nil {
			t.Fatalf("申告できるはず: %v", err)
		}
		if rep.CarID != "car-1" || rep.CircleID != "circle-1" || rep.ReservationID != "res-1" || rep.MemberID != "member-1" {
			t.Errorf("申告が車・サークル・予約・申告者に紐づいていない: %+v", rep)
		}
		if rep.Note != "駅前の駐車場" || !rep.ReportedAt.Equal(after(2)) {
			t.Errorf("メモと申告時刻が記録されていない: %+v", rep)
		}
	})

	t.Run("メモは空でも申告できる", func(t *testing.T) {
		if _, err := entity.NewLocationReport(inUseReservation(t), "member-1", mapLink(t), "", after(2)); err != nil {
			t.Errorf("申告できるはず: %v", err)
		}
	})

	t.Run("本人以外は申告できない", func(t *testing.T) {
		_, err := entity.NewLocationReport(inUseReservation(t), "member-2", mapLink(t), "", after(2))
		if !errors.Is(err, domain.ErrNotOwner) {
			t.Errorf("ErrNotOwner を期待: %v", err)
		}
	})

	t.Run("利用を開始していない予約では申告できない", func(t *testing.T) {
		_, err := entity.NewLocationReport(newReserved(t), "member-1", mapLink(t), "", after(2))
		if !errors.Is(err, domain.ErrInvalidStatus) {
			t.Errorf("ErrInvalidStatus を期待: %v", err)
		}
	})

	t.Run("返却済みの予約では申告できない", func(t *testing.T) {
		r := inUseReservation(t)
		_ = r.Return("member-1", 10, "", after(3))
		_, err := entity.NewLocationReport(r, "member-1", mapLink(t), "", after(4))
		if !errors.Is(err, domain.ErrInvalidStatus) {
			t.Errorf("ErrInvalidStatus を期待: %v", err)
		}
	})

	t.Run("メモが 200 文字を超えると申告できない", func(t *testing.T) {
		_, err := entity.NewLocationReport(inUseReservation(t), "member-1", mapLink(t), strings.Repeat("あ", 201), after(2))
		if !errors.Is(err, domain.ErrNoteTooLong) {
			t.Errorf("ErrNoteTooLong を期待: %v", err)
		}
	})

	t.Run("メモがちょうど 200 文字なら申告できる", func(t *testing.T) {
		_, err := entity.NewLocationReport(inUseReservation(t), "member-1", mapLink(t), strings.Repeat("あ", 200), after(2))
		if err != nil {
			t.Errorf("申告できるはず: %v", err)
		}
	})
}
