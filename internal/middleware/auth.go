package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"weather-api/internal/auth"
)

type contextKey string

const claimsKey contextKey = "claims"

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				writeErr(w, http.StatusUnauthorized, "missing or invalid token")
				return
			}

			claims, err := auth.ParseToken(strings.TrimPrefix(header, "Bearer "), secret)
			if err != nil {
				writeErr(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := ClaimsFrom(r.Context())
			if claims == nil || claims.Role != role {
				writeErr(w, http.StatusForbidden, "forbidden")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func ClaimsFrom(ctx context.Context) *auth.Claims {
	v, _ := ctx.Value(claimsKey).(*auth.Claims)
	return v
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
