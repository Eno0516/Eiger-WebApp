package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/internal/infrastructure/database/sqlcgen"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

// Compile-time check
var _ domainrepo.UserProfileRepository = (*UserProfileRepository)(nil)

type UserProfileRepository struct {
	q *sqlcgen.Queries
}

func NewUserProfileRepository(db *pgxpool.Pool) *UserProfileRepository {
	return &UserProfileRepository{q: sqlcgen.New(db)}
}

func (r *UserProfileRepository) FindByUserID(ctx context.Context, userID string) (*entity.UserProfile, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	row, err := r.q.GetUserProfileByUserID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toProfileEntity(row), nil
}

func (r *UserProfileRepository) Create(ctx context.Context, profile *entity.UserProfile) (*entity.UserProfile, error) {
	uid, err := parseUUID(profile.UserID)
	if err != nil {
		return nil, apperrors.NewBadRequest("invalid user id", err)
	}
	row, err := r.q.CreateUserProfile(ctx, sqlcgen.CreateUserProfileParams{
		UserID:      uid,
		DisplayName: profile.DisplayName,
		University:  nullableString(profile.University),
		Faculty:     nullableString(profile.Faculty),
		Grade:       nullableInt16(profile.Grade),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, apperrors.ErrConflict
		}
		return nil, err
	}
	return toProfileEntity(row), nil
}

func (r *UserProfileRepository) Update(ctx context.Context, profile *entity.UserProfile) (*entity.UserProfile, error) {
	uid, err := parseUUID(profile.UserID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	row, err := r.q.UpdateUserProfile(ctx, sqlcgen.UpdateUserProfileParams{
		DisplayName: profile.DisplayName,
		University:  nullableString(profile.University),
		Faculty:     nullableString(profile.Faculty),
		Grade:       nullableInt16(profile.Grade),
		UserID:      uid,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toProfileEntity(row), nil
}

// toProfileEntity は sqlcgen.UserProfile をドメインエンティティに変換する。
func toProfileEntity(row sqlcgen.UserProfile) *entity.UserProfile {
	var grade int
	if row.Grade != nil {
		grade = int(*row.Grade)
	}
	return &entity.UserProfile{
		UserID:      uuidToString(row.UserID),
		DisplayName: row.DisplayName,
		University:  derefString(row.University),
		Faculty:     derefString(row.Faculty),
		Grade:       grade,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}

// nullableString は空文字列を nil に変換する。
func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// nullableInt16 は 0 を nil に変換する。
func nullableInt16(n int) *int16 {
	if n == 0 {
		return nil
	}
	v := int16(n)
	return &v
}

// derefString はポインタが nil の場合に空文字列を返す。
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
