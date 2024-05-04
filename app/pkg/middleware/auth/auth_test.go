package auth

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"plato/app/pkg/middleware/auth/handler/mocks"
	"testing"
)

// TestAuthMiddleware tests the AuthMiddleware function to ensure it handles the request appropriately based on authentication success or failure.
func TestAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Setup
	authenticator := mocks.NewMockAuthenticator(ctrl)
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // What you expect to send after successful authentication
	})

	middleware := &AuthMiddleware{
		authenticator: authenticator,
	}

	handler := middleware.Handler(nextHandler)
	tests := []struct {
		name           string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "Successful Authentication",
			mockSetup: func() {
				authenticator.EXPECT().
					Authenticate(gomock.Any()).
					Return(context.Background(), nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Failed Authentication",
			mockSetup: func() {
				authenticator.EXPECT().
					Authenticate(gomock.Any()).
					Return(nil, fmt.Errorf("unauthorized"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup the request and response recorder
			req := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()

			// Setup the specific test case behavior
			tc.mockSetup()

			// Execute the middleware with the request
			handler.ServeHTTP(rr, req)

			// Assert that the status code is what we expect
			assert.Equal(t, tc.expectedStatus, rr.Code, "Expected response status code to match")
		})
	}
}
