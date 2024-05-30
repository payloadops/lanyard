package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/payloadops/plato/app/config"
)

// Claims represents the JWT claims containing the standard claims, user ID, and organization ID.
type Claims struct {
	jwt.StandardClaims
	OrgID string `json:"org"`
}

// AuthMiddleware returns a middleware function that validates the JWT token from the Authorization header.
// It sets the user ID and organization ID in the request context if the token is valid.
func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "Invalid Authorization Header", http.StatusUnauthorized)
				return
			}

			// Parse and validate the token
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Set the user and org context
			ctx := context.WithValue(r.Context(), "OrgID", claims.OrgID)
			ctx = context.WithValue(ctx, "UserID", claims.Subject)

			// Call the next handler with the new context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
