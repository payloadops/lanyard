package router

import (
	"context"
	"net/http"
	"plato/app/pkg/client/auth"
	dynamodao "plato/app/pkg/dao/dynamo"
	"strings"
)

// AuthMiddleware checks for the presence of an API key in the request header
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for API Key in the header
		apiKey := r.Header.Get("x-api-key")
		if apiKey != "" {
			apikeyRecord, err := dynamodao.GetApiKey(apiKey)
			if err == nil && apikeyRecord != nil && apikeyRecord.Active && hasScope(apikeyRecord.Scopes, r.Method, r.URL.Path) {
				ctx := context.WithValue(r.Context(), "projectId", apikeyRecord.ProjectId)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}

		// Check for OAuth Token in the Authorization header
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate OAuth token and extract user info
			claims, err := auth.ValidateOAuthToken(token)
			if err == nil && claims != nil {
				// Add user ID to the context
				ctx := context.WithValue(r.Context(), "userID", claims.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}

		// If no valid authentication details are found, return an unauthorized error
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func hasScope(scopes []string, method, path string) bool {
	requiredScope := method + ":" + path
	for _, p := range scopes {
		if p == requiredScope {
			return true
		}
	}
	return false
}
