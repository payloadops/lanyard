package router

import (
	"context"
	"fmt"
	"net/http"

	"plato/app/pkg/auth"
	"plato/app/pkg/service/apikey"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userCtx context.Context
		var err error

		apiKey := r.Header.Get("x-api-key")
		if apiKey != "" {
			userCtx, err = validateAPIKey(r, apiKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
		}

		// if userCtx == nil {
		// 	authHeader := r.Header.Get("Authorization")
		// 	if strings.HasPrefix(authHeader, "Bearer ") {
		// 		userCtx, err = validateOAuthToken(r.Context(), strings.TrimPrefix(authHeader, "Bearer "))
		// 		if err != nil {
		// 			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		// 			return
		// 		}
		// 	}
		// }

		if userCtx == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(userCtx))
	})
}

func validateAPIKey(r *http.Request, apiKey string) (context.Context, error) {
	apikeyService := apikey.NewService()
	apikeyRecord, err := apikeyService.GetApiKey(r.Context(), apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get api key: %w", err)
	}
	if apikeyRecord == nil || !apikeyRecord.Active {
		return nil, fmt.Errorf("invalid or inactive api key")
	}
	if !validateScopes(apikeyRecord.Scopes, r.Method, r.URL.Path) {
		return nil, fmt.Errorf("insufficient permissions for %s %s", r.Method, r.URL.Path)
	}

	userCtx := context.WithValue(r.Context(), auth.OrgContext{}, apikeyRecord.OrgId)
	return context.WithValue(userCtx, auth.ProjectContext{}, apikeyRecord.ProjectId), nil
}

// func validateOAuthToken(ctx context.Context, token string) (context.Context, error) {

// 	if err != nil {
// 		return nil, err
// 	}
// 	userCtx := context.WithValue(ctx, auth.OrgId{}, claims.OrgId)
// 	return context.WithValue(userCtx, auth.UserId{}, claims.UserId), nil
// }

func validateScopes(scopes []string, method string, path string) bool {
	// Implement your scope validation logic here
	// This example always returns true, replace with your authorization checks
	return true
}
