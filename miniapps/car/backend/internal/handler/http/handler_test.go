package httphandler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	httphandler "github.com/enoro/eiger-webapp/miniapps/car/backend/internal/handler/http"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/testutil"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/car"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/location"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/membership"
	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/usecase/reservation"
)

// headerIdentity は、本体の JWT ミドルウェアの代わりに、ヘッダーから「誰か」を取り出す fake。
type headerIdentity struct{}

func (headerIdentity) MemberID(r *http.Request) (string, bool) {
	id := r.Header.Get("X-Test-Member")
	return id, id != ""
}

type failingDirectory struct{}

func (failingDirectory) IsMember(context.Context, string, string) (bool, error) {
	return false, errors.New("db down: secret-detail")
}

func newHandler(e *testutil.Env, dir membership.Directory) http.Handler {
	return httphandler.New(httphandler.UseCases{
		ListCars:          car.NewList(e.Cars, dir),
		CreateReservation: reservation.NewCreate(e.Cars, e.Reservations, dir, e.Clock, e.IDs),
		CancelReservation: reservation.NewCancel(e.Reservations, dir),
		StartReservation:  reservation.NewStart(e.Reservations, dir),
		ReturnReservation: reservation.NewReturn(e.Reservations, e.Locations, dir, e.Clock),
		ListByRange:       reservation.NewListByRange(e.Reservations, dir),
		ListMine:          reservation.NewListMine(e.Reservations, dir),
		ReportLocation:    location.NewReport(e.Reservations, e.Locations, dir, e.Clock),
		GetLocation:       location.NewGet(e.Cars, e.Locations, dir),
	}, headerIdentity{})
}

// call はリクエストを送り、ステータスと JSON(オブジェクト)を返す。member が空なら未認証。
func call(t *testing.T, h http.Handler, method, path, member, body string) (int, map[string]any) {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if member != "" {
		req.Header.Set("X-Test-Member", member)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	var out map[string]any
	if rec.Body.Len() > 0 {
		if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
			t.Fatalf("レスポンスが JSON ではない: %q", rec.Body.String())
		}
	}
	return rec.Code, out
}

func errorCode(body map[string]any) string {
	e, _ := body["error"].(map[string]any)
	code, _ := e["code"].(string)
	return code
}

func rfc3339(t time.Time) string { return t.Format(time.RFC3339) }

func createBody(e *testutil.Env, startH, endH int) string {
	return `{"start":"` + rfc3339(e.After(startH)) + `","end":"` + rfc3339(e.After(endH)) + `"}`
}

func TestAuth(t *testing.T) {
	e := testutil.NewEnv(t)
	status, body := call(t, newHandler(e, e.Members), "GET", "/circles/circle-1/cars", "", "")
	if status != 401 || errorCode(body) != "UNAUTHORIZED" {
		t.Errorf("本人確認できないリクエストは 401 UNAUTHORIZED のはず: %d %v", status, body)
	}
}

func TestListCars(t *testing.T) {
	t.Run("メンバーは車の一覧を取得できる", func(t *testing.T) {
		e := testutil.NewEnv(t)
		status, body := call(t, newHandler(e, e.Members), "GET", "/circles/circle-1/cars", testutil.Member1, "")
		if status != 200 {
			t.Fatalf("200 のはず: %d %v", status, body)
		}
		cars, _ := body["cars"].([]any)
		first, _ := cars[0].(map[string]any)
		if len(cars) != 1 || first["id"] != "car-1" || first["name"] != "サークルカー" || first["available"] != true {
			t.Errorf("車の内容が想定と違う: %v", body)
		}
	})

	t.Run("サークルのメンバーでない人は 403", func(t *testing.T) {
		e := testutil.NewEnv(t)
		status, body := call(t, newHandler(e, e.Members), "GET", "/circles/circle-1/cars", testutil.Outsider, "")
		if status != 403 || errorCode(body) != "NOT_MEMBER" {
			t.Errorf("403 NOT_MEMBER のはず: %d %v", status, body)
		}
	})
}

func TestCreateReservation(t *testing.T) {
	t.Run("メンバーは予約でき、201 と予約の内容が返る", func(t *testing.T) {
		e := testutil.NewEnv(t)
		status, body := call(t, newHandler(e, e.Members), "POST", "/cars/car-1/reservations", testutil.Member1, createBody(e, 1, 3))
		if status != 201 {
			t.Fatalf("201 のはず: %d %v", status, body)
		}
		if body["id"] != "id-1" || body["status"] != "reserved" || body["carId"] != "car-1" || body["memberId"] != testutil.Member1 ||
			body["start"] != rfc3339(e.After(1)) || body["end"] != rfc3339(e.After(3)) {
			t.Errorf("予約の内容が想定と違う: %v", body)
		}
	})

	t.Run("時間帯が重なると 409", func(t *testing.T) {
		e := testutil.NewEnv(t)
		e.SeedReservation(t, testutil.Member2, 2, 4)
		status, body := call(t, newHandler(e, e.Members), "POST", "/cars/car-1/reservations", testutil.Member1, createBody(e, 3, 5))
		if status != 409 || errorCode(body) != "RESERVATION_CONFLICT" {
			t.Errorf("409 RESERVATION_CONFLICT のはず: %d %v", status, body)
		}
	})

	t.Run("過去の時刻からは 400", func(t *testing.T) {
		e := testutil.NewEnv(t)
		status, body := call(t, newHandler(e, e.Members), "POST", "/cars/car-1/reservations", testutil.Member1, createBody(e, -2, 1))
		if status != 400 || errorCode(body) != "PAST_PERIOD" {
			t.Errorf("400 PAST_PERIOD のはず: %d %v", status, body)
		}
	})

	t.Run("JSON が不正なら 400", func(t *testing.T) {
		e := testutil.NewEnv(t)
		status, body := call(t, newHandler(e, e.Members), "POST", "/cars/car-1/reservations", testutil.Member1, `{"start":`)
		if status != 400 || errorCode(body) != "INVALID_REQUEST" {
			t.Errorf("400 INVALID_REQUEST のはず: %d %v", status, body)
		}
	})

	t.Run("日時の形式が不正なら 400", func(t *testing.T) {
		e := testutil.NewEnv(t)
		status, body := call(t, newHandler(e, e.Members), "POST", "/cars/car-1/reservations", testutil.Member1, `{"start":"明日","end":"明後日"}`)
		if status != 400 || errorCode(body) != "INVALID_REQUEST" {
			t.Errorf("400 INVALID_REQUEST のはず: %d %v", status, body)
		}
	})

	t.Run("存在しない車は 404", func(t *testing.T) {
		e := testutil.NewEnv(t)
		status, body := call(t, newHandler(e, e.Members), "POST", "/cars/none/reservations", testutil.Member1, createBody(e, 1, 3))
		if status != 404 || errorCode(body) != "NOT_FOUND" {
			t.Errorf("404 NOT_FOUND のはず: %d %v", status, body)
		}
	})
}

func TestCalendar(t *testing.T) {
	t.Run("期間と重なる予約を取得できる", func(t *testing.T) {
		e := testutil.NewEnv(t)
		e.SeedReservation(t, testutil.Member2, 1, 3)
		e.SeedReservation(t, testutil.Member1, 40, 42) // 期間外
		path := "/circles/circle-1/reservations?from=" + rfc3339(e.After(0)) + "&to=" + rfc3339(e.After(24))
		status, body := call(t, newHandler(e, e.Members), "GET", strings.ReplaceAll(path, "+", "%2B"), testutil.Member1, "")
		if status != 200 {
			t.Fatalf("200 のはず: %d %v", status, body)
		}
		list, _ := body["reservations"].([]any)
		if len(list) != 1 {
			t.Errorf("期間内の 1 件のはず: %v", body)
		}
	})

	t.Run("from / to が無い・不正なら 400", func(t *testing.T) {
		e := testutil.NewEnv(t)
		for _, path := range []string{
			"/circles/circle-1/reservations",
			"/circles/circle-1/reservations?from=あした&to=あさって",
		} {
			status, body := call(t, newHandler(e, e.Members), "GET", path, testutil.Member1, "")
			if status != 400 || errorCode(body) != "INVALID_REQUEST" {
				t.Errorf("%s: 400 INVALID_REQUEST のはず: %d %v", path, status, body)
			}
		}
	})
}

func TestListMine(t *testing.T) {
	e := testutil.NewEnv(t)
	e.SeedReservation(t, testutil.Member1, 1, 3)
	e.SeedReservation(t, testutil.Member2, 5, 7)
	status, body := call(t, newHandler(e, e.Members), "GET", "/circles/circle-1/reservations/mine", testutil.Member1, "")
	list, _ := body["reservations"].([]any)
	if status != 200 || len(list) != 1 {
		t.Errorf("自分の 1 件のはず: %d %v", status, body)
	}
}

func TestCancelAndStart(t *testing.T) {
	t.Run("本人はキャンセルできる", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedReservation(t, testutil.Member1, 1, 3)
		status, body := call(t, newHandler(e, e.Members), "POST", "/reservations/"+r.ID+"/cancel", testutil.Member1, "")
		if status != 200 || body["status"] != "cancelled" {
			t.Errorf("200 で cancelled のはず: %d %v", status, body)
		}
	})

	t.Run("本人以外は 403", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedReservation(t, testutil.Member1, 1, 3)
		status, body := call(t, newHandler(e, e.Members), "POST", "/reservations/"+r.ID+"/cancel", testutil.Member2, "")
		if status != 403 || errorCode(body) != "NOT_OWNER" {
			t.Errorf("403 NOT_OWNER のはず: %d %v", status, body)
		}
	})

	t.Run("存在しない予約は 404", func(t *testing.T) {
		e := testutil.NewEnv(t)
		status, _ := call(t, newHandler(e, e.Members), "POST", "/reservations/none/cancel", testutil.Member1, "")
		if status != 404 {
			t.Errorf("404 のはず: %d", status)
		}
	})

	t.Run("本人は利用を開始でき、開始後のキャンセルは 409", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedReservation(t, testutil.Member1, 1, 3)
		h := newHandler(e, e.Members)
		if status, body := call(t, h, "POST", "/reservations/"+r.ID+"/start", testutil.Member1, ""); status != 200 || body["status"] != "in_use" {
			t.Fatalf("200 で in_use のはず: %d %v", status, body)
		}
		status, body := call(t, h, "POST", "/reservations/"+r.ID+"/cancel", testutil.Member1, "")
		if status != 409 || errorCode(body) != "INVALID_STATUS" {
			t.Errorf("409 INVALID_STATUS のはず: %d %v", status, body)
		}
	})
}

func TestReturn(t *testing.T) {
	t.Run("返却できる。走行距離とメモ、駐車場所が反映される", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedInUse(t, testutil.Member1, 1, 3)
		h := newHandler(e, e.Members)
		status, body := call(t, h, "POST", "/reservations/"+r.ID+"/return", testutil.Member1,
			`{"distanceKm":42,"note":"給油済み","parkingLink":"https://maps.app.goo.gl/AbCdEf123","parkingNote":"駅前"}`)
		if status != 200 || body["status"] != "returned" || body["distanceKm"] != float64(42) || body["returnNote"] != "給油済み" {
			t.Fatalf("返却の内容が想定と違う: %d %v", status, body)
		}
		status, body = call(t, h, "GET", "/cars/car-1/location", testutil.Member2, "")
		if status != 200 || body["link"] != "https://maps.app.goo.gl/AbCdEf123" || body["note"] != "駅前" {
			t.Errorf("駐車場所が位置として見られるはず: %d %v", status, body)
		}
	})

	t.Run("駐車場所のリンクが不正なら 400", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedInUse(t, testutil.Member1, 1, 3)
		status, body := call(t, newHandler(e, e.Members), "POST", "/reservations/"+r.ID+"/return", testutil.Member1,
			`{"distanceKm":1,"parkingLink":"https://evil.example.com/maps"}`)
		if status != 400 || errorCode(body) != "INVALID_LOCATION_LINK" {
			t.Errorf("400 INVALID_LOCATION_LINK のはず: %d %v", status, body)
		}
	})
}

func TestLocation(t *testing.T) {
	t.Run("利用中の本人は位置を申告でき、メンバーは取得できる", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedInUse(t, testutil.Member1, 1, 3)
		h := newHandler(e, e.Members)
		status, body := call(t, h, "POST", "/reservations/"+r.ID+"/location", testutil.Member1,
			`{"link":"https://maps.app.goo.gl/AbCdEf123","note":"駅前"}`)
		if status != 201 || body["carId"] != "car-1" || body["memberId"] != testutil.Member1 {
			t.Fatalf("201 のはず: %d %v", status, body)
		}
		status, body = call(t, h, "GET", "/cars/car-1/location", testutil.Member2, "")
		if status != 200 || body["note"] != "駅前" || body["reportedAt"] != rfc3339(e.Clock.T) {
			t.Errorf("取得できるはず: %d %v", status, body)
		}
	})

	t.Run("利用を開始していない予約での申告は 409", func(t *testing.T) {
		e := testutil.NewEnv(t)
		r := e.SeedReservation(t, testutil.Member1, 1, 3)
		status, body := call(t, newHandler(e, e.Members), "POST", "/reservations/"+r.ID+"/location", testutil.Member1,
			`{"link":"https://maps.app.goo.gl/AbCdEf123"}`)
		if status != 409 || errorCode(body) != "INVALID_STATUS" {
			t.Errorf("409 INVALID_STATUS のはず: %d %v", status, body)
		}
	})

	t.Run("申告がまだない車は 404", func(t *testing.T) {
		e := testutil.NewEnv(t)
		status, _ := call(t, newHandler(e, e.Members), "GET", "/cars/car-1/location", testutil.Member1, "")
		if status != 404 {
			t.Errorf("404 のはず: %d", status)
		}
	})

	t.Run("サークルのメンバーでない人は位置を見られない(403)", func(t *testing.T) {
		e := testutil.NewEnv(t)
		status, _ := call(t, newHandler(e, e.Members), "GET", "/cars/car-1/location", testutil.Outsider, "")
		if status != 403 {
			t.Errorf("403 のはず: %d", status)
		}
	})
}

func TestInternalError(t *testing.T) {
	e := testutil.NewEnv(t)
	status, body := call(t, newHandler(e, failingDirectory{}), "GET", "/circles/circle-1/cars", testutil.Member1, "")
	if status != 500 || errorCode(body) != "INTERNAL" {
		t.Fatalf("500 INTERNAL のはず: %d %v", status, body)
	}
	if raw, _ := json.Marshal(body); strings.Contains(string(raw), "secret-detail") {
		t.Errorf("内部エラーの詳細をレスポンスに出してはいけない: %s", raw)
	}
}
