package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/payloadops/plato/app/config"
	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/dal/mocks"
	"github.com/stretchr/testify/assert"
)

// generateTestToken generates a JWT token for testing with the specified secret, user ID, and organization ID.
func generateTestToken(secret, userID, orgID string, expired bool) string {
	expiresAt := time.Now().Add(time.Hour).Unix()
	if expired {
		expiresAt = time.Now().Add(-time.Hour).Unix()
	}
	claims := &Claims{
		OrgID: orgID,
		StandardClaims: jwt.StandardClaims{
			Subject:   userID,
			ExpiresAt: expiresAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func TestJWTAuthMiddleware(t *testing.T) {
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
			tokenString:    generateTestToken(cfg.JWTSecret, "user1", "org1", false),
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
			name:           "Invalid Authorization Header",
			tokenString:    "Bearer ",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Token",
			tokenString:    "invalid.token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Expired Token",
			tokenString:    generateTestToken(cfg.JWTSecret, "user1", "org1", true),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing sub",
			tokenString:    generateTestToken(cfg.JWTSecret, "", "org1", true),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing org",
			tokenString:    generateTestToken(cfg.JWTSecret, "user1", "", true),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Use(JWTAuthMiddleware(cfg, zap.NewNop()))
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				userID := r.Context().Value("userID").(string)
				orgID := r.Context().Value("orgID").(string)

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

func TestAPIKeyAuthMiddleware(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAPIKeyDBClient := mocks.NewMockAPIKeyManager(mockCtrl)

	cfg := &config.Config{}
	tests := []struct {
		name              string
		authHeader        string
		expectedStatus    int
		expectedProjectID string
		expectedOrgID     string
		setupMocks        func()
	}{
		{
			name:              "Valid API Key",
			authHeader:        "validClientID:validSecret",
			expectedStatus:    http.StatusOK,
			expectedProjectID: "project123",
			expectedOrgID:     "org123",
			setupMocks: func() {
				mockAPIKeyDBClient.EXPECT().
					GetAPIKey(gomock.Any(), "validClientID").
					Return(&dal.APIKey{Secret: "validSecret", Deleted: false, ProjectID: "project123", OrgID: "org123"}, nil).Times(1)
			},
		},
		{
			name:              "Missing Authorization Header",
			authHeader:        "",
			expectedStatus:    http.StatusUnauthorized,
			expectedProjectID: "",
			expectedOrgID:     "",
			setupMocks:        nil,
		},
		{
			name:              "Invalid Authorization Format",
			authHeader:        "invalidFormat",
			expectedStatus:    http.StatusUnauthorized,
			expectedProjectID: "",
			expectedOrgID:     "",
			setupMocks:        nil,
		},
		{
			name:              "Non-existent API Key",
			authHeader:        "nonexistentClientID:randomSecret",
			expectedStatus:    http.StatusUnauthorized,
			expectedProjectID: "",
			expectedOrgID:     "",
			setupMocks: func() {
				mockAPIKeyDBClient.EXPECT().
					GetAPIKey(gomock.Any(), "nonexistentClientID").
					Return(nil, nil).Times(1) // Simulating key not found
			},
		},
		{
			name:              "Deleted API Key",
			authHeader:        "deletedClientID:anySecret",
			expectedStatus:    http.StatusUnauthorized,
			expectedProjectID: "",
			expectedOrgID:     "",
			setupMocks: func() {
				mockAPIKeyDBClient.EXPECT().
					GetAPIKey(gomock.Any(), "deletedClientID").
					Return(&dal.APIKey{Secret: "anySecret", Deleted: true}, nil).Times(1)
			},
		},
		{
			name:              "Invalid Client Secret",
			authHeader:        "validClientID:invalidSecret",
			expectedStatus:    http.StatusUnauthorized,
			expectedProjectID: "",
			expectedOrgID:     "",
			setupMocks: func() {
				mockAPIKeyDBClient.EXPECT().
					GetAPIKey(gomock.Any(), "validClientID").
					Return(&dal.APIKey{Secret: "validClientSecret", Deleted: false}, nil).Times(1)
			},
		},
		{
			name:              "Database Error",
			authHeader:        "validClientID:validSecret",
			expectedStatus:    http.StatusUnauthorized,
			expectedProjectID: "",
			expectedOrgID:     "",
			setupMocks: func() {
				mockAPIKeyDBClient.EXPECT().
					GetAPIKey(gomock.Any(), "validClientID").
					Return(nil, nil).Times(1) // Simulating key not found
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMocks != nil {
				tt.setupMocks()
			}
			r := chi.NewRouter()
			r.Use(APIKeyAuthMiddleware(cfg, zap.NewNop(), mockAPIKeyDBClient))
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				projectID, _ := r.Context().Value("projectID").(string) // Safely handle nil
				orgID, _ := r.Context().Value("orgID").(string)         // Safely handle nil

				assert.Equal(t, tt.expectedProjectID, projectID)
				assert.Equal(t, tt.expectedOrgID, orgID)
				w.WriteHeader(http.StatusOK)
			})

			req, _ := http.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
