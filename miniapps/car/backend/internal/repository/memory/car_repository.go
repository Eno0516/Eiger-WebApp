// Package memory は、repository の in-memory 実装。テストと、DB なしの開発用。
package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
)

type CarRepository struct {
	mu    sync.RWMutex
	items map[string]entity.Car
}

func NewCarRepository() *CarRepository {
	return &CarRepository{items: map[string]entity.Car{}}
}

func (r *CarRepository) FindByID(_ context.Context, id string) (entity.Car, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	car, ok := r.items[id]
	if !ok {
		return entity.Car{}, domain.ErrNotFound
	}
	return car, nil
}

func (r *CarRepository) ListByCircle(_ context.Context, circleID string) ([]entity.Car, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []entity.Car
	for _, c := range r.items {
		if c.CircleID == circleID {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (r *CarRepository) Save(_ context.Context, car entity.Car) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[car.ID] = car
	return nil
}
