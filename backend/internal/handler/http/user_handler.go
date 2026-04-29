package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/enoro/eiger-webapp/backend/internal/handler/http/dto"
	"github.com/enoro/eiger-webapp/backend/internal/handler/http/middleware"
	"github.com/enoro/eiger-webapp/backend/internal/usecase/user"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type UserHandler struct {
	createUC *user.CreateUserUseCase
	getUC    *user.GetUserUseCase
	log      zerolog.Logger
}

func NewUserHandler(createUC *user.CreateUserUseCase, getUC *user.GetUserUseCase, log zerolog.Logger) *UserHandler {
	return &UserHandler{createUC: createUC, getUC: getUC, log: log}
}

// GetMe GET /users/me
// JWT から userID を取得し、自分のユーザー情報（is_registered など）を返す。
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	out, err := h.getUC.ExecuteByID(r.Context(), claims.UserID)
	if err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	h.writeJSON(w, http.StatusOK, dto.MeResponse{
		ID:            out.User.ID,
		IsRegistered:  out.User.IsRegistered,
		IsSystemAdmin: out.User.IsSystemAdmin,
	})
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.LineAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Create はシステム管理者向けの手動ユーザー作成。LineUserID のみ受け取る。
	out, err := h.createUC.Execute(r.Context(), user.CreateUserInput{LineUserID: req.IDToken})
	if err != nil {
		if apperrors.HTTPStatusCode(err) >= 500 {
			h.log.Error().Err(err).Msg("create user error")
		}
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	h.writeJSON(w, http.StatusCreated, dto.UserFromEntity(out.User))
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	out, err := h.getUC.ExecuteByID(r.Context(), id)
	if err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}

	h.writeJSON(w, http.StatusOK, dto.UserFromEntity(out.User))
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	out, err := h.getUC.ExecuteAll(r.Context())
	if err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}

	resp := make([]dto.UserResponse, len(out.Users))
	for i, u := range out.Users {
		resp[i] = dto.UserFromEntity(u)
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// UpdateDisplayName PATCH /users/:id/display-name
func (h *UserHandler) UpdateDisplayName(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateDisplayNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok || claims.UserID != id {
		h.writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	// TODO: UpdateDisplayNameUseCase を実装後、ここで呼び出す
	h.writeError(w, http.StatusNotImplemented, "not implemented")
}

func (h *UserHandler) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.log.Error().Err(err).Msg("encode response")
	}
}

func (h *UserHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, dto.ErrorResponse{Error: message})
}
