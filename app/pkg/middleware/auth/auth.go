package auth

import (
	"net/http"
	"plato/app/pkg/middleware/auth/handler"
)

// AuthMiddleware represents the dependencies and settings of the authentication middleware.
type AuthMiddleware struct {
	authenticator handler.Authenticator
}

// NewAuthMiddleware creates a new instance of AuthMiddlewareStruct.
func NewHandler() *AuthMiddleware {
	return &AuthMiddleware{
		authenticator: handler.NewHandler(),
	}
}

// Handler creates an HTTP middleware for authenticating users either through API keys or OAuth tokens.
func (a *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := a.authenticator.Authenticate(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(userCtx))
	})
}
