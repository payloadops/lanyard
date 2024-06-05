package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/golang-jwt/jwt"
	"github.com/payloadops/plato/app/config"
	"github.com/payloadops/plato/app/dal"
)

// Claims represents the JWT claims containing the standard claims, user ID, and organization ID.
type Claims struct {
	jwt.StandardClaims
	OrgID string `json:"org"`
}

func APIKeyAuthMiddleware(cfg *config.Config, logger *zap.Logger, apiKeyDBClient *dal.APIKeyDBClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := middleware.GetReqID(r.Context())
			// Extract the token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
				return
			}

			clientID, clientSecret := strings.Split(authHeader, ":")[0], strings.Split(authHeader, ":")[1]

			if clientID == "" {
				http.Error(w, "Client ID Not Found in Authorization Header", http.StatusUnauthorized)
				return
			}
			if clientSecret == "" {
				http.Error(w, "Client Secret Not Found in Authorization Header", http.StatusUnauthorized)
				return
			}

			key, err := apiKeyDBClient.GetAPIKeyByID(r.Context(), clientID)

			if err != nil && key.Deleted {
				logger.Error("deleted key",
					zap.String("requestID", requestID),
					zap.Error(err),
				)

				http.Error(w, "Cannot Use Deleted API Key", http.StatusUnauthorized)
				return
			}

			if clientSecret != key.Secret {
				logger.Error("invalid secret",
					zap.String("requestID", requestID),
					zap.Error(err),
				)

				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Set the user and org context
			ctx := context.WithValue(r.Context(), "orgID", key.OrgID)
			ctx = context.WithValue(ctx, "projectID", key.ProjectID)

			// Call the next handler with the new context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// JWTAuthMiddleware returns a middleware function that validates the JWT token from the Authorization header.
// It sets the user ID and organization ID in the request context if the token is valid.
func JWTAuthMiddleware(cfg *config.Config, logger *zap.Logger) func(http.Handler) http.Handler {
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

			if claims.OrgID == "" {
				err := fmt.Errorf("required field '%s' is empty value", "org")
				logger.Error("failed to parse org from claims",
					zap.String("requestID", requestID),
					zap.Error(err),
				)

				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if claims.Subject == "" {
				err := fmt.Errorf("required field '%s' is empty value", "sub")
				logger.Error("failed to parse sub from claims",
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
