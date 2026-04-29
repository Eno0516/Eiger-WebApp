package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/internal/infrastructure/database/sqlcgen"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

// Compile-time check: UserRepository は domain の UserRepository interface を満たす
var _ domainrepo.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	q *sqlcgen.Queries
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{q: sqlcgen.New(db)}
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	row, err := r.q.GetUserByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toEntity(row), nil
}

func (r *UserRepository) FindByLineUserID(ctx context.Context, lineUserID string) (*entity.User, error) {
	row, err := r.q.GetUserByLineUserID(ctx, lineUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toEntity(row), nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
	rows, err := r.q.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]*entity.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, toEntity(row))
	}
	return users, nil
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	row, err := r.q.CreateUser(ctx, sqlcgen.CreateUserParams{
		LineUserID:    user.LineUserID,
		IsSystemAdmin: user.IsSystemAdmin,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return nil, apperrors.ErrConflict
		}
		return nil, err
	}
	return toEntity(row), nil
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	uid, err := parseUUID(user.ID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	row, err := r.q.UpdateUser(ctx, sqlcgen.UpdateUserParams{
		IsRegistered: user.IsRegistered,
		ID:           uid,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toEntity(row), nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	uid, err := parseUUID(id)
	if err != nil {
		return apperrors.ErrNotFound
	}
	return r.q.DeleteUser(ctx, uid)
}

// toEntity は sqlcgen.User をドメインエンティティに変換する。
func toEntity(row sqlcgen.User) *entity.User {
	return &entity.User{
		ID:            uuidToString(row.ID),
		LineUserID:    row.LineUserID,
		IsSystemAdmin: row.IsSystemAdmin,
		IsRegistered:  row.IsRegistered,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}
}

// uuidToString は pgtype.UUID を文字列に変換する。
func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	b := u.Bytes
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// parseUUID は文字列を pgtype.UUID に変換する。
func parseUUID(s string) (pgtype.UUID, error) {
	var uid pgtype.UUID
	err := uid.Scan(s)
	return uid, err
}
