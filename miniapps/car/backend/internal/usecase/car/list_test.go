package car_test

import (
	"context"
	"errors"
	"testing"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/testutil"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/car"
)

func TestList(t *testing.T) {
	ctx := context.Background()

	t.Run("メンバーは、自分のサークルの車だけを見られる", func(t *testing.T) {
		e := testutil.NewEnv(t)
		_ = e.Cars.Save(ctx, entity.Car{ID: "car-other", CircleID: "circle-2", Name: "別サークルの車", Available: true})

		got, err := car.NewList(e.Cars, e.Members).Execute(ctx, car.ListInput{CircleID: testutil.CircleID, MemberID: testutil.Member1})
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0].ID != testutil.CarID {
			t.Errorf("car-1 だけのはず: %+v", got)
		}
	})

	t.Run("サークルのメンバーでない人は見られない", func(t *testing.T) {
		e := testutil.NewEnv(t)
		_, err := car.NewList(e.Cars, e.Members).Execute(ctx, car.ListInput{CircleID: testutil.CircleID, MemberID: testutil.Outsider})
		if !errors.Is(err, domain.ErrNotMember) {
			t.Errorf("ErrNotMember を期待: %v", err)
		}
	})
}
