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

var _ domainrepo.EventRepository = (*EventRepository)(nil)

type EventRepository struct {
	mu    sync.RWMutex
	store map[string]*entity.Event
}

func NewEventRepository() *EventRepository {
	return &EventRepository{store: make(map[string]*entity.Event)}
}

func (r *EventRepository) Create(ctx context.Context, ev *entity.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ev.ID == "" {
		ev.ID = uuid.NewString()
	}
	ev.CreatedAt = time.Now()
	r.store[ev.ID] = ev
	return nil
}

func (r *EventRepository) FindByID(ctx context.Context, id string) (*entity.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ev, ok := r.store[id]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	return ev, nil
}

func (r *EventRepository) FindByCircle(ctx context.Context, circleID string) ([]*entity.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Event
	for _, ev := range r.store {
		if ev.CircleID == circleID {
			result = append(result, ev)
		}
	}
	return result, nil
}

func (r *EventRepository) Update(ctx context.Context, ev *entity.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[ev.ID]; !ok {
		return apperrors.ErrNotFound
	}
	r.store[ev.ID] = ev
	return nil
}

func (r *EventRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return apperrors.ErrNotFound
	}
	delete(r.store, id)
	return nil
}
