package httphandler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
)

// errorTable は、domain のエラーを HTTP のステータスとエラーコードに対応づける。
var errorTable = []struct {
	err    error
	status int
	code   string
}{
	{domain.ErrNotFound, http.StatusNotFound, "NOT_FOUND"},
	{domain.ErrNotMember, http.StatusForbidden, "NOT_MEMBER"},
	{domain.ErrNotOwner, http.StatusForbidden, "NOT_OWNER"},
	{domain.ErrInvalidPeriod, http.StatusBadRequest, "INVALID_PERIOD"},
	{domain.ErrPastPeriod, http.StatusBadRequest, "PAST_PERIOD"},
	{domain.ErrInvalidDistance, http.StatusBadRequest, "INVALID_DISTANCE"},
	{domain.ErrInvalidLocationLink, http.StatusBadRequest, "INVALID_LOCATION_LINK"},
	{domain.ErrNoteTooLong, http.StatusBadRequest, "NOTE_TOO_LONG"},
	{domain.ErrReservationConflict, http.StatusConflict, "RESERVATION_CONFLICT"},
	{domain.ErrInvalidStatus, http.StatusConflict, "INVALID_STATUS"},
	{domain.ErrCarUnavailable, http.StatusConflict, "CAR_UNAVAILABLE"},
}

type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeErrorCode(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorBody{Error: errorDetail{Code: code, Message: message}})
}

// writeError は、domain のエラーを対応する HTTP レスポンスにする。
// 未知のエラーは 500 とし、内部の詳細はレスポンスに出さずログにだけ残す。
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	for _, e := range errorTable {
		if errors.Is(err, e.err) {
			writeErrorCode(w, e.status, e.code, e.err.Error())
			return
		}
	}
	slog.ErrorContext(r.Context(), "車の予約: 想定外のエラー", "error", err, "method", r.Method, "path", r.URL.Path)
	writeErrorCode(w, http.StatusInternalServerError, "INTERNAL", "サーバーでエラーが発生しました")
}
