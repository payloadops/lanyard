package handler

import (
	"context"
	"net/http"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_handler.go "plato/app/pkg/middleware/auth/handler" Authenticator

// Authenticator defines the interface for an Handler that handles authentication.
type Authenticator interface {
	Authenticate(r *http.Request) (context.Context, error)
}
