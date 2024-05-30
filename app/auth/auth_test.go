package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/payloadops/plato/app/config"
	"github.com/stretchr/testify/assert"
)

// generateTestToken generates a JWT token for testing with the specified secret, user ID, and organization ID.
func generateTestToken(secret, userID, orgID string) string {
	claims := &Claims{
		UserID: userID,
		OrgID:  orgID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func TestAuthMiddleware(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "testsecret",
	}

	tests := []struct {
		name           string
		tokenString    string
		expectedStatus int
		expectedUserID string
		expectedOrgID  string
	}{
		{
			name:           "Valid Token",
			tokenString:    generateTestToken(cfg.JWTSecret, "user1", "org1"),
			expectedStatus: http.StatusOK,
			expectedUserID: "user1",
			expectedOrgID:  "org1",
		},
		{
			name:           "Missing Authorization Header",
			tokenString:    "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Token",
			tokenString:    "invalid.token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Use(AuthMiddleware(cfg))
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				userID := r.Context().Value("UserID").(string)
				orgID := r.Context().Value("OrgID").(string)

				assert.Equal(t, tt.expectedUserID, userID)
				assert.Equal(t, tt.expectedOrgID, orgID)
				w.WriteHeader(http.StatusOK)
			})

			req, _ := http.NewRequest("GET", "/", nil)
			if tt.tokenString != "" {
				req.Header.Set("Authorization", "Bearer "+tt.tokenString)
			}
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
