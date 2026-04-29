package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/backend/internal/repository/memory"
	"github.com/enoro/eiger-webapp/backend/internal/usecase/user"
)

func TestCreateUserUseCase_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := memory.NewUserRepository()
		profileRepo := memory.NewUserProfileRepository()
		uc := user.NewCreateUserUseCase(repo, profileRepo)

		out, err := uc.Execute(context.Background(), user.CreateUserInput{
			LineUserID: "line-user-abc",
		})

		require.NoError(t, err)
		assert.Equal(t, "line-user-abc", out.User.LineUserID)
		assert.NotEmpty(t, out.User.ID)
		assert.False(t, out.User.IsRegistered)
	})

	t.Run("validation error - missing line user id", func(t *testing.T) {
		repo := memory.NewUserRepository()
		profileRepo := memory.NewUserProfileRepository()
		uc := user.NewCreateUserUseCase(repo, profileRepo)

		_, err := uc.Execute(context.Background(), user.CreateUserInput{
			LineUserID: "",
		})

		assert.ErrorContains(t, err, entity.ErrUserLineUserIDRequired.Error())
	})

	t.Run("conflict - duplicate line_user_id", func(t *testing.T) {
		repo := memory.NewUserRepository()
		profileRepo := memory.NewUserProfileRepository()
		uc := user.NewCreateUserUseCase(repo, profileRepo)

		input := user.CreateUserInput{LineUserID: "line-user-abc"}
		_, err := uc.Execute(context.Background(), input)
		require.NoError(t, err)

		_, err = uc.Execute(context.Background(), input)
		assert.Error(t, err)
	})
}
