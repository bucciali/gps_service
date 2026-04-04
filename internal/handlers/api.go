package handlers

import (
	"encoding/json"
	"map-backend/internal/db"
	"map-backend/internal/response"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type CreatePointRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	TypeID      string  `json:"type_id"`
}

type UpdatePointRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	TypeID      string  `json:"type_id"`
}

func GetPointsHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		points, err := db.GetPoints(r.Context(), pool)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, "failed to get points")
			return
		}
		response.WriteJSON(w, http.StatusOK, map[string]any{
			"points": points,
		})
	}
}

func CreatePointHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreatePointRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		defer r.Body.Close()

		if strings.TrimSpace(req.Name) == "" {
			response.WriteError(w, http.StatusBadRequest, "name is required")
			return
		}

		userID := r.Context().Value("user_id").(string)

		p := db.Point{
			Name:        req.Name,
			Description: req.Description,
			Latitude:    req.Latitude,
			Longitude:   req.Longitude,
			TypeID:      req.TypeID,
			CreatedBy:   userID,
		}

		created, err := db.CreatePoint(r.Context(), pool, p)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, "failed to create point")
			return
		}
		response.WriteJSON(w, http.StatusCreated, map[string]any{
			"point": created,
		})
	}
}

func UpdatePointHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var req UpdatePointRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		defer r.Body.Close()

		p := db.Point{
			PointId:     id,
			Name:        req.Name,
			Description: req.Description,
			Latitude:    req.Latitude,
			Longitude:   req.Longitude,
			TypeID:      req.TypeID,
		}

		if err := db.UpdatePoint(r.Context(), pool, p); err != nil {
			if err.Error() == "point is not found" {
				response.WriteError(w, http.StatusNotFound, "point not found")
				return
			}
			response.WriteError(w, http.StatusInternalServerError, "failed to update point")
			return
		}
		response.WriteJSON(w, http.StatusOK, map[string]any{
			"message": "updated",
		})
	}
}

func DeletePointHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if err := db.DeletePoint(r.Context(), pool, id); err != nil {
			if err.Error() == "point is not found" {
				response.WriteError(w, http.StatusNotFound, "point not found")
				return
			}
			response.WriteError(w, http.StatusInternalServerError, "failed to delete point")
			return
		}
		response.WriteJSON(w, http.StatusOK, map[string]any{
			"message": "deleted",
		})
	}
}

func PingDbHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := pool.Ping(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("problems with pinging")
			response.WriteError(w, http.StatusInternalServerError, "database ping failed")
			return
		}
		response.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}

}
