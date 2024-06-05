package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/golang-jwt/jwt"
	"github.com/payloadops/plato/app/config"
	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/utils"
)

// Claims represents the JWT claims containing the standard claims, user ID, and organization ID.
type Claims struct {
	jwt.StandardClaims
	OrgID string `json:"org"`
}

func APIKeyAuthMiddleware(cfg *config.Config, logger *zap.Logger, apiKeyManager dal.APIKeyManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := middleware.GetReqID(r.Context())
			// Extract the token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
				return
			}

			// Check if the Authorization header is using Basic Auth
			if !strings.HasPrefix(authHeader, "Basic ") {
				http.Error(w, "Invalid Authorization Header", http.StatusUnauthorized)
				return
			}

			// Decode the Base64 encoded credentials
			base64Credentials := strings.TrimPrefix(authHeader, "Basic ")
			decodedCredentials, err := base64.StdEncoding.DecodeString(base64Credentials)
			if err != nil {
				http.Error(w, "Invalid Base64 Encoding", http.StatusUnauthorized)
				return
			}

			// Split the credentials into clientID and clientSecret
			credentials := strings.SplitN(string(decodedCredentials), ":", 2)
			if len(credentials) != 2 {
				http.Error(w, "Invalid Authorization Header Format", http.StatusUnauthorized)
				return
			}

			clientID, clientSecret := credentials[0], credentials[1]
			key, err := apiKeyManager.GetAPIKeyByID(r.Context(), clientID)
			if err != nil {
				logger.Error("failed to get API key",
					zap.String("requestID", requestID),
					zap.Error(err),
				)

				http.Error(w, "Unexpected Error", http.StatusInternalServerError)
				return
			}

			if key == nil {
				http.Error(w, "Invalid API Key", http.StatusUnauthorized)
				return
			}

			if key.Deleted {
				logger.Error("attempted use of deleted API key", zap.String("requestID", requestID))
				http.Error(w, "Invalid API Key", http.StatusUnauthorized)
				return
			}

			if !utils.SecureCompare(clientSecret, key.Secret) {
				logger.Error("invalid API key secret",
					zap.String("requestID", requestID),
					zap.Error(err),
				)

				http.Error(w, "Invalid API Key", http.StatusUnauthorized)
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
