package memory

import (
	"context"
	"sync"
	"time"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

var _ domainrepo.EventParticipantRepository = (*EventParticipantRepository)(nil)

type EventParticipantRepository struct {
	mu    sync.RWMutex
	store map[string]*entity.EventParticipant // key: eventID + ":" + userID
}

func NewEventParticipantRepository() *EventParticipantRepository {
	return &EventParticipantRepository{store: make(map[string]*entity.EventParticipant)}
}

func epKey(eventID, userID string) string { return eventID + ":" + userID }

func (r *EventParticipantRepository) Upsert(ctx context.Context, p *entity.EventParticipant) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	r.store[epKey(p.EventID, p.UserID)] = p
	return nil
}

func (r *EventParticipantRepository) BulkUpsert(ctx context.Context, participants []*entity.EventParticipant) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for _, p := range participants {
		if p.CreatedAt.IsZero() {
			p.CreatedAt = now
		}
		r.store[epKey(p.EventID, p.UserID)] = p
	}
	return nil
}

func (r *EventParticipantRepository) FindByEvent(ctx context.Context, eventID string) ([]*entity.EventParticipant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.EventParticipant
	for _, p := range r.store {
		if p.EventID == eventID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *EventParticipantRepository) FindByEventAndUser(ctx context.Context, eventID, userID string) (*entity.EventParticipant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.store[epKey(eventID, userID)]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	return p, nil
}

func (r *EventParticipantRepository) FindConfirmedByEvent(ctx context.Context, eventID string) ([]*entity.EventParticipant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.EventParticipant
	for _, p := range r.store {
		if p.EventID == eventID && p.Status == entity.ParticipantStatusConfirmed {
			result = append(result, p)
		}
	}
	return result, nil
}
