package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/backend/internal/handler/http/dto"
	"github.com/enoro/eiger-webapp/backend/internal/usecase/auth"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type AdminHandler struct {
	updateRoleUC *auth.UpdateUserRoleUseCase
	log          zerolog.Logger
}

func NewAdminHandler(updateRoleUC *auth.UpdateUserRoleUseCase, log zerolog.Logger) *AdminHandler {
	return &AdminHandler{updateRoleUC: updateRoleUC, log: log}
}

// UpdateUserRole PATCH /admin/users/:id/role
func (h *AdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// TODO: updateRoleUC.Execute を呼び出す
	// TODO: apperrors.HTTPStatusCode でエラーハンドリングする
	// TODO: 200 + メッセージを返す
	_ = id
	panic("not implemented")
}

func (h *AdminHandler) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.log.Error().Err(err).Msg("encode response")
	}
}

func (h *AdminHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, dto.ErrorResponse{Error: message})
}

// import guard
var _ = apperrors.HTTPStatusCode
var _ = entity.RoleAdmin
