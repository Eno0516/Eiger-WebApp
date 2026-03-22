package http

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/enoro/eiger-webapp/backend/internal/handler/http/dto"
	"github.com/enoro/eiger-webapp/backend/internal/usecase/auth"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type AuthHandler struct {
	registerUC           *auth.RegisterUseCase
	loginUC              *auth.LoginUseCase
	verifyEmailUC        *auth.VerifyEmailUseCase
	resendVerificationUC *auth.ResendVerificationUseCase
	forgotPasswordUC     *auth.ForgotPasswordUseCase
	resetPasswordUC      *auth.ResetPasswordUseCase
	log                  zerolog.Logger
}

func NewAuthHandler(
	registerUC *auth.RegisterUseCase,
	loginUC *auth.LoginUseCase,
	verifyEmailUC *auth.VerifyEmailUseCase,
	resendVerificationUC *auth.ResendVerificationUseCase,
	forgotPasswordUC *auth.ForgotPasswordUseCase,
	resetPasswordUC *auth.ResetPasswordUseCase,
	log zerolog.Logger,
) *AuthHandler {
	return &AuthHandler{
		registerUC:           registerUC,
		loginUC:              loginUC,
		verifyEmailUC:        verifyEmailUC,
		resendVerificationUC: resendVerificationUC,
		forgotPasswordUC:     forgotPasswordUC,
		resetPasswordUC:      resetPasswordUC,
		log:                  log,
	}
}

// Register POST /auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// TODO: registerUC.Execute を呼び出す
	// TODO: 201 を返す（ボディはメッセージのみ）
	panic("not implemented")
}

// Login POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// TODO: loginUC.Execute を呼び出す
	// TODO: 200 + dto.LoginResponse{Token: out.Token} を返す
	panic("not implemented")
}

// VerifyEmail GET /auth/verify-email?token=xxx
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		h.writeError(w, http.StatusBadRequest, "token is required")
		return
	}
	// TODO: verifyEmailUC.Execute を呼び出す
	// TODO: 200 + メッセージを返す
	panic("not implemented")
}

// ResendVerification POST /auth/resend-verification
func (h *AuthHandler) ResendVerification(w http.ResponseWriter, r *http.Request) {
	var req dto.ResendVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// TODO: resendVerificationUC.Execute を呼び出す
	// TODO: 200 + メッセージを返す
	panic("not implemented")
}

// ForgotPassword POST /auth/forgot-password
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// TODO: forgotPasswordUC.Execute を呼び出す
	// TODO: 200 + メッセージを返す（ユーザー存在有無に関わらず同じレスポンス）
	panic("not implemented")
}

// ResetPassword POST /auth/reset-password
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// TODO: resetPasswordUC.Execute を呼び出す
	// TODO: 200 + メッセージを返す
	panic("not implemented")
}

func (h *AuthHandler) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.log.Error().Err(err).Msg("encode response")
	}
}

func (h *AuthHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, dto.ErrorResponse{Error: message})
}

// import guard
var _ = apperrors.HTTPStatusCode
