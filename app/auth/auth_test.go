package auth

import (
	"encoding/base64"
	"fmt"
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

func TestJWTAuthMiddleware(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "secret",
	}

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedOrgID  string
		expectedUserID string
		setupMocks     func() string
	}{
		{
			name:           "Valid JWT",
			expectedStatus: http.StatusOK,
			expectedOrgID:  "org123",
			expectedUserID: "user123",
			setupMocks: func() string {
				claims := Claims{
					OrgID: "org123",
					StandardClaims: jwt.StandardClaims{
						Subject:   "user123",
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))
				return "Bearer " + tokenString
			},
		},
		{
			name:           "Missing Authorization Header",
			expectedStatus: http.StatusUnauthorized,
			setupMocks:     nil,
		},
		{
			name:           "Invalid Authorization Format",
			authHeader:     "invalidFormat",
			expectedStatus: http.StatusUnauthorized,
			setupMocks:     nil,
		},
		{
			name:           "Invalid Signing Method",
			expectedStatus: http.StatusUnauthorized,
			setupMocks: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
					Subject:   "user123",
					ExpiresAt: time.Now().Add(time.Hour).Unix(),
				})
				tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))
				return "Bearer " + tokenString
			},
		},
		{
			name:           "Invalid JWT",
			expectedStatus: http.StatusUnauthorized,
			setupMocks: func() string {
				return "Bearer invalidToken"
			},
		},
		{
			name:           "Expired JWT",
			expectedStatus: http.StatusUnauthorized,
			setupMocks: func() string {
				claims := Claims{
					OrgID: "org123",
					StandardClaims: jwt.StandardClaims{
						Subject:   "user123",
						ExpiresAt: time.Now().Add(-time.Hour).Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))
				return "Bearer " + tokenString
			},
		},
		{
			name:           "Missing OrgID",
			expectedStatus: http.StatusUnauthorized,
			setupMocks: func() string {
				claims := Claims{
					StandardClaims: jwt.StandardClaims{
						Subject:   "user123",
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))
				return "Bearer " + tokenString
			},
		},
		{
			name:           "Missing UserID",
			expectedStatus: http.StatusUnauthorized,
			setupMocks: func() string {
				claims := Claims{
					OrgID: "org123",
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))
				return "Bearer " + tokenString
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authHeader := tt.authHeader
			if tt.setupMocks != nil {
				authHeader = tt.setupMocks()
			}

			r := chi.NewRouter()
			r.Use(JWTAuthMiddleware(cfg, zap.NewNop()))
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				orgID, _ := r.Context().Value("orgID").(string)   // Safely handle nil
				userID, _ := r.Context().Value("userID").(string) // Safely handle nil

				assert.Equal(t, tt.expectedOrgID, orgID)
				assert.Equal(t, tt.expectedUserID, userID)
				w.WriteHeader(http.StatusOK)
			})

			req, _ := http.NewRequest("GET", "/", nil)
			if authHeader != "" {
				req.Header.Set("Authorization", authHeader)
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

	mockAPIKeyManager := mocks.NewMockAPIKeyManager(mockCtrl)

	cfg := &config.Config{}
	tests := []struct {
		name              string
		authHeader        string
		expectedStatus    int
		expectedServiceID string
		expectedOrgID     string
		setupMocks        func()
	}{
		{
			name:              "Valid API Key",
			authHeader:        "Basic " + base64.StdEncoding.EncodeToString([]byte("validClientID:validSecret")),
			expectedStatus:    http.StatusOK,
			expectedServiceID: "service123",
			expectedOrgID:     "org123",
			setupMocks: func() {
				mockAPIKeyManager.EXPECT().
					GetAPIKey(gomock.Any(), "validClientID").
					Return(&dal.APIKey{Secret: "validSecret", Deleted: false, ServiceID: "service123", OrgID: "org123"}, nil).Times(1)
			},
		},
		{
			name:              "Missing Authorization Header",
			authHeader:        "",
			expectedStatus:    http.StatusUnauthorized,
			expectedServiceID: "",
			expectedOrgID:     "",
			setupMocks:        nil,
		},
		{
			name:              "Invalid Authorization Format",
			authHeader:        "invalidFormat",
			expectedStatus:    http.StatusUnauthorized,
			expectedServiceID: "",
			expectedOrgID:     "",
			setupMocks:        nil,
		},
		{
			name:              "Non-existent API Key",
			authHeader:        "Basic " + base64.StdEncoding.EncodeToString([]byte("nonexistentClientID:randomSecret")),
			expectedStatus:    http.StatusUnauthorized,
			expectedServiceID: "",
			expectedOrgID:     "",
			setupMocks: func() {
				mockAPIKeyManager.EXPECT().
					GetAPIKey(gomock.Any(), "nonexistentClientID").
					Return(nil, nil).Times(1) // Simulating key not found
			},
		},
		{
			name:              "Deleted API Key",
			authHeader:        "Basic " + base64.StdEncoding.EncodeToString([]byte("deletedClientID:anySecret")),
			expectedStatus:    http.StatusUnauthorized,
			expectedServiceID: "",
			expectedOrgID:     "",
			setupMocks: func() {
				mockAPIKeyManager.EXPECT().
					GetAPIKey(gomock.Any(), "deletedClientID").
					Return(&dal.APIKey{Secret: "anySecret", Deleted: true}, nil).Times(1)
			},
		},
		{
			name:              "Invalid Client Secret",
			authHeader:        "Basic " + base64.StdEncoding.EncodeToString([]byte("validClientID:invalidSecret")),
			expectedStatus:    http.StatusUnauthorized,
			expectedServiceID: "",
			expectedOrgID:     "",
			setupMocks: func() {
				mockAPIKeyManager.EXPECT().
					GetAPIKey(gomock.Any(), "validClientID").
					Return(&dal.APIKey{Secret: "validClientSecret", Deleted: false}, nil).Times(1)
			},
		},
		{
			name:              "Database Error",
			authHeader:        "Basic " + base64.StdEncoding.EncodeToString([]byte("validClientID:validSecret")),
			expectedStatus:    http.StatusInternalServerError,
			expectedServiceID: "",
			expectedOrgID:     "",
			setupMocks: func() {
				mockAPIKeyManager.EXPECT().
					GetAPIKey(gomock.Any(), "validClientID").
					Return(nil, fmt.Errorf("database error")).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			r := chi.NewRouter()
			r.Use(APIKeyAuthMiddleware(cfg, zap.NewNop(), mockAPIKeyManager))
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				serviceID, _ := r.Context().Value("serviceID").(string) // Safely handle nil
				orgID, _ := r.Context().Value("orgID").(string)         // Safely handle nil

				assert.Equal(t, tt.expectedServiceID, serviceID)
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
