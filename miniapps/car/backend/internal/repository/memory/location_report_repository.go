package memory

import (
	"context"
	"sync"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
)

type LocationReportRepository struct {
	mu    sync.RWMutex
	byCar map[string][]entity.LocationReport
}

func NewLocationReportRepository() *LocationReportRepository {
	return &LocationReportRepository{byCar: map[string][]entity.LocationReport{}}
}

func (r *LocationReportRepository) Save(_ context.Context, rep entity.LocationReport) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byCar[rep.CarID] = append(r.byCar[rep.CarID], rep)
	return nil
}

// LatestByCar は申告時刻が最も新しいものを返す。同時刻なら、後から保存したものを返す。
func (r *LocationReportRepository) LatestByCar(_ context.Context, carID string) (entity.LocationReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	reports := r.byCar[carID]
	if len(reports) == 0 {
		return entity.LocationReport{}, domain.ErrNotFound
	}
	latest := reports[0]
	for _, rep := range reports[1:] {
		if !rep.ReportedAt.Before(latest.ReportedAt) {
			latest = rep
		}
	}
	return latest, nil
}
