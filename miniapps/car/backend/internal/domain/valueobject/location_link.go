package valueobject

import (
	"net/url"
	"strings"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
)

const maxLocationLinkLength = 2000

// LocationLink は、車の位置として申告する Google Maps のリンク。
// 任意の URL を貼られてフィッシングなどに使われないよう、許可したドメインとパスだけを受け付ける。
type LocationLink struct {
	raw string
}

// NewLocationLink はリンクを検証して作る。許可外なら domain.ErrInvalidLocationLink を返す。
func NewLocationLink(raw string) (LocationLink, error) {
	if raw == "" || len(raw) > maxLocationLinkLength {
		return LocationLink{}, domain.ErrInvalidLocationLink
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Port() != "" {
		return LocationLink{}, domain.ErrInvalidLocationLink
	}
	if !isAllowedMapsLocation(strings.ToLower(u.Hostname()), u.Path) {
		return LocationLink{}, domain.ErrInvalidLocationLink
	}
	return LocationLink{raw: raw}, nil
}

func (l LocationLink) String() string { return l.raw }

func isAllowedMapsLocation(host, path string) bool {
	switch host {
	case "google.com", "www.google.com":
		return path == "/maps" || strings.HasPrefix(path, "/maps/")
	case "maps.google.com":
		return true
	case "maps.app.goo.gl":
		return len(path) > 1
	case "goo.gl":
		return strings.HasPrefix(path, "/maps/")
	}
	return false
}
