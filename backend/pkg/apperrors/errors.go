package apperrors

import (
	"errors"
	"net/http"
)

// AppError はアプリケーション全体で使う構造化エラー。
// HTTP ステータスコードとメッセージを保持し、エラーチェーンをサポートする。
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// New は HTTP ステータスコードとメッセージを持つ AppError を生成する。
// センチネルエラーの定義に使う。
func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Wrap は既存のエラーをラップして AppError を生成する。
// 直接使用は非推奨。下記の型付きコンストラクタを使うこと。
func Wrap(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// --- センチネルエラー ---
// repository 層が DB エラーをこれらに変換する。
// errors.Is() / errors.As() で比較する場合は Wrap 系を使うこと。

var (
	ErrNotFound     = New(http.StatusNotFound, "not found")
	ErrBadRequest   = New(http.StatusBadRequest, "bad request")
	ErrConflict     = New(http.StatusConflict, "conflict")
	ErrUnauthorized = New(http.StatusUnauthorized, "unauthorized")
	ErrForbidden    = New(http.StatusForbidden, "forbidden")
	ErrInternal     = New(http.StatusInternalServerError, "internal server error")
)

// --- 型付きコンストラクタ ---
// usecase 層でエラーを生成するときはこれらを使い、raw な HTTP ステータスコードを書かない。

func NewBadRequest(message string, err error) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message, Err: err}
}

func NewNotFound(message string, err error) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message, Err: err}
}

func NewConflict(message string, err error) *AppError {
	return &AppError{Code: http.StatusConflict, Message: message, Err: err}
}

func NewUnauthorized(message string, err error) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: message, Err: err}
}

func NewForbidden(message string, err error) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: message, Err: err}
}

func NewInternal(message string, err error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: message, Err: err}
}

// HTTPStatusCode は err から HTTP ステータスコードを取得する。
// AppError でなければ 500 を返す。
// handler 層でのみ呼び出す。
func HTTPStatusCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return http.StatusInternalServerError
}

// HTTPMessage は err から クライアントへ返すメッセージを取得する。
// AppError でなければ "internal server error" を返す（内部エラーを隠蔽）。
// handler 層でのみ呼び出す。
func HTTPMessage(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Message
	}
	return "internal server error"
}
