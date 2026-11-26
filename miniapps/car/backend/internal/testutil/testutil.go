// Package testutil は、usecase のテストで使う fake と、共通の準備(環境)を提供する。
// テスト以外から import しないこと。
package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/valueobject"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/repository/memory"
)

const (
	CircleID = "circle-1"
	CarID    = "car-1"
	Member1  = "member-1"
	Member2  = "member-2"
	// Outsider は、どのサークルのメンバーでもない人。
	Outsider = "outsider"
)

// FakeDirectory は、本体が実装する「サークルのメンバー判定」の fake。
type FakeDirectory struct {
	members map[string]map[string]bool
}

func NewDirectory() *FakeDirectory {
	return &FakeDirectory{members: map[string]map[string]bool{}}
}

func (d *FakeDirectory) Add(circleID, memberID string) {
	if d.members[circleID] == nil {
		d.members[circleID] = map[string]bool{}
	}
	d.members[circleID][memberID] = true
}

func (d *FakeDirectory) IsMember(_ context.Context, circleID, memberID string) (bool, error) {
	return d.members[circleID][memberID], nil
}

// FixedClock は、常に同じ時刻を返す時計。
type FixedClock struct{ T time.Time }

func (c FixedClock) Now() time.Time { return c.T }

// SeqID は、"id-1", "id-2", ... と採番する。
type SeqID struct{ n int }

func (s *SeqID) NewID() string {
	s.n++
	return fmt.Sprintf("id-%d", s.n)
}

// Env は、usecase のテストの共通の準備。
//   - サークル circle-1 に、メンバー member-1 と member-2 がいる(outsider はいない)
//   - circle-1 に、利用できる車 car-1 が 1 台ある
//   - 現在時刻は 2026-12-01 09:00 UTC で固定
type Env struct {
	Cars         *memory.CarRepository
	Reservations *memory.ReservationRepository
	Locations    *memory.LocationReportRepository
	Members      *FakeDirectory
	Clock        FixedClock
	IDs          *SeqID
	Car          entity.Car
}

func NewEnv(t *testing.T) *Env {
	t.Helper()
	e := &Env{
		Cars:         memory.NewCarRepository(),
		Reservations: memory.NewReservationRepository(),
		Locations:    memory.NewLocationReportRepository(),
		Members:      NewDirectory(),
		Clock:        FixedClock{T: time.Date(2026, 12, 1, 9, 0, 0, 0, time.UTC)},
		IDs:          &SeqID{},
		Car:          entity.Car{ID: CarID, CircleID: CircleID, Name: "サークルカー", PlateNumber: "品川 300 あ 12-34", Available: true},
	}
	e.Members.Add(CircleID, Member1)
	e.Members.Add(CircleID, Member2)
	if err := e.Cars.Save(context.Background(), e.Car); err != nil {
		t.Fatal(err)
	}
	return e
}

// After は、現在時刻から hours 時間後の時刻を返す。
func (e *Env) After(hours int) time.Time { return e.Clock.T.Add(time.Duration(hours) * time.Hour) }

func (e *Env) Period(t *testing.T, startH, endH int) valueobject.TimeRange {
	t.Helper()
	p, err := valueobject.NewTimeRange(e.After(startH), e.After(endH))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// SeedReservation は、memberID が car-1 を startH〜endH 時間後で予約した状態(reserved)を保存して返す。
func (e *Env) SeedReservation(t *testing.T, memberID string, startH, endH int) entity.Reservation {
	t.Helper()
	r, err := entity.NewReservation(e.IDs.NewID(), e.Car, memberID, e.Period(t, startH, endH), nil, e.Clock.T)
	if err != nil {
		t.Fatal(err)
	}
	if err := e.Reservations.Save(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	return r
}

// SeedInUse は、SeedReservation した予約を、利用中(in_use)にして保存して返す。
func (e *Env) SeedInUse(t *testing.T, memberID string, startH, endH int) entity.Reservation {
	t.Helper()
	r := e.SeedReservation(t, memberID, startH, endH)
	if err := r.Start(memberID); err != nil {
		t.Fatal(err)
	}
	if err := e.Reservations.Save(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	return r
}
