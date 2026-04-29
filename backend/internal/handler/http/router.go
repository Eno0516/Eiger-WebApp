package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"

	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/internal/handler/http/middleware"
	jwtpkg "github.com/enoro/eiger-webapp/backend/pkg/jwt"
)

func NewRouter(
	userHandler *UserHandler,
	authHandler *AuthHandler,
	adminHandler *AdminHandler,
	circleHandler *CircleHandler,
	eventHandler *EventHandler,
	circleMemberRepo domainrepo.CircleMemberRepository,
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
			r.Post("/line", authHandler.LineAuth)
			r.Post("/logout", authHandler.Logout)
		})

		// ログイン必須
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(jwtGen))

			// 初期登録
			r.Post("/auth/signup", authHandler.Register)

			// 自分のユーザー情報
			r.Get("/users/me", userHandler.GetMe)

			// サークル一覧・参加申請（一般ユーザー向け）
			r.Get("/circles", circleHandler.GetAll)
			r.Post("/circles/{id}/apply", circleHandler.Apply)

			// ユーザー情報
			r.Get("/users/{id}", userHandler.GetByID)
			r.Patch("/users/{id}/display-name", userHandler.UpdateDisplayName)
		})

		// サークル admin のみ（circle_admin が自サークルを操作する）
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(jwtGen))
			r.Use(middleware.RequireCircleAdmin(circleMemberRepo))

			r.Get("/circles/{id}/members", circleHandler.GetMembers)
			r.Get("/circles/{id}/members/{userId}/form-response", circleHandler.GetFormResponse)
			r.Patch("/circles/{id}/members/{userId}/role", circleHandler.UpdateMemberRole)
		})

		// ログイン必須（user 向けイベント参加申請）
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(jwtGen))
			r.Post("/circles/{circleId}/events/{eventId}/apply", eventHandler.ApplyEvent)
		})

		// 管理画面（system_admin または circle_admin）
		r.Route("/admin", func(r chi.Router) {
			r.Use(middleware.JWTAuth(jwtGen))

			// サークル一覧: system_admin=全件、circle_admin=自分のサークルのみ
			// RequireSystemAdminOrCircleAdmin をここでは使わず、usecase 側でフィルタリング
			// （circle_admin かどうかは listCirclesUC 内で判定）
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireSystemAdmin())
				r.Post("/circles", adminHandler.CreateCircle)
				r.Delete("/circles/{id}", adminHandler.DeleteCircle)
			})

			// system_admin または circle_admin（サークル所属チェックあり）
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireSystemAdminOrCircleAdmin(circleMemberRepo))
				r.Get("/circles", adminHandler.GetCircles)
				r.Get("/circles/{id}", adminHandler.GetCircle)
				r.Get("/circles/{id}/members", adminHandler.GetCircleMembers)
				r.Patch("/circles/{id}/members/{userId}/role", adminHandler.UpdateCircleMemberRole)

				// イベント管理
				r.Get("/circles/{id}/events", eventHandler.ListEvents)
				r.Post("/circles/{id}/events", eventHandler.CreateEvent)
				r.Get("/circles/{id}/events/{eventId}/participants", eventHandler.GetParticipants)
				r.Patch("/circles/{id}/events/{eventId}/participants", eventHandler.BulkUpdateParticipants)
				r.Patch("/circles/{id}/events/{eventId}/participants/{userId}", eventHandler.UpdateParticipantStatus)
				r.Post("/circles/{id}/events/{eventId}/confirm", eventHandler.ConfirmEvent)
			})
		})
	})

	return r
}
