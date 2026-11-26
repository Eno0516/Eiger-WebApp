package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/valueobject"
)

type ReservationRepository struct {
	mu    sync.RWMutex
	items map[string]entity.Reservation
}

func NewReservationRepository() *ReservationRepository {
	return &ReservationRepository{items: map[string]entity.Reservation{}}
}

// Save は、他の予約と時間帯が重なる場合に domain.ErrReservationConflict を返す。
// postgres 実装(#9)では、DB の排他制約が同じ役割を担う。
func (r *ReservationRepository) Save(_ context.Context, res entity.Reservation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if res.Status.OccupiesCar() {
		for _, other := range r.items {
			if other.ID != res.ID && other.CarID == res.CarID && other.Status.OccupiesCar() && other.Period.Overlaps(res.Period) {
				return domain.ErrReservationConflict
			}
		}
	}
	r.items[res.ID] = res
	return nil
}

func (r *ReservationRepository) FindByID(_ context.Context, id string) (entity.Reservation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res, ok := r.items[id]
	if !ok {
		return entity.Reservation{}, domain.ErrNotFound
	}
	return res, nil
}

func (r *ReservationRepository) ListByCarInRange(_ context.Context, carID string, period valueobject.TimeRange) ([]entity.Reservation, error) {
	return r.list(func(res entity.Reservation) bool {
		return res.CarID == carID && res.Period.Overlaps(period)
	}), nil
}

func (r *ReservationRepository) ListByCircleInRange(_ context.Context, circleID string, period valueobject.TimeRange) ([]entity.Reservation, error) {
	return r.list(func(res entity.Reservation) bool {
		return res.CircleID == circleID && res.Period.Overlaps(period)
	}), nil
}

func (r *ReservationRepository) ListByMember(_ context.Context, circleID, memberID string) ([]entity.Reservation, error) {
	return r.list(func(res entity.Reservation) bool {
		return res.CircleID == circleID && res.MemberID == memberID
	}), nil
}

// list は条件に合う予約を、開始の早い順(同時刻は ID 順)で返す。
func (r *ReservationRepository) list(match func(entity.Reservation) bool) []entity.Reservation {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []entity.Reservation
	for _, res := range r.items {
		if match(res) {
			out = append(out, res)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if !out[i].Period.Start().Equal(out[j].Period.Start()) {
			return out[i].Period.Start().Before(out[j].Period.Start())
		}
		return out[i].ID < out[j].ID
	})
	return out
}
