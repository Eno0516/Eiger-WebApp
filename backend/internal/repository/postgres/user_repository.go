package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/backend/internal/infrastructure/database/sqlcgen"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type UserRepository struct {
	q *sqlcgen.Queries
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{q: sqlcgen.New(db)}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	// TODO: q.CreateUser を呼び出して user に ID / CreatedAt / UpdatedAt を詰める
	// TODO: pgx の unique violation エラーを apperrors.ErrConflict に変換する
	panic("not implemented")
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	// TODO: q.GetUserByID を呼び出す
	// TODO: pgx.ErrNoRows を apperrors.ErrNotFound に変換する
	// TODO: toEntity() でドメインエンティティに変換して返す
	panic("not implemented")
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	// TODO: q.GetUserByEmail を呼び出す
	// TODO: pgx.ErrNoRows を apperrors.ErrNotFound に変換する
	// TODO: toEntity() でドメインエンティティに変換して返す
	panic("not implemented")
}

// FindByVerificationToken はメール認証トークンからユーザーを取得する。
// 期限切れのトークンは DB クエリレベルで除外される。
func (r *UserRepository) FindByVerificationToken(ctx context.Context, token string) (*entity.User, error) {
	// TODO: q.GetUserByVerificationToken を呼び出す
	// TODO: pgx.ErrNoRows を apperrors.ErrNotFound に変換する
	// TODO: toEntity() でドメインエンティティに変換して返す
	panic("not implemented")
}

// FindByResetToken はパスワードリセットトークンからユーザーを取得する。
// 期限切れのトークンは DB クエリレベルで除外される。
func (r *UserRepository) FindByResetToken(ctx context.Context, token string) (*entity.User, error) {
	// TODO: q.GetUserByResetToken を呼び出す
	// TODO: pgx.ErrNoRows を apperrors.ErrNotFound に変換する
	// TODO: toEntity() でドメインエンティティに変換して返す
	panic("not implemented")
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
	// TODO: q.ListUsers を呼び出す
	// TODO: 結果を []entity.User にマップして返す
	panic("not implemented")
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	// TODO: q.UpdateUser を呼び出す
	// TODO: pgx.ErrNoRows を apperrors.ErrNotFound に変換する
	// TODO: 返却された updated_at を user に詰める
	panic("not implemented")
}

// UpdateRole は指定ユーザーのロールを更新する。admin 専用操作。
func (r *UserRepository) UpdateRole(ctx context.Context, id string, role entity.UserRole) error {
	// TODO: q.UpdateUserRole を呼び出す
	// TODO: pgx.ErrNoRows を apperrors.ErrNotFound に変換する
	panic("not implemented")
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	// TODO: q.DeleteUser を呼び出す
	// TODO: RowsAffected() == 0 なら apperrors.ErrNotFound を返す
	panic("not implemented")
}

// toEntity は sqlcgen.User をドメインエンティティに変換する。
func toEntity(row sqlcgen.User) *entity.User {
	// TODO: 全フィールドをマップする
	// TODO: nullable フィールド（*string / *time.Time）はそのままポインタで渡す
	panic("not implemented")
}

// Compile-time check: UserRepository は domain.UserRepository を満たすこと。
var _ interface {
	Create(context.Context, *entity.User) error
	FindByID(context.Context, string) (*entity.User, error)
	FindByEmail(context.Context, string) (*entity.User, error)
	FindByVerificationToken(context.Context, string) (*entity.User, error)
	FindByResetToken(context.Context, string) (*entity.User, error)
	FindAll(context.Context) ([]*entity.User, error)
	Update(context.Context, *entity.User) error
	UpdateRole(context.Context, string, entity.UserRole) error
	Delete(context.Context, string) error
} = (*UserRepository)(nil)

// Ensure apperrors is used (import guard).
var _ = apperrors.ErrNotFound
