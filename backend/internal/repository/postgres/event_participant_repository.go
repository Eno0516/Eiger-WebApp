package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/internal/infrastructure/database/sqlcgen"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

// Compile-time check: EventParticipantRepository は domain の EventParticipantRepository interface を満たす
var _ domainrepo.EventParticipantRepository = (*EventParticipantRepository)(nil)

type EventParticipantRepository struct {
	q *sqlcgen.Queries
}

func NewEventParticipantRepository(db *pgxpool.Pool) *EventParticipantRepository {
	return &EventParticipantRepository{q: sqlcgen.New(db)}
}

func (r *EventParticipantRepository) Upsert(ctx context.Context, p *entity.EventParticipant) error {
	eid, err := parseUUID(p.EventID)
	if err != nil {
		return apperrors.NewBadRequest("invalid event id", err)
	}
	uid, err := parseUUID(p.UserID)
	if err != nil {
		return apperrors.NewBadRequest("invalid user id", err)
	}
	_, err = r.q.UpsertEventParticipant(ctx, sqlcgen.UpsertEventParticipantParams{
		EventID: eid,
		UserID:  uid,
		Status:  string(p.Status),
	})
	return err
}

func (r *EventParticipantRepository) BulkUpsert(ctx context.Context, participants []*entity.EventParticipant) error {
	for _, p := range participants {
		if err := r.Upsert(ctx, p); err != nil {
			return err
		}
	}
	return nil
}

func (r *EventParticipantRepository) FindByEvent(ctx context.Context, eventID string) ([]*entity.EventParticipant, error) {
	eid, err := parseUUID(eventID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	rows, err := r.q.ListEventParticipantsByEvent(ctx, eid)
	if err != nil {
		return nil, err
	}
	return toParticipantEntities(rows), nil
}

func (r *EventParticipantRepository) FindByEventAndUser(ctx context.Context, eventID, userID string) (*entity.EventParticipant, error) {
	eid, err := parseUUID(eventID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	row, err := r.q.GetEventParticipantByEventAndUser(ctx, sqlcgen.GetEventParticipantByEventAndUserParams{
		EventID: eid,
		UserID:  uid,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toParticipantEntity(row), nil
}

func (r *EventParticipantRepository) FindConfirmedByEvent(ctx context.Context, eventID string) ([]*entity.EventParticipant, error) {
	eid, err := parseUUID(eventID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	rows, err := r.q.ListConfirmedEventParticipantsByEvent(ctx, eid)
	if err != nil {
		return nil, err
	}
	return toParticipantEntities(rows), nil
}

func toParticipantEntity(row sqlcgen.EventParticipant) *entity.EventParticipant {
	return &entity.EventParticipant{
		EventID:   uuidToString(row.EventID),
		UserID:    uuidToString(row.UserID),
		Status:    entity.ParticipantStatus(row.Status),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}

func toParticipantEntities(rows []sqlcgen.EventParticipant) []*entity.EventParticipant {
	result := make([]*entity.EventParticipant, 0, len(rows))
	for _, row := range rows {
		result = append(result, toParticipantEntity(row))
	}
	return result
}
