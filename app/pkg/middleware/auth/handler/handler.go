package handler

import (
	"context"
	"fmt"
	"net/http"
	"plato/app/pkg/service/token"
	"regexp"
	"strings"

	"plato/app/pkg/auth"
	"plato/app/pkg/service/apikey"
)

var _ Authenticator = (*Handler)(nil)

// Handler handles authentication processes for various authentication schemes.
type Handler struct {
	apiKeyService apikey.ApiKeyService
	tokenService  token.TokenService
}

// NewHandler initializes a new instance of Handler.
func NewHandler() *Handler {
	return &Handler{
		apiKeyService: apikey.NewService(),
		tokenService:  token.NewService(),
	}
}

// Authenticate determines the user context based on the request's authorization method.
func (a *Handler) Authenticate(r *http.Request) (context.Context, error) {
	apiKey := r.Header.Get("x-api-key")
	if apiKey != "" {
		return a.validateAPIKey(r, apiKey)
	}

	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return a.validateOAuthToken(r, strings.TrimPrefix(authHeader, "Bearer "))
	}

	return nil, fmt.Errorf("unauthorized")
}

// validateAPIKey validates the API key and returns the user context if successful.
func (a *Handler) validateAPIKey(r *http.Request, apiKey string) (context.Context, error) {
	apikeyRecord, err := a.apiKeyService.GetApiKey(r.Context(), apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	if apikeyRecord == nil || !apikeyRecord.Active {
		return nil, fmt.Errorf("invalid or inactive API key")
	}

	if !validateScopes(apikeyRecord.Scopes, r.Method, r.URL.Path) {
		return nil, fmt.Errorf("insufficient permissions for %s %s", r.Method, r.URL.Path)
	}

	userCtx := context.WithValue(r.Context(), auth.OrgContext{}, apikeyRecord.OrgId)
	return context.WithValue(userCtx, auth.ProjectContext{}, apikeyRecord.ProjectId), nil
}

// validateOAuthToken validates the OAuth token and returns the user context if successful.
func (a *Handler) validateOAuthToken(r *http.Request, token string) (context.Context, error) {
	claims, err := a.tokenService.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	userCtx := context.WithValue(r.Context(), auth.OrgContext{}, claims.OrgId)
	return context.WithValue(userCtx, auth.ProjectContext{}, claims.UserId), nil
}

// validateScopes checks if the provided scopes allow the requested operation based on the HTTP method and URL path.
// Scopes should be in the format "METHOD:path", e.g., "GET:/users" or "POST:/users/{userId}".
func validateScopes(scopes []string, method string, path string) bool {
	requiredScope := fmt.Sprintf("%s:%s", method, path)
	for _, scope := range scopes {
		if matchScope(scope, requiredScope) {
			return true
		}
	}

	return false
}

// matchScope checks if a single scope matches the required scope. This function supports basic wildcard matching.
func matchScope(scope, requiredScope string) bool {
	// Simple wildcard match: replace wildcards in scope with regex equivalents
	// Convert "GET:/items/*" to "^GET:/items/.*$"
	// Note: For more sophisticated patterns, consider using actual regex.
	pattern := "^" + strings.ReplaceAll(strings.ReplaceAll(scope, "*", ".*"), ":", ":") + "$"
	matched, err := regexp.MatchString(pattern, requiredScope)
	if err != nil {
		return false // in case of regex syntax error, fail closed
	}

	return matched
}
