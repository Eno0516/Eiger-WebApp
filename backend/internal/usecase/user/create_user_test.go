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
		uc := user.NewCreateUserUseCase(repo)

		out, err := uc.Execute(context.Background(), user.CreateUserInput{
			Name:  "Alice",
			Email: "alice@example.com",
		})

		require.NoError(t, err)
		assert.Equal(t, "Alice", out.User.Name)
		assert.Equal(t, "alice@example.com", out.User.Email)
		assert.NotEmpty(t, out.User.ID)
	})

	t.Run("validation error - missing name", func(t *testing.T) {
		repo := memory.NewUserRepository()
		uc := user.NewCreateUserUseCase(repo)

		_, err := uc.Execute(context.Background(), user.CreateUserInput{
			Name:  "",
			Email: "alice@example.com",
		})

		assert.ErrorContains(t, err, entity.ErrUserNameRequired.Error())
	})

	t.Run("conflict - duplicate email", func(t *testing.T) {
		repo := memory.NewUserRepository()
		uc := user.NewCreateUserUseCase(repo)

		input := user.CreateUserInput{Name: "Alice", Email: "alice@example.com"}
		_, err := uc.Execute(context.Background(), input)
		require.NoError(t, err)

		_, err = uc.Execute(context.Background(), input)
		assert.Error(t, err)
	})
}
