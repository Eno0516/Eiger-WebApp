package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/internal/infrastructure/database/sqlcgen"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

// Compile-time check: EventRepository は domain の EventRepository interface を満たす
var _ domainrepo.EventRepository = (*EventRepository)(nil)

type EventRepository struct {
	q *sqlcgen.Queries
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{q: sqlcgen.New(db)}
}

func (r *EventRepository) Create(ctx context.Context, event *entity.Event) error {
	cid, err := parseUUID(event.CircleID)
	if err != nil {
		return apperrors.NewBadRequest("invalid circle id", err)
	}
	_, err = r.q.CreateEvent(ctx, sqlcgen.CreateEventParams{
		CircleID:     cid,
		Title:        event.Title,
		Date:         timeToPgDate(event.Date),
		Description:  event.Description,
		FormID:       nullableUUID(event.FormID),
		LiffUrl:      event.LiffURL,
		LineGroupUrl: event.LineGroupURL,
	})
	return err
}

func (r *EventRepository) FindByID(ctx context.Context, id string) (*entity.Event, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	row, err := r.q.GetEventByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toEventEntity(row), nil
}

func (r *EventRepository) FindByCircle(ctx context.Context, circleID string) ([]*entity.Event, error) {
	cid, err := parseUUID(circleID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	rows, err := r.q.ListEventsByCircle(ctx, cid)
	if err != nil {
		return nil, err
	}
	events := make([]*entity.Event, 0, len(rows))
	for _, row := range rows {
		events = append(events, toEventEntity(row))
	}
	return events, nil
}

func (r *EventRepository) Update(ctx context.Context, event *entity.Event) error {
	uid, err := parseUUID(event.ID)
	if err != nil {
		return apperrors.ErrNotFound
	}
	_, err = r.q.UpdateEvent(ctx, sqlcgen.UpdateEventParams{
		Title:        event.Title,
		Date:         timeToPgDate(event.Date),
		Description:  event.Description,
		FormID:       nullableUUID(event.FormID),
		LiffUrl:      event.LiffURL,
		LineGroupUrl: event.LineGroupURL,
		ConfirmedAt:  nullableTimestamptz(event.ConfirmedAt),
		ID:           uid,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.ErrNotFound
	}
	return err
}

func (r *EventRepository) Delete(ctx context.Context, id string) error {
	uid, err := parseUUID(id)
	if err != nil {
		return apperrors.ErrNotFound
	}
	return r.q.DeleteEvent(ctx, uid)
}

func toEventEntity(row sqlcgen.Event) *entity.Event {
	var formID *string
	if row.FormID != nil {
		s := uuidToString(*row.FormID)
		formID = &s
	}
	var confirmedAt *time.Time
	if row.ConfirmedAt != nil && row.ConfirmedAt.Valid {
		t := row.ConfirmedAt.Time
		confirmedAt = &t
	}
	return &entity.Event{
		ID:           uuidToString(row.ID),
		CircleID:     uuidToString(row.CircleID),
		Title:        row.Title,
		Date:         row.Date.Time,
		Description:  row.Description,
		FormID:       formID,
		LiffURL:      row.LiffUrl,
		LineGroupURL: row.LineGroupUrl,
		ConfirmedAt:  confirmedAt,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}

func timeToPgDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: !t.IsZero()}
}

func nullableUUID(s *string) *pgtype.UUID {
	if s == nil {
		return nil
	}
	uid, err := parseUUID(*s)
	if err != nil {
		return nil
	}
	return &uid
}

func nullableTimestamptz(t *time.Time) *pgtype.Timestamptz {
	if t == nil {
		return nil
	}
	v := pgtype.Timestamptz{Time: *t, Valid: true}
	return &v
}
