package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type UserRepository struct {
	mu    sync.RWMutex
	store map[string]*entity.User
	seq   int
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		store: make(map[string]*entity.User),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.seq++
	user.ID = fmt.Sprintf("%d", r.seq)
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	r.store[user.ID] = user
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.store[id]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.store {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, apperrors.ErrNotFound
}

// FindByVerificationToken はメール認証トークンからユーザーを検索する（テスト用）。
func (r *UserRepository) FindByVerificationToken(ctx context.Context, token string) (*entity.User, error) {
	// TODO: store を走査して EmailVerificationToken が一致するユーザーを返す
	// TODO: EmailVerificationExpiresAt が過去なら apperrors.ErrNotFound を返す
	panic("not implemented")
}

// FindByResetToken はパスワードリセットトークンからユーザーを検索する（テスト用）。
func (r *UserRepository) FindByResetToken(ctx context.Context, token string) (*entity.User, error) {
	// TODO: store を走査して PasswordResetToken が一致するユーザーを返す
	// TODO: PasswordResetExpiresAt が過去なら apperrors.ErrNotFound を返す
	panic("not implemented")
}

// UpdateRole は指定ユーザーのロールを更新する（テスト用）。
func (r *UserRepository) UpdateRole(ctx context.Context, id string, role entity.UserRole) error {
	// TODO: store から id で取得し Role を更新する
	// TODO: 存在しない場合は apperrors.ErrNotFound を返す
	panic("not implemented")
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*entity.User, 0, len(r.store))
	for _, u := range r.store {
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[user.ID]; !ok {
		return apperrors.ErrNotFound
	}
	user.UpdatedAt = time.Now()
	r.store[user.ID] = user
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[id]; !ok {
		return apperrors.ErrNotFound
	}
	delete(r.store, id)
	return nil
}
