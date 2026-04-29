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

// Compile-time check: CircleRepository は domain の CircleRepository interface を満たす
var _ domainrepo.CircleRepository = (*CircleRepository)(nil)

type CircleRepository struct {
	q *sqlcgen.Queries
}

func NewCircleRepository(db *pgxpool.Pool) *CircleRepository {
	return &CircleRepository{q: sqlcgen.New(db)}
}

func (r *CircleRepository) Create(ctx context.Context, circle *entity.Circle) error {
	_, err := r.q.CreateCircle(ctx, sqlcgen.CreateCircleParams{
		Name:               circle.Name,
		Description:        circle.Description,
		ApplicationFormUrl: circle.ApplicationFormURL,
		SpreadsheetID:      circle.SpreadsheetID,
	})
	return err
}

func (r *CircleRepository) FindByID(ctx context.Context, id string) (*entity.Circle, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	row, err := r.q.GetCircleByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toCircleEntity(row), nil
}

func (r *CircleRepository) FindAll(ctx context.Context) ([]*entity.Circle, error) {
	rows, err := r.q.ListCircles(ctx)
	if err != nil {
		return nil, err
	}
	circles := make([]*entity.Circle, 0, len(rows))
	for _, row := range rows {
		circles = append(circles, toCircleEntity(row))
	}
	return circles, nil
}

func (r *CircleRepository) Update(ctx context.Context, circle *entity.Circle) error {
	uid, err := parseUUID(circle.ID)
	if err != nil {
		return apperrors.ErrNotFound
	}
	_, err = r.q.UpdateCircle(ctx, sqlcgen.UpdateCircleParams{
		Name:               circle.Name,
		Description:        circle.Description,
		ApplicationFormUrl: circle.ApplicationFormURL,
		SpreadsheetID:      circle.SpreadsheetID,
		ID:                 uid,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.ErrNotFound
	}
	return err
}

func (r *CircleRepository) Delete(ctx context.Context, id string) error {
	uid, err := parseUUID(id)
	if err != nil {
		return apperrors.ErrNotFound
	}
	return r.q.DeleteCircle(ctx, uid)
}

func toCircleEntity(row sqlcgen.Circle) *entity.Circle {
	return &entity.Circle{
		ID:                 uuidToString(row.ID),
		Name:               row.Name,
		Description:        row.Description,
		ApplicationFormURL: row.ApplicationFormUrl,
		SpreadsheetID:      row.SpreadsheetID,
		CreatedAt:          row.CreatedAt.Time,
		UpdatedAt:          row.UpdatedAt.Time,
	}
}
