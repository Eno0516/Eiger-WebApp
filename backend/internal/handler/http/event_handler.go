package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/internal/handler/http/dto"
	"github.com/enoro/eiger-webapp/backend/internal/handler/http/middleware"
	"github.com/enoro/eiger-webapp/backend/internal/usecase/event"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

// EventHandler はイベント関連の操作を扱う。
// admin ルート（circle_admin / system_admin）と user ルート（apply）の両方を担当する。
type EventHandler struct {
	createUC        *event.CreateEventUseCase
	listUC          *event.ListEventsUseCase
	applyUC         *event.ApplyEventUseCase
	updateStatusUC  *event.UpdateParticipantStatusUseCase
	bulkUpdateUC    *event.BulkUpdateParticipantsUseCase
	confirmUC       *event.ConfirmEventUseCase
	participantRepo domainrepo.EventParticipantRepository
	log             zerolog.Logger
}

func NewEventHandler(
	createUC *event.CreateEventUseCase,
	listUC *event.ListEventsUseCase,
	applyUC *event.ApplyEventUseCase,
	updateStatusUC *event.UpdateParticipantStatusUseCase,
	bulkUpdateUC *event.BulkUpdateParticipantsUseCase,
	confirmUC *event.ConfirmEventUseCase,
	participantRepo domainrepo.EventParticipantRepository,
	log zerolog.Logger,
) *EventHandler {
	return &EventHandler{
		createUC:        createUC,
		listUC:          listUC,
		applyUC:         applyUC,
		updateStatusUC:  updateStatusUC,
		bulkUpdateUC:    bulkUpdateUC,
		confirmUC:       confirmUC,
		participantRepo: participantRepo,
		log:             log,
	}
}

// ListEvents GET /admin/circles/:id/events
func (h *EventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	circleID := chi.URLParam(r, "id")
	events, err := h.listUC.Execute(r.Context(), event.ListEventsInput{CircleID: circleID})
	if err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	resp := make([]dto.EventResponse, len(events))
	for i, e := range events {
		resp[i] = dto.EventFromEntity(e)
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// CreateEvent POST /admin/circles/:id/events
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	circleID := chi.URLParam(r, "id")
	var req dto.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid date format (YYYY-MM-DD)")
		return
	}
	ev, err := h.createUC.Execute(r.Context(), event.CreateEventInput{
		CircleID:     circleID,
		Title:        req.Title,
		Date:         date,
		Description:  req.Description,
		FormID:       req.FormID,
		LiffURL:      req.LiffURL,
		LineGroupURL: req.LineGroupURL,
	})
	if err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	h.writeJSON(w, http.StatusCreated, dto.EventFromEntity(ev))
}

// GetParticipants GET /admin/circles/:id/events/:eventId/participants
func (h *EventHandler) GetParticipants(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventId")
	participants, err := h.participantRepo.FindByEvent(r.Context(), eventID)
	if err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	resp := make([]dto.EventParticipantResponse, len(participants))
	for i, p := range participants {
		resp[i] = dto.EventParticipantFromEntity(p)
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// UpdateParticipantStatus PATCH /admin/circles/:id/events/:eventId/participants/:userId
func (h *EventHandler) UpdateParticipantStatus(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventId")
	userID := chi.URLParam(r, "userId")
	var req dto.UpdateParticipantStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.updateStatusUC.Execute(r.Context(), event.UpdateParticipantStatusInput{
		EventID: eventID,
		UserID:  userID,
		Status:  entity.ParticipantStatus(req.Status),
	}); err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	h.writeJSON(w, http.StatusOK, dto.MessageResponse{Message: "status updated"})
}

// BulkUpdateParticipants PATCH /admin/circles/:id/events/:eventId/participants
func (h *EventHandler) BulkUpdateParticipants(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventId")
	var req dto.BulkUpdateParticipantsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	items := make([]event.BulkParticipantUpdate, len(req.Participants))
	for i, p := range req.Participants {
		items[i] = event.BulkParticipantUpdate{
			UserID: p.UserID,
			Status: entity.ParticipantStatus(p.Status),
		}
	}
	if err := h.bulkUpdateUC.Execute(r.Context(), event.BulkUpdateParticipantsInput{
		EventID:     eventID,
		Participant: items,
	}); err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	h.writeJSON(w, http.StatusOK, dto.MessageResponse{Message: "participants updated"})
}

// ConfirmEvent POST /admin/circles/:id/events/:eventId/confirm
func (h *EventHandler) ConfirmEvent(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventId")
	if err := h.confirmUC.Execute(r.Context(), event.ConfirmEventInput{EventID: eventID}); err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	h.writeJSON(w, http.StatusOK, dto.MessageResponse{Message: "event confirmed"})
}

// ApplyEvent POST /circles/:circleId/events/:eventId/apply （user ルート）
func (h *EventHandler) ApplyEvent(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventId")
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := h.applyUC.Execute(r.Context(), event.ApplyEventInput{
		EventID: eventID,
		UserID:  claims.UserID,
	}); err != nil {
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	h.writeJSON(w, http.StatusOK, dto.MessageResponse{Message: "applied"})
}

func (h *EventHandler) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.log.Error().Err(err).Msg("encode response")
	}
}

func (h *EventHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, dto.ErrorResponse{Error: message})
}

// import guard
var _ = apperrors.HTTPStatusCode
