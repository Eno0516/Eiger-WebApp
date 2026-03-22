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
