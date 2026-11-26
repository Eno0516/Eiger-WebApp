// Package httphandler は、車の予約ミニアプリの HTTP ハンドラ。
// usecase を呼んで結果を JSON にするだけで、ルールは持たない。
// ルーティングは標準ライブラリの http.ServeMux(Go 1.22 のパターン)で行い、パスは本体から見て相対。
package httphandler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/car"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/location"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/reservation"
)

const maxBodyBytes = 1 << 20

// Identity は、リクエストから「誰が操作しているか」を取り出す。
// 本体の認証ミドルウェア(JWT)が実装して渡す。取り出せなければ ok=false(未認証)。
type Identity interface {
	MemberID(r *http.Request) (memberID string, ok bool)
}

type UseCases struct {
	ListCars          *car.ListUseCase
	CreateReservation *reservation.CreateUseCase
	CancelReservation *reservation.CancelUseCase
	StartReservation  *reservation.StartUseCase
	ReturnReservation *reservation.ReturnUseCase
	ListByRange       *reservation.ListByRangeUseCase
	ListMine          *reservation.ListMineUseCase
	ReportLocation    *location.ReportUseCase
	GetLocation       *location.GetUseCase
}

type handler struct {
	uc UseCases
	id Identity
}

// New は、全ルートを持つ http.Handler を返す。すべてのルートで本人確認が必要。
func New(uc UseCases, id Identity) http.Handler {
	h := &handler{uc: uc, id: id}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /circles/{circleID}/cars", h.authed(h.listCars))
	mux.HandleFunc("GET /circles/{circleID}/reservations", h.authed(h.listByRange))
	mux.HandleFunc("GET /circles/{circleID}/reservations/mine", h.authed(h.listMine))
	mux.HandleFunc("POST /cars/{carID}/reservations", h.authed(h.createReservation))
	mux.HandleFunc("GET /cars/{carID}/location", h.authed(h.getLocation))
	mux.HandleFunc("POST /reservations/{id}/cancel", h.authed(h.cancelReservation))
	mux.HandleFunc("POST /reservations/{id}/start", h.authed(h.startReservation))
	mux.HandleFunc("POST /reservations/{id}/return", h.authed(h.returnReservation))
	mux.HandleFunc("POST /reservations/{id}/location", h.authed(h.reportLocation))
	return mux
}

// authed は、本人確認できたリクエストだけを next に渡す。
func (h *handler) authed(next func(w http.ResponseWriter, r *http.Request, memberID string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memberID, ok := h.id.MemberID(r)
		if !ok {
			writeErrorCode(w, http.StatusUnauthorized, "UNAUTHORIZED", "ログインが必要です")
			return
		}
		next(w, r, memberID)
	}
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(v); err != nil {
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "リクエストの形式が不正です")
		return false
	}
	return true
}

func (h *handler) listCars(w http.ResponseWriter, r *http.Request, memberID string) {
	cars, err := h.uc.ListCars.Execute(r.Context(), car.ListInput{CircleID: r.PathValue("circleID"), MemberID: memberID})
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"cars": toCarsJSON(cars)})
}

func (h *handler) listByRange(w http.ResponseWriter, r *http.Request, memberID string) {
	from, errFrom := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
	to, errTo := time.Parse(time.RFC3339, r.URL.Query().Get("to"))
	if errFrom != nil || errTo != nil {
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "from と to を RFC3339 形式で指定してください")
		return
	}
	list, err := h.uc.ListByRange.Execute(r.Context(), reservation.ListByRangeInput{
		CircleID: r.PathValue("circleID"), MemberID: memberID, From: from, To: to,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"reservations": toReservationsJSON(list)})
}

func (h *handler) listMine(w http.ResponseWriter, r *http.Request, memberID string) {
	list, err := h.uc.ListMine.Execute(r.Context(), reservation.ListMineInput{CircleID: r.PathValue("circleID"), MemberID: memberID})
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"reservations": toReservationsJSON(list)})
}

func (h *handler) createReservation(w http.ResponseWriter, r *http.Request, memberID string) {
	var req struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}
	if !decode(w, r, &req) {
		return
	}
	res, err := h.uc.CreateReservation.Execute(r.Context(), reservation.CreateInput{
		CarID: r.PathValue("carID"), MemberID: memberID, Start: req.Start, End: req.End,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, toReservationJSON(res))
}

func (h *handler) cancelReservation(w http.ResponseWriter, r *http.Request, memberID string) {
	res, err := h.uc.CancelReservation.Execute(r.Context(), reservation.CancelInput{ReservationID: r.PathValue("id"), MemberID: memberID})
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, toReservationJSON(res))
}

func (h *handler) startReservation(w http.ResponseWriter, r *http.Request, memberID string) {
	res, err := h.uc.StartReservation.Execute(r.Context(), reservation.StartInput{ReservationID: r.PathValue("id"), MemberID: memberID})
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, toReservationJSON(res))
}

func (h *handler) returnReservation(w http.ResponseWriter, r *http.Request, memberID string) {
	var req struct {
		DistanceKm  int    `json:"distanceKm"`
		Note        string `json:"note"`
		ParkingLink string `json:"parkingLink"`
		ParkingNote string `json:"parkingNote"`
	}
	if !decode(w, r, &req) {
		return
	}
	res, err := h.uc.ReturnReservation.Execute(r.Context(), reservation.ReturnInput{
		ReservationID: r.PathValue("id"), MemberID: memberID,
		DistanceKm: req.DistanceKm, Note: req.Note, ParkingLink: req.ParkingLink, ParkingNote: req.ParkingNote,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, toReservationJSON(res))
}

func (h *handler) reportLocation(w http.ResponseWriter, r *http.Request, memberID string) {
	var req struct {
		Link string `json:"link"`
		Note string `json:"note"`
	}
	if !decode(w, r, &req) {
		return
	}
	rep, err := h.uc.ReportLocation.Execute(r.Context(), location.ReportInput{
		ReservationID: r.PathValue("id"), MemberID: memberID, Link: req.Link, Note: req.Note,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, toLocationJSON(rep))
}

func (h *handler) getLocation(w http.ResponseWriter, r *http.Request, memberID string) {
	rep, err := h.uc.GetLocation.Execute(r.Context(), location.GetInput{CarID: r.PathValue("carID"), MemberID: memberID})
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, toLocationJSON(rep))
}
