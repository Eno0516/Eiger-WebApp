package valueobject_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/valueobject"
)

func TestNewLocationLink_許可するリンク(t *testing.T) {
	allowed := []struct{ name, url string }{
		{"Google Maps の通常のURL", "https://www.google.com/maps/place/Tokyo+Station/@35.68,139.76,17z"},
		{"www なしの google.com/maps", "https://google.com/maps?q=35.68,139.76"},
		{"maps.google.com", "https://maps.google.com/?q=35.68,139.76"},
		{"スマホの共有で出る短縮URL", "https://maps.app.goo.gl/AbCdEf123"},
		{"goo.gl/maps の短縮URL", "https://goo.gl/maps/AbCdEf123"},
	}
	for _, c := range allowed {
		t.Run(c.name, func(t *testing.T) {
			link, err := valueobject.NewLocationLink(c.url)
			if err != nil {
				t.Fatalf("許可されるはず: %v", err)
			}
			if link.String() != c.url {
				t.Errorf("URL がそのまま保持されるはず: %q", link.String())
			}
		})
	}
}

func TestNewLocationLink_拒否するリンク(t *testing.T) {
	rejected := []struct{ name, url string }{
		{"空文字", ""},
		{"http は拒否する", "http://www.google.com/maps/place/x"},
		{"許可外のドメインは拒否する", "https://evil.example.com/maps"},
		{"許可ドメインを先頭に含むだけの別ドメインは拒否する", "https://google.com.evil.example.com/maps"},
		{"ユーザー情報で許可ドメインを装うURLは拒否する", "https://google.com@evil.example.com/maps"},
		{"google.com でも maps 以外のパスは拒否する", "https://www.google.com/search?q=maps"},
		{"goo.gl でも maps 以外のパスは拒否する", "https://goo.gl/abc"},
		{"javascript スキームは拒否する", "javascript:alert(1)"},
		{"長すぎるURLは拒否する", "https://maps.app.goo.gl/" + strings.Repeat("a", 2100)},
	}
	for _, c := range rejected {
		t.Run(c.name, func(t *testing.T) {
			_, err := valueobject.NewLocationLink(c.url)
			if !errors.Is(err, domain.ErrInvalidLocationLink) {
				t.Errorf("ErrInvalidLocationLink を期待: %v", err)
			}
		})
	}
}
