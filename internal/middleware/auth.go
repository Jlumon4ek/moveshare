package middleware

import (
	"context"
	"moveshare/internal/services"
	"net/http"
	"strings"
)

type contextKey string

const (
	ContextUserIDKey contextKey = "userID"
)

func AuthMiddleware(jwtService services.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := jwtService.ValidateToken(token)
			if err != nil {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ContextUserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
