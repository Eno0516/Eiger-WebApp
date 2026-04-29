package http

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/enoro/eiger-webapp/backend/internal/handler/http/dto"
	"github.com/enoro/eiger-webapp/backend/internal/handler/http/middleware"
	"github.com/enoro/eiger-webapp/backend/internal/usecase/auth"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type AuthHandler struct {
	lineAuthUC *auth.LineAuthUseCase
	registerUC *auth.RegisterUseCase
	log        zerolog.Logger
}

func NewAuthHandler(lineAuthUC *auth.LineAuthUseCase, registerUC *auth.RegisterUseCase, log zerolog.Logger) *AuthHandler {
	return &AuthHandler{lineAuthUC: lineAuthUC, registerUC: registerUC, log: log}
}

// LineAuth POST /auth/line
func (h *AuthHandler) LineAuth(w http.ResponseWriter, r *http.Request) {
	var req dto.LineAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	out, err := h.lineAuthUC.Execute(r.Context(), auth.LineAuthInput{IDToken: req.IDToken})
	if err != nil {
		if apperrors.HTTPStatusCode(err) >= 500 {
			h.log.Error().Err(err).Msg("line auth error")
		}
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	h.writeJSON(w, http.StatusOK, dto.LoginResponse{Token: out.Token, IsNewUser: out.IsNewUser})
}

// Logout POST /auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// JWT はステートレスなためサーバー側の処理は不要。
	// フロント側で localStorage のトークンを削除する。
	h.writeJSON(w, http.StatusOK, dto.MessageResponse{Message: "logged out"})
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

// Register POST /auth/signup
// JWT から userID を取得し、初期登録（UserProfile 作成 + Sheets 書き込み）を行う。
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	out, err := h.registerUC.Execute(r.Context(), auth.RegisterInput{
		UserID:      claims.UserID,
		DisplayName: req.DisplayName,
		University:  req.University,
		Faculty:     req.Faculty,
		Grade:       req.Grade,
	})
	if err != nil {
		if apperrors.HTTPStatusCode(err) >= 500 {
			h.log.Error().Err(err).Msg("register error")
		}
		h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
		return
	}
	h.writeJSON(w, http.StatusCreated, dto.RegisterResponse{DisplayName: out.Profile.DisplayName})
}
