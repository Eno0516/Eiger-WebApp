package memory_test

import (
	"testing"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/repository/memory"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/repository/repotest"
)

func TestReservationRepository_契約(t *testing.T) {
	repotest.RunReservationRepositoryContract(t, func() repository.ReservationRepository {
		return memory.NewReservationRepository()
	})
}
