package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/backend/internal/handler/http/dto"
	"github.com/enoro/eiger-webapp/backend/internal/handler/http/middleware"
	"github.com/enoro/eiger-webapp/backend/internal/usecase/auth"
	"github.com/enoro/eiger-webapp/backend/internal/usecase/circle"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type CircleHandler struct {
	applyUC          *circle.ApplyCircleUseCase
	updateMemberRole *auth.UpdateCircleMemberRoleUseCase
	listCirclesUC    *circle.ListCirclesUseCase
	getMembersUC     *circle.GetCircleMembersUseCase
	log              zerolog.Logger
}

func NewCircleHandler(
	applyUC *circle.ApplyCircleUseCase,
	updateMemberRole *auth.UpdateCircleMemberRoleUseCase,
	listCirclesUC *circle.ListCirclesUseCase,
	getMembersUC *circle.GetCircleMembersUseCase,
	log zerolog.Logger,
) *CircleHandler {
	return &CircleHandler{
		applyUC:          applyUC,
		updateMemberRole: updateMemberRole,
		listCirclesUC:    listCirclesUC,
		getMembersUC:     getMembersUC,
		log:              log,
	}
}

// GetAll GET /circles
func (h *CircleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	circles, err := h.listCirclesUC.Execute(r.Context(), circle.ListCirclesInput{
		RequesterUserID: claims.UserID,
		IsSystemAdmin:   claims.IsSystemAdmin,
	})
	if err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	resp := make([]dto.CircleResponse, len(circles))
	for i, c := range circles {
		resp[i] = dto.CircleFromEntity(c)
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// Apply POST /circles/:id/apply
func (h *CircleHandler) Apply(w http.ResponseWriter, r *http.Request) {
	circleID := chi.URLParam(r, "id")
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := h.applyUC.Execute(r.Context(), circle.ApplyCircleInput{
		CircleID: circleID,
		UserID:   claims.UserID,
	}); err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	h.writeJSON(w, http.StatusOK, dto.MessageResponse{Message: "applied"})
}

// GetMembers GET /circles/:id/members
func (h *CircleHandler) GetMembers(w http.ResponseWriter, r *http.Request) {
	circleID := chi.URLParam(r, "id")
	members, err := h.getMembersUC.Execute(r.Context(), circle.GetCircleMembersInput{CircleID: circleID})
	if err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	resp := make([]dto.CircleMemberResponse, len(members))
	for i, m := range members {
		resp[i] = dto.CircleMemberFromEntity(m)
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// GetFormResponse GET /circles/:id/members/:userId/form-response
func (h *CircleHandler) GetFormResponse(w http.ResponseWriter, r *http.Request) {
	// TODO: Google Sheets API で対象ユーザーの UUID に紐づく Form 回答を取得して返す
	h.writeError(w, http.StatusNotImplemented, "not implemented")
}

// UpdateMemberRole PATCH /circles/:id/members/:userId/role
func (h *CircleHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	circleID := chi.URLParam(r, "id")
	userID := chi.URLParam(r, "userId")

	var req dto.UpdateCircleMemberRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.updateMemberRole.Execute(r.Context(), auth.UpdateCircleMemberRoleInput{
		CircleID:     circleID,
		TargetUserID: userID,
		NewRole:      entity.CircleRole(req.Role),
	}); err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	h.writeJSON(w, http.StatusOK, dto.MessageResponse{Message: "role updated"})
}

func (h *CircleHandler) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.log.Error().Err(err).Msg("encode response")
	}
}

func (h *CircleHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, dto.ErrorResponse{Error: message})
}
