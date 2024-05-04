package handler

import (
	"context"
	"errors"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	dbdal "plato/app/pkg/dal/postgres"
	apikeymocks "plato/app/pkg/service/apikey/mocks"
	"plato/app/pkg/service/token"
	tokenmocks "plato/app/pkg/service/token/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateAPIKey checks edge cases for API key validation
func TestValidateAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)

	authenticator := &Handler{}
	apiKeyService := apikeymocks.NewMockApiKeyService(ctrl)
	authenticator.apiKeyService = apiKeyService

	testCases := []struct {
		name         string
		apiKey       string
		apiKeyRecord *dbdal.ApiKeyItem
		expectError  bool
		errorMessage string
	}{
		{
			name:         "Valid API Key",
			apiKey:       "validKey",
			apiKeyRecord: &dbdal.ApiKeyItem{Active: true, OrgId: "org1", ProjectId: "proj1", Scopes: []string{"GET:/protected"}},
			expectError:  false,
		},
		{
			name:         "Inactive API Key",
			apiKey:       "inactiveKey",
			apiKeyRecord: &dbdal.ApiKeyItem{Active: false},
			expectError:  true,
			errorMessage: "invalid or inactive API key",
		},
		{
			name:         "API Key Not Found",
			apiKey:       "notFoundKey",
			expectError:  true,
			errorMessage: "failed to get API key: key not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			request := httptest.NewRequest("GET", "/protected", nil)
			if tc.expectError {
				apiKeyService.EXPECT().
					GetApiKey(ctx, tc.apiKey).
					Return(nil, errors.New(tc.errorMessage))
			} else {
				apiKeyService.EXPECT().
					GetApiKey(ctx, tc.apiKey).
					Return(tc.apiKeyRecord, nil)
			}

			_, err := authenticator.validateAPIKey(request, tc.apiKey)
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestValidateOAuthToken checks edge cases for OAuth token validation
func TestValidateOAuthToken(t *testing.T) {
	ctrl := gomock.NewController(t)

	authenticator := &Handler{}
	tokenService := tokenmocks.NewMockTokenService(ctrl)
	authenticator.tokenService = tokenService

	testCases := []struct {
		name         string
		token        string
		tokenClaims  *token.TokenClaims
		expectError  bool
		errorMessage string
	}{
		{
			name:        "Valid Token",
			token:       "validToken",
			tokenClaims: &token.TokenClaims{OrgId: "org1", UserId: "user1"},
			expectError: false,
		},
		{
			name:         "Invalid Token",
			token:        "invalidToken",
			expectError:  true,
			errorMessage: "failed to validate token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/protected", nil)
			if tc.expectError {
				tokenService.EXPECT().
					ValidateToken(tc.token).
					Return(nil, errors.New(tc.errorMessage))
			} else {
				tokenService.EXPECT().
					ValidateToken(tc.token).
					Return(tc.tokenClaims, nil)
			}

			_, err := authenticator.validateOAuthToken(request, tc.token)
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestValidateScopes checks edge cases for scope validation
func TestValidateScopes(t *testing.T) {
	assert.True(t, validateScopes([]string{"GET:/items/*"}, "GET", "/items/123"), "Scope should match")
	assert.False(t, validateScopes([]string{"POST:/items"}, "GET", "/items"), "Scope should not match")
}
