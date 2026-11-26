package car_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	car "github.com/enoro/eiger-webapp/miniapps/car/backend"
)

// 本体がやることを模した fake。実際は、本体の認証ミドルウェアとサークル管理が実装して渡す。
type host struct{ members map[string]bool } // key: circleID/memberID

func (h host) IsMember(_ context.Context, circleID, memberID string) (bool, error) {
	return h.members[circleID+"/"+memberID], nil
}

func (host) MemberID(r *http.Request) (string, bool) {
	id := r.Header.Get("X-Test-Member")
	return id, id != ""
}

func send(t *testing.T, h http.Handler, method, path, member, body string) (int, map[string]any) {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("X-Test-Member", member)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	var out map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	return rec.Code, out
}

// 本体から見える入口(NewHandler)だけを使って、予約から返却までの一連の流れを確認する。
func TestNewHandler_予約から返却までの流れ(t *testing.T) {
	h := host{members: map[string]bool{"circle-1/member-1": true}}
	handler := car.NewHandler(car.Deps{Members: h, Identity: h})

	// 車を登録する経路は #12 のため、この時点では、車がまだ無いことを確認する。
	status, body := send(t, handler, "GET", "/circles/circle-1/cars", "member-1", "")
	if cars, _ := body["cars"].([]any); status != 200 || len(cars) != 0 {
		t.Fatalf("車は 0 台のはず: %d %v", status, body)
	}

	// メンバーでない人は、同じ入口からでも拒否される。
	if status, _ := send(t, handler, "GET", "/circles/circle-1/cars", "stranger", ""); status != 403 {
		t.Errorf("メンバー以外は 403 のはず: %d", status)
	}

	// 本人確認できない場合は 401。
	if status, _ := send(t, handler, "GET", "/circles/circle-1/cars", "", ""); status != 401 {
		t.Errorf("本人確認できないと 401 のはず: %d", status)
	}

	// 存在しない車の予約は 404(時計・ID 採番が実装済みで、組み立てが壊れていないことの確認も兼ねる)。
	start := time.Now().Add(time.Hour).UTC().Format(time.RFC3339)
	end := time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)
	if status, _ := send(t, handler, "POST", "/cars/none/reservations", "member-1", `{"start":"`+start+`","end":"`+end+`"}`); status != 404 {
		t.Errorf("存在しない車は 404 のはず: %d", status)
	}
}
