package router

import (
	"context"
	"net/http"
	"plato/app/pkg/client/auth"
	"plato/app/pkg/service/apikey"
	"strings"
)

// AuthMiddleware checks for the presence of an API key in the request header
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apikeyService := apikey.NewService()
		apiKey := r.Header.Get("x-api-key")
		if apiKey != "" {
			apikeyRecord, err := apikeyService.GetAPIKey(r.Context(), apiKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err == nil && apikeyRecord != nil && apikeyRecord.Active && validateScopes(apikeyRecord.Scopes, r.Method, r.URL.Path) {
				ctx := context.WithValue(r.Context(), "projectId", apikeyRecord.ApiKey)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
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
				return
			}
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func validateScopes(scopes []string, method string, path string) bool {
	// scopeHierarchy := map[string]int{
	// 	"read":  1,
	// 	"write": 2,
	// 	"admin": 3,
	// }

	// scopeMethods := map[string]string{
	// 	"prompts": "prompts",
	// 	"keys":    "keys",
	// 	"users":   "users",
	// 	"teams":   "",
	// }

	return true
}
