package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"

	"github.com/enoro/eiger-webapp/backend/internal/handler/http/middleware"
)

func NewRouter(userHandler *UserHandler, log zerolog.Logger) http.Handler {
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
		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.GetAll)
			r.Post("/", userHandler.Create)
			r.Get("/{id}", userHandler.GetByID)
		})
	})

	return r
}
