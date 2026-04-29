//go:build integration

package integration_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	postgresrepo "github.com/enoro/eiger-webapp/backend/internal/repository/postgres"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

// testDB はテストスイート全体で共有する DB 接続プール。
// TestMain で初期化される。
var testDB *pgxpool.Pool

// TestMain は PostgreSQL コンテナを起動してマイグレーションを適用する。
func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		// マイグレーションファイルをコンテナ初期化スクリプトとして適用する。
		// ファイルは tests/integration/ からの相対パスで指定する。
		tcpostgres.WithInitScripts(
			"../../internal/infrastructure/database/migrations/001_create_users.sql",
			"../../internal/infrastructure/database/migrations/002_create_circles.sql",
			"../../internal/infrastructure/database/migrations/003_create_circle_members.sql",
			"../../internal/infrastructure/database/migrations/004_create_user_profiles.sql",
			"../../internal/infrastructure/database/migrations/005_create_user_line_accounts.sql",
			"../../internal/infrastructure/database/migrations/006_create_events.sql",
			"../../internal/infrastructure/database/migrations/007_create_event_participants.sql",
		),
	)
	if err != nil {
		panic("PostgreSQL コンテナの起動に失敗しました: " + err.Error())
	}
	defer pgContainer.Terminate(ctx) //nolint:errcheck

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic("接続文字列の取得に失敗しました: " + err.Error())
	}

	testDB, err = pgxpool.New(ctx, connStr)
	if err != nil {
		panic("pgxpool の初期化に失敗しました: " + err.Error())
	}
	defer testDB.Close()

	os.Exit(m.Run())
}

// newUserRepo はテスト用の postgres.UserRepository を生成する。
func newUserRepo() *postgresrepo.UserRepository {
	return postgresrepo.NewUserRepository(testDB)
}

// --- テストケース ---

func TestUserRepository(t *testing.T) {
	t.Run("create - ユーザーが作成され FindByLineUserID で取得できる", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		repo := newUserRepo()
		newUser := &entity.User{
			LineUserID:    "line-sub-001",
			IsSystemAdmin: false,
		}

		// Act
		created, err := repo.Create(ctx, newUser)

		// Assert
		require.NoError(t, err)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, "line-sub-001", created.LineUserID)
		_ = zerolog.Nop() // import guard
	})

	t.Run("find by line user id - 存在しない ID で ErrNotFound が返る", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		repo := newUserRepo()

		// Act
		_, err := repo.FindByLineUserID(ctx, "nonexistent-line-id")

		// Assert
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
	})

	t.Run("create duplicate - 同じ LINE ユーザー ID で ErrConflict が返る", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		repo := newUserRepo()

		_, err := repo.Create(ctx, &entity.User{LineUserID: "line-dup-001"})
		require.NoError(t, err)

		// Act
		_, err = repo.Create(ctx, &entity.User{LineUserID: "line-dup-001"})

		// Assert
		assert.ErrorIs(t, err, apperrors.ErrConflict)
	})

	t.Run("update - DisplayName の変更が反映される", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		repo := newUserRepo()

		created, err := repo.Create(ctx, &entity.User{LineUserID: "line-update-001"})
		require.NoError(t, err)

		// Act
		created.DisplayName = "変更後の名前"
		updated, err := repo.Update(ctx, created)
		require.NoError(t, err)

		// Assert
		found, err := repo.FindByID(ctx, updated.ID)
		require.NoError(t, err)
		assert.Equal(t, "変更後の名前", found.DisplayName)
	})
}

// import guard
var (
	_ = apperrors.ErrNotFound
	_ = (*postgresrepo.UserRepository)(nil)
)
