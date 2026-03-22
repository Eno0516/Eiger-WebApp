package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/enoro/eiger-webapp/backend/internal/handler/http/dto"
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

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	out, err := h.createUC.Execute(r.Context(), user.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, dto.UserFromEntity(out.User))
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	out, err := h.getUC.ExecuteByID(r.Context(), id)
	if err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, dto.UserFromEntity(out.User))
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	out, err := h.getUC.ExecuteAll(r.Context())
	if err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), err.Error())
		return
	}

	resp := make([]dto.UserResponse, len(out.Users))
	for i, u := range out.Users {
		resp[i] = dto.UserFromEntity(u)
	}
	h.writeJSON(w, http.StatusOK, resp)
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
