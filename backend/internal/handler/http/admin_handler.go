package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/internal/handler/http/dto"
	"github.com/enoro/eiger-webapp/backend/internal/handler/http/middleware"
	"github.com/enoro/eiger-webapp/backend/internal/usecase/auth"
	"github.com/enoro/eiger-webapp/backend/internal/usecase/circle"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

// AdminHandler は system_admin・circle_admin 向けの管理操作を扱う。
type AdminHandler struct {
	createCircleUC     *circle.CreateCircleUseCase
	listCirclesUC      *circle.ListCirclesUseCase
	deleteCircleUC     *circle.DeleteCircleUseCase
	getCircleMembersUC *circle.GetCircleMembersUseCase
	updateMemberRoleUC *auth.UpdateCircleMemberRoleUseCase
	circleRepo         domainrepo.CircleRepository
	log                zerolog.Logger
}

func NewAdminHandler(
	createCircleUC *circle.CreateCircleUseCase,
	listCirclesUC *circle.ListCirclesUseCase,
	deleteCircleUC *circle.DeleteCircleUseCase,
	getCircleMembersUC *circle.GetCircleMembersUseCase,
	updateMemberRoleUC *auth.UpdateCircleMemberRoleUseCase,
	circleRepo domainrepo.CircleRepository,
	log zerolog.Logger,
) *AdminHandler {
	return &AdminHandler{
		createCircleUC:     createCircleUC,
		listCirclesUC:      listCirclesUC,
		deleteCircleUC:     deleteCircleUC,
		getCircleMembersUC: getCircleMembersUC,
		updateMemberRoleUC: updateMemberRoleUC,
		circleRepo:         circleRepo,
		log:                log,
	}
}

// GetCircles GET /admin/circles
// system_admin: 全件、circle_admin: 自分が所属するサークルのみ
func (h *AdminHandler) GetCircles(w http.ResponseWriter, r *http.Request) {
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

// CreateCircle POST /admin/circles （system_admin のみ）
func (h *AdminHandler) CreateCircle(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCircleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	c, err := h.createCircleUC.Execute(r.Context(), circle.CreateCircleInput{
		Name:               req.Name,
		Description:        req.Description,
		ApplicationFormURL: req.ApplicationFormURL,
		SpreadsheetID:      req.SpreadsheetID,
	})
	if err != nil {
		if apperrors.HTTPStatusCode(err) >= 500 {
			h.log.Error().Err(err).Msg("create circle error")
		}
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	h.writeJSON(w, http.StatusCreated, dto.CircleFromEntity(c))
}

// GetCircle GET /admin/circles/:id
func (h *AdminHandler) GetCircle(w http.ResponseWriter, r *http.Request) {
	circleID := chi.URLParam(r, "id")
	c, err := h.circleRepo.FindByID(r.Context(), circleID)
	if err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	h.writeJSON(w, http.StatusOK, dto.CircleFromEntity(c))
}

// DeleteCircle DELETE /admin/circles/:id （system_admin のみ）
func (h *AdminHandler) DeleteCircle(w http.ResponseWriter, r *http.Request) {
	circleID := chi.URLParam(r, "id")
	if err := h.deleteCircleUC.Execute(r.Context(), circle.DeleteCircleInput{CircleID: circleID}); err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetCircleMembers GET /admin/circles/:id/members
func (h *AdminHandler) GetCircleMembers(w http.ResponseWriter, r *http.Request) {
	circleID := chi.URLParam(r, "id")
	members, err := h.getCircleMembersUC.Execute(r.Context(), circle.GetCircleMembersInput{CircleID: circleID})
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

// UpdateCircleMemberRole PATCH /admin/circles/:id/members/:userId/role
func (h *AdminHandler) UpdateCircleMemberRole(w http.ResponseWriter, r *http.Request) {
	circleID := chi.URLParam(r, "id")
	userID := chi.URLParam(r, "userId")

	var req dto.UpdateCircleMemberRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.updateMemberRoleUC.Execute(r.Context(), auth.UpdateCircleMemberRoleInput{
		CircleID:     circleID,
		TargetUserID: userID,
		NewRole:      entity.CircleRole(req.Role),
	}); err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	h.writeJSON(w, http.StatusOK, dto.MessageResponse{Message: "role updated"})
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
