package auth

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
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
func AuthMiddleware(cfg *config.Config, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := middleware.GetReqID(r.Context())
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
				// Ensure the token method conforms to "alg" expected value
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					logger.Warn("unexpected signing method",
						zap.String("requestID", requestID),
					)

					return nil, errors.New("unexpected signing method")
				}

				return []byte(cfg.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				logger.Error("invalid token",
					zap.String("requestID", requestID),
					zap.Error(err),
				)

				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Set the user and org context
			ctx := context.WithValue(r.Context(), "orgID", claims.OrgID)
			ctx = context.WithValue(ctx, "userID", claims.Subject)

			// Call the next handler with the new context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
