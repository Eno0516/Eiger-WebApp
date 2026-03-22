package apperrors

import (
	"errors"
	"net/http"
)

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

func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

var (
	ErrNotFound     = New(http.StatusNotFound, "not found")
	ErrBadRequest   = New(http.StatusBadRequest, "bad request")
	ErrConflict     = New(http.StatusConflict, "conflict")
	ErrUnauthorized = New(http.StatusUnauthorized, "unauthorized")
	ErrInternal     = New(http.StatusInternalServerError, "internal server error")
)

func HTTPStatusCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return http.StatusInternalServerError
}
