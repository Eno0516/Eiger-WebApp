// Package domain は、車の予約ミニアプリのドメイン共通のエラーを定義する。
// このパッケージと配下のパッケージは、標準ライブラリ以外に依存しない。
package domain

import "errors"

var (
	ErrNotFound            = errors.New("対象が見つかりません")
	ErrInvalidPeriod       = errors.New("期間が不正です(終了は開始より後にしてください)")
	ErrPastPeriod          = errors.New("過去の時刻からは予約できません")
	ErrCarUnavailable      = errors.New("この車は現在利用できません")
	ErrReservationConflict = errors.New("その時間帯はすでに予約されています")
	ErrNotMember           = errors.New("このサークルのメンバーではありません")
	ErrNotOwner            = errors.New("予約した本人のみ操作できます")
	ErrInvalidStatus       = errors.New("現在の状態ではこの操作はできません")
	ErrInvalidDistance     = errors.New("走行距離が不正です")
	ErrInvalidLocationLink = errors.New("位置のリンクが不正です(Google Maps のリンクを貼ってください)")
	ErrNoteTooLong         = errors.New("メモが長すぎます")
)
