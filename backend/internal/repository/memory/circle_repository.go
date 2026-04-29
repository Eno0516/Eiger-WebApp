package memory

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

var _ domainrepo.CircleRepository = (*CircleRepository)(nil)

type CircleRepository struct {
	mu    sync.RWMutex
	store map[string]*entity.Circle
}

func NewCircleRepository() *CircleRepository {
	return &CircleRepository{store: make(map[string]*entity.Circle)}
}

func (r *CircleRepository) Create(ctx context.Context, circle *entity.Circle) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if circle.ID == "" {
		circle.ID = uuid.NewString()
	}
	now := time.Now()
	circle.CreatedAt = now
	circle.UpdatedAt = now
	r.store[circle.ID] = circle
	return nil
}

func (r *CircleRepository) FindByID(ctx context.Context, id string) (*entity.Circle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.store[id]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	return c, nil
}

func (r *CircleRepository) FindAll(ctx context.Context) ([]*entity.Circle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*entity.Circle, 0, len(r.store))
	for _, c := range r.store {
		result = append(result, c)
	}
	return result, nil
}

func (r *CircleRepository) Update(ctx context.Context, circle *entity.Circle) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[circle.ID]; !ok {
		return apperrors.ErrNotFound
	}
	circle.UpdatedAt = time.Now()
	r.store[circle.ID] = circle
	return nil
}

func (r *CircleRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return apperrors.ErrNotFound
	}
	delete(r.store, id)
	return nil
}
