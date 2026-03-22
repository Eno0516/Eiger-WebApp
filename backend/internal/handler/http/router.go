package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	"github.com/enoro/eiger-webapp/backend/internal/handler/http/middleware"
	jwtpkg "github.com/enoro/eiger-webapp/backend/pkg/jwt"
)

func NewRouter(
	userHandler *UserHandler,
	authHandler *AuthHandler,
	adminHandler *AdminHandler,
	jwtGen *jwtpkg.Generator,
	log zerolog.Logger,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(middleware.CORS)
	r.Use(middleware.Logger(log))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// 認証不要
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Get("/verify-email", authHandler.VerifyEmail)
			r.Post("/resend-verification", authHandler.ResendVerification)
			r.Post("/forgot-password", authHandler.ForgotPassword)
			r.Post("/reset-password", authHandler.ResetPassword)
		})

		// pending 以上（ログイン必須）
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(jwtGen))

			r.Route("/users", func(r chi.Router) {
				r.Get("/", userHandler.GetAll)
				r.Post("/", userHandler.Create)
				r.Get("/{id}", userHandler.GetByID)
			})
		})

		// admin のみ
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(jwtGen))
			r.Use(middleware.RequireRole(entity.RoleAdmin))

			r.Route("/admin", func(r chi.Router) {
				r.Patch("/users/{id}/role", adminHandler.UpdateUserRole)
			})
		})
	})

	return r
}
