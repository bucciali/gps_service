package router

import (
	"map-backend/internal/auth"
	"map-backend/internal/handlers"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(jm *auth.JWTManager, pool *pgxpool.Pool) chi.Router {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware(jm))
		r.Get("/api/points", handlers.GetPointsHandler(pool))
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware(jm))
		r.Use(auth.RequireRole("admin"))
		r.Post("/api/points", handlers.CreatePointHandler(pool))
		r.Put("/api/points/{id}", handlers.UpdatePointHandler(pool))
		r.Delete("/api/points/{id}", handlers.DeletePointHandler(pool))
	})

	return r
}
