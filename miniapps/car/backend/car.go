// Package car は、「車の予約」ミニアプリの、本体から見える唯一の入口。
// 本体は NewHandler で受け取った http.Handler を、好きなパスに Mount するだけでよい
// (例: chi なら r.Mount("/api/v1/miniapps/car", car.NewHandler(deps)))。
// ミニアプリの内部(internal/)には、本体からは触れない。
package car

import (
	"context"
	"net/http"

	httphandler "github.com/enoro/eiger-webapp/miniapps/car/backend/internal/handler/http"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/infrastructure/system"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/repository/memory"
	carusecase "github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/car"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/location"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/reservation"
)

// MemberDirectory は、「そのサークルのメンバーか」の判定。本体(サークル管理)が実装して渡す。
type MemberDirectory interface {
	IsMember(ctx context.Context, circleID, memberID string) (bool, error)
}

// Identity は、リクエストから「誰が操作しているか」を取り出す。本体の認証(JWT)が実装して渡す。
// 取り出せなければ ok=false(未認証)を返す。
type Identity interface {
	MemberID(r *http.Request) (memberID string, ok bool)
}

// Deps は、本体がミニアプリに渡すもの。
type Deps struct {
	Members  MemberDirectory
	Identity Identity
}

// NewHandler は、車の予約 API の http.Handler を返す。
// TODO(#9): 現状は memory repository(再起動で消える)。PostgreSQL 実装ができたら、Deps で DB を受け取って差し替える。
func NewHandler(deps Deps) http.Handler {
	cars := memory.NewCarRepository()
	reservations := memory.NewReservationRepository()
	locations := memory.NewLocationReportRepository()
	clock := system.Clock{}
	ids := system.RandomID{}

	return httphandler.New(httphandler.UseCases{
		ListCars:          carusecase.NewList(cars, deps.Members),
		CreateReservation: reservation.NewCreate(cars, reservations, deps.Members, clock, ids),
		CancelReservation: reservation.NewCancel(reservations, deps.Members),
		StartReservation:  reservation.NewStart(reservations, deps.Members),
		ReturnReservation: reservation.NewReturn(reservations, locations, deps.Members, clock),
		ListByRange:       reservation.NewListByRange(reservations, deps.Members),
		ListMine:          reservation.NewListMine(reservations, deps.Members),
		ReportLocation:    location.NewReport(reservations, locations, deps.Members, clock),
		GetLocation:       location.NewGet(cars, locations, deps.Members),
	}, deps.Identity)
}
