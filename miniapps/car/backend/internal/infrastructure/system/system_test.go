package system_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/infrastructure/system"
)

func TestRandomID(t *testing.T) {
	t.Run("32 文字の 16 進数で、呼ぶたびに別の値になる", func(t *testing.T) {
		gen := system.RandomID{}
		seen := map[string]bool{}
		for i := 0; i < 1000; i++ {
			id := gen.NewID()
			if !regexp.MustCompile(`^[0-9a-f]{32}$`).MatchString(id) {
				t.Fatalf("32 文字の 16 進数のはず: %q", id)
			}
			if seen[id] {
				t.Fatalf("重複した: %q", id)
			}
			seen[id] = true
		}
	})
}

func TestClock(t *testing.T) {
	t.Run("現在時刻を返す", func(t *testing.T) {
		before := time.Now()
		got := system.Clock{}.Now()
		if got.Before(before) || time.Since(got) > time.Second {
			t.Errorf("現在時刻のはず: %v", got)
		}
	})
}
