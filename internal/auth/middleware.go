package auth

import (
	"gps_service/internal/response"
	"net/http"
	"strings"
)

func AuthMiddleware(JWTManager *JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.WriteError(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}
			if !strings.HasPrefix(authHeader, "Bearer ") {
				response.WriteError(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				response.WriteError(w, http.StatusUnauthorized, "missing token")
				return
			}
			claims, err := JWTManager.ParseToken(tokenString)
			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}
			ctx := WithUserRole(r.Context(), claims.UserID, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currentRole, ok := GetRoleFromContext(r.Context())
			if !ok {
				response.WriteError(w, http.StatusUnauthorized, "cant auth role")
				return
			}
			if currentRole != role {
				response.WriteError(w, http.StatusForbidden, "dont have rights")
				return

			}
			next.ServeHTTP(w, r)

		})
	}
}
