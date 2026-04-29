//go:build integration

package integration_test

import (
	"context"
	"testing"
	"time"

	postgresrepo "github.com/enoro/eiger-webapp/backend/internal/repository/postgres"
	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newEventRepo() *postgresrepo.EventRepository {
	return postgresrepo.NewEventRepository(testDB)
}

func newEventParticipantRepo() *postgresrepo.EventParticipantRepository {
	return postgresrepo.NewEventParticipantRepository(testDB)
}

// seedTestEvent は統合テスト用イベントを DB に作成して返すヘルパー。
func seedTestEvent(t *testing.T, circleID, title string) *entity.Event {
	t.Helper()
	repo := newEventRepo()
	ev := &entity.Event{
		CircleID: circleID,
		Title:    title,
		Date:     time.Now().Add(24 * time.Hour),
	}
	err := repo.Create(context.Background(), ev)
	require.NoError(t, err)
	// Create は ID を返さないため FindByCircle で取得する
	all, err := repo.FindByCircle(context.Background(), circleID)
	require.NoError(t, err)
	for _, e := range all {
		if e.Title == title {
			return e
		}
	}
	t.Fatalf("seedTestEvent: created event %q not found", title)
	return nil
}

// =============================================================================
// EventRepository
// =============================================================================

func TestEventRepository(t *testing.T) {
	t.Run("create - イベントが作成され FindByCircle で取得できる", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		c := seedTestCircle(t, "イベントテストサークルA")
		repo := newEventRepo()

		// Act
		err := repo.Create(ctx, &entity.Event{
			CircleID: c.ID,
			Title:    "テストイベントA",
			Date:     time.Now().Add(48 * time.Hour),
		})

		// Assert
		require.NoError(t, err)
		events, err := repo.FindByCircle(ctx, c.ID)
		require.NoError(t, err)
		require.Len(t, events, 1)
		assert.Equal(t, "テストイベントA", events[0].Title)
	})

	t.Run("find by id - 作成後に取得できる", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		c := seedTestCircle(t, "イベントテストサークルB")
		ev := seedTestEvent(t, c.ID, "テストイベントB")

		// Act
		found, err := newEventRepo().FindByID(ctx, ev.ID)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "テストイベントB", found.Title)
		assert.Equal(t, c.ID, found.CircleID)
	})

	t.Run("find by id - 存在しない ID で ErrNotFound が返る", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		repo := newEventRepo()

		// Act
		_, err := repo.FindByID(ctx, "00000000-0000-0000-0000-000000000000")

		// Assert
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
	})

	t.Run("update - confirmed_at の設定が反映される", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		c := seedTestCircle(t, "イベントテストサークルC")
		ev := seedTestEvent(t, c.ID, "テストイベントC")
		repo := newEventRepo()

		// Act
		now := time.Now()
		ev.ConfirmedAt = &now
		err := repo.Update(ctx, ev)

		// Assert
		require.NoError(t, err)
		found, err := repo.FindByID(ctx, ev.ID)
		require.NoError(t, err)
		assert.NotNil(t, found.ConfirmedAt)
	})

	t.Run("delete - 削除後に FindByID で ErrNotFound が返る", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		c := seedTestCircle(t, "イベントテストサークルD")
		ev := seedTestEvent(t, c.ID, "テストイベントD")
		repo := newEventRepo()

		// Act
		err := repo.Delete(ctx, ev.ID)

		// Assert
		require.NoError(t, err)
		_, err = repo.FindByID(ctx, ev.ID)
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
	})
}

// =============================================================================
// EventParticipantRepository
// =============================================================================

func TestEventParticipantRepository(t *testing.T) {
	t.Run("upsert - 参加申請が作成され FindByEventAndUser で取得できる", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u := seedTestUser(t, "line-participant-001")
		c := seedTestCircle(t, "参加者テストサークルA")
		ev := seedTestEvent(t, c.ID, "参加者テストイベントA")
		repo := newEventParticipantRepo()

		// Act
		err := repo.Upsert(ctx, &entity.EventParticipant{
			EventID: ev.ID,
			UserID:  u.ID,
			Status:  entity.ParticipantStatusApplied,
		})

		// Assert
		require.NoError(t, err)
		found, err := repo.FindByEventAndUser(ctx, ev.ID, u.ID)
		require.NoError(t, err)
		assert.Equal(t, entity.ParticipantStatusApplied, found.Status)
	})

	t.Run("upsert - ステータス変更が反映される", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u := seedTestUser(t, "line-participant-002")
		c := seedTestCircle(t, "参加者テストサークルB")
		ev := seedTestEvent(t, c.ID, "参加者テストイベントB")
		repo := newEventParticipantRepo()
		require.NoError(t, repo.Upsert(ctx, &entity.EventParticipant{EventID: ev.ID, UserID: u.ID, Status: entity.ParticipantStatusApplied}))

		// Act
		err := repo.Upsert(ctx, &entity.EventParticipant{EventID: ev.ID, UserID: u.ID, Status: entity.ParticipantStatusConfirmed})

		// Assert
		require.NoError(t, err)
		found, err := repo.FindByEventAndUser(ctx, ev.ID, u.ID)
		require.NoError(t, err)
		assert.Equal(t, entity.ParticipantStatusConfirmed, found.Status)
	})

	t.Run("find by event and user - 存在しない組み合わせで ErrNotFound が返る", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u := seedTestUser(t, "line-participant-003")
		c := seedTestCircle(t, "参加者テストサークルC")
		ev := seedTestEvent(t, c.ID, "参加者テストイベントC")
		repo := newEventParticipantRepo()

		// Act
		_, err := repo.FindByEventAndUser(ctx, ev.ID, u.ID)

		// Assert
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
	})

	t.Run("find by event - イベントの全参加者が返る", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u1 := seedTestUser(t, "line-participant-004")
		u2 := seedTestUser(t, "line-participant-005")
		c := seedTestCircle(t, "参加者テストサークルD")
		ev := seedTestEvent(t, c.ID, "参加者テストイベントD")
		repo := newEventParticipantRepo()
		require.NoError(t, repo.Upsert(ctx, &entity.EventParticipant{EventID: ev.ID, UserID: u1.ID, Status: entity.ParticipantStatusApplied}))
		require.NoError(t, repo.Upsert(ctx, &entity.EventParticipant{EventID: ev.ID, UserID: u2.ID, Status: entity.ParticipantStatusConfirmed}))

		// Act
		participants, err := repo.FindByEvent(ctx, ev.ID)

		// Assert
		require.NoError(t, err)
		assert.Len(t, participants, 2)
	})

	t.Run("find confirmed by event - confirmed のみ返る", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u1 := seedTestUser(t, "line-participant-006")
		u2 := seedTestUser(t, "line-participant-007")
		c := seedTestCircle(t, "参加者テストサークルE")
		ev := seedTestEvent(t, c.ID, "参加者テストイベントE")
		repo := newEventParticipantRepo()
		require.NoError(t, repo.Upsert(ctx, &entity.EventParticipant{EventID: ev.ID, UserID: u1.ID, Status: entity.ParticipantStatusApplied}))
		require.NoError(t, repo.Upsert(ctx, &entity.EventParticipant{EventID: ev.ID, UserID: u2.ID, Status: entity.ParticipantStatusConfirmed}))

		// Act
		confirmed, err := repo.FindConfirmedByEvent(ctx, ev.ID)

		// Assert
		require.NoError(t, err)
		assert.Len(t, confirmed, 1)
		assert.Equal(t, entity.ParticipantStatusConfirmed, confirmed[0].Status)
	})

	t.Run("bulk upsert - 複数参加者のステータスが一括更新される", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		u1 := seedTestUser(t, "line-participant-008")
		u2 := seedTestUser(t, "line-participant-009")
		c := seedTestCircle(t, "参加者テストサークルF")
		ev := seedTestEvent(t, c.ID, "参加者テストイベントF")
		repo := newEventParticipantRepo()
		require.NoError(t, repo.Upsert(ctx, &entity.EventParticipant{EventID: ev.ID, UserID: u1.ID, Status: entity.ParticipantStatusApplied}))
		require.NoError(t, repo.Upsert(ctx, &entity.EventParticipant{EventID: ev.ID, UserID: u2.ID, Status: entity.ParticipantStatusApplied}))

		// Act
		err := repo.BulkUpsert(ctx, []*entity.EventParticipant{
			{EventID: ev.ID, UserID: u1.ID, Status: entity.ParticipantStatusConfirmed},
			{EventID: ev.ID, UserID: u2.ID, Status: entity.ParticipantStatusRejected},
		})

		// Assert
		require.NoError(t, err)
		p1, err := repo.FindByEventAndUser(ctx, ev.ID, u1.ID)
		require.NoError(t, err)
		assert.Equal(t, entity.ParticipantStatusConfirmed, p1.Status)
		p2, err := repo.FindByEventAndUser(ctx, ev.ID, u2.ID)
		require.NoError(t, err)
		assert.Equal(t, entity.ParticipantStatusRejected, p2.Status)
	})
}
