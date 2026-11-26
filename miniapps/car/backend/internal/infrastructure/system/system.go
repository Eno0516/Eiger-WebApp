// Package system は、時計と ID 採番の本番用の実装。
package system

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// Clock は現在時刻を返す。
type Clock struct{}

func (Clock) Now() time.Time { return time.Now() }

// RandomID は、128 ビットの乱数を 16 進数 32 文字にした ID を採番する。外部ライブラリは使わない。
type RandomID struct{}

func (RandomID) NewID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand が失敗するのは、OS の乱数源が使えないとき。続行できないため panic する。
		panic("乱数の取得に失敗しました: " + err.Error())
	}
	return hex.EncodeToString(b)
}
