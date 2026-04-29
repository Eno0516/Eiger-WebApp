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

// Compile-time check: CircleMemberRepository は domain の CircleMemberRepository interface を満たす
var _ domainrepo.CircleMemberRepository = (*CircleMemberRepository)(nil)

type CircleMemberRepository struct {
	q *sqlcgen.Queries
}

func NewCircleMemberRepository(db *pgxpool.Pool) *CircleMemberRepository {
	return &CircleMemberRepository{q: sqlcgen.New(db)}
}

func (r *CircleMemberRepository) Upsert(ctx context.Context, member *entity.CircleMember) error {
	cid, err := parseUUID(member.CircleID)
	if err != nil {
		return apperrors.NewBadRequest("invalid circle id", err)
	}
	uid, err := parseUUID(member.UserID)
	if err != nil {
		return apperrors.NewBadRequest("invalid user id", err)
	}
	_, err = r.q.UpsertCircleMember(ctx, sqlcgen.UpsertCircleMemberParams{
		CircleID: cid,
		UserID:   uid,
		Role:     string(member.Role),
	})
	return err
}

func (r *CircleMemberRepository) FindByCircleAndUser(ctx context.Context, circleID, userID string) (*entity.CircleMember, error) {
	cid, err := parseUUID(circleID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	row, err := r.q.GetCircleMemberByCircleAndUser(ctx, sqlcgen.GetCircleMemberByCircleAndUserParams{
		CircleID: cid,
		UserID:   uid,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toCircleMemberEntity(row), nil
}

func (r *CircleMemberRepository) FindByCircle(ctx context.Context, circleID string) ([]*entity.CircleMember, error) {
	cid, err := parseUUID(circleID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	rows, err := r.q.ListCircleMembersByCircle(ctx, cid)
	if err != nil {
		return nil, err
	}
	return toCircleMemberEntities(rows), nil
}

func (r *CircleMemberRepository) FindByUser(ctx context.Context, userID string) ([]*entity.CircleMember, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	rows, err := r.q.ListCircleMembersByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	return toCircleMemberEntities(rows), nil
}

func (r *CircleMemberRepository) FindAdminsByCircle(ctx context.Context, circleID string) ([]*entity.CircleMember, error) {
	cid, err := parseUUID(circleID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	rows, err := r.q.ListAdminsByCircle(ctx, cid)
	if err != nil {
		return nil, err
	}
	return toCircleMemberEntities(rows), nil
}

func (r *CircleMemberRepository) Delete(ctx context.Context, circleID, userID string) error {
	cid, err := parseUUID(circleID)
	if err != nil {
		return apperrors.ErrNotFound
	}
	uid, err := parseUUID(userID)
	if err != nil {
		return apperrors.ErrNotFound
	}
	return r.q.DeleteCircleMember(ctx, sqlcgen.DeleteCircleMemberParams{
		CircleID: cid,
		UserID:   uid,
	})
}

func toCircleMemberEntity(row sqlcgen.CircleMember) *entity.CircleMember {
	return &entity.CircleMember{
		CircleID: uuidToString(row.CircleID),
		UserID:   uuidToString(row.UserID),
		Role:     entity.CircleRole(row.Role),
		JoinedAt: row.JoinedAt.Time,
	}
}

func toCircleMemberEntities(rows []sqlcgen.CircleMember) []*entity.CircleMember {
	result := make([]*entity.CircleMember, 0, len(rows))
	for _, row := range rows {
		result = append(result, toCircleMemberEntity(row))
	}
	return result
}
