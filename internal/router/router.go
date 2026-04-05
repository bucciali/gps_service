package router

import (
	"gps_service/internal/auth"
	"gps_service/internal/handlers"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(jm *auth.JWTManager, pool *pgxpool.Pool) chi.Router {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Post("/api/auth/login", handlers.DummyLoginHandler(jm, pool))

	r.Get("/api/points", handlers.GetPointsHandler(pool))

	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware(jm))
		r.Use(auth.RequireRole("admin"))
		r.Post("/api/points", handlers.CreatePointHandler(pool))
		r.Put("/api/points/{id}", handlers.UpdatePointHandler(pool))
		r.Delete("/api/points/{id}", handlers.DeletePointHandler(pool))
	})

	return r
}
