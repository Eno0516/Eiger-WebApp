package memory

import (
	"context"
	"sync"
	"time"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

// コンパイル時に UserProfileRepository を実装していることを保証する
var _ domainrepo.UserProfileRepository = (*UserProfileRepository)(nil)

type UserProfileRepository struct {
	mu    sync.RWMutex
	store map[string]*entity.UserProfile // key: userID
}

func NewUserProfileRepository() *UserProfileRepository {
	return &UserProfileRepository{
		store: make(map[string]*entity.UserProfile),
	}
}

func (r *UserProfileRepository) FindByUserID(_ context.Context, userID string) (*entity.UserProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.store[userID]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	return p, nil
}

func (r *UserProfileRepository) Create(_ context.Context, profile *entity.UserProfile) (*entity.UserProfile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[profile.UserID]; ok {
		return nil, apperrors.ErrConflict
	}
	now := time.Now()
	profile.CreatedAt = now
	profile.UpdatedAt = now
	r.store[profile.UserID] = profile
	return profile, nil
}

func (r *UserProfileRepository) Update(_ context.Context, profile *entity.UserProfile) (*entity.UserProfile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[profile.UserID]; !ok {
		return nil, apperrors.ErrNotFound
	}
	profile.UpdatedAt = time.Now()
	r.store[profile.UserID] = profile
	return profile, nil
}
