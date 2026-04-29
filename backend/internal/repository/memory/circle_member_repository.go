package memory

import (
	"context"
	"sync"
	"time"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

var _ domainrepo.CircleMemberRepository = (*CircleMemberRepository)(nil)

type CircleMemberRepository struct {
	mu    sync.RWMutex
	store map[string]*entity.CircleMember // key: circleID + ":" + userID
}

func NewCircleMemberRepository() *CircleMemberRepository {
	return &CircleMemberRepository{store: make(map[string]*entity.CircleMember)}
}

func key(circleID, userID string) string { return circleID + ":" + userID }

func (r *CircleMemberRepository) Upsert(ctx context.Context, member *entity.CircleMember) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if member.JoinedAt.IsZero() {
		member.JoinedAt = time.Now()
	}
	r.store[key(member.CircleID, member.UserID)] = member
	return nil
}

func (r *CircleMemberRepository) FindByCircleAndUser(ctx context.Context, circleID, userID string) (*entity.CircleMember, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.store[key(circleID, userID)]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	return m, nil
}

func (r *CircleMemberRepository) FindByCircle(ctx context.Context, circleID string) ([]*entity.CircleMember, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.CircleMember
	for _, m := range r.store {
		if m.CircleID == circleID {
			result = append(result, m)
		}
	}
	return result, nil
}

func (r *CircleMemberRepository) FindByUser(ctx context.Context, userID string) ([]*entity.CircleMember, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.CircleMember
	for _, m := range r.store {
		if m.UserID == userID {
			result = append(result, m)
		}
	}
	return result, nil
}

func (r *CircleMemberRepository) FindAdminsByCircle(ctx context.Context, circleID string) ([]*entity.CircleMember, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.CircleMember
	for _, m := range r.store {
		if m.CircleID == circleID && m.Role == entity.CircleRoleAdmin {
			result = append(result, m)
		}
	}
	return result, nil
}

func (r *CircleMemberRepository) Delete(ctx context.Context, circleID, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := key(circleID, userID)
	if _, ok := r.store[k]; !ok {
		return apperrors.ErrNotFound
	}
	delete(r.store, k)
	return nil
}
