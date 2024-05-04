package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestService_ValidateToken(t *testing.T) {
	// Setup the environment variable
	secretKey := "supersecret"
	os.Setenv("JWT_SECRET_KEY", secretKey)

	// Setup
	svc := NewService()
	signingKey := []byte(secretKey)
	wrongKey := []byte("supersecret1")

	// Create a valid token for testing
	validClaims := &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // valid for 1 hour
		},
		OrgId:  "org123",
		UserId: "user456",
	}
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims)
	validTokenString, err := validToken.SignedString(signingKey)
	assert.NoError(t, err)

	// Test cases
	tests := []struct {
		name        string
		tokenString string
		expectError bool
		errorString string
	}{
		{
			name:        "Valid token",
			tokenString: validTokenString,
			expectError: false,
		},
		{
			name:        "Expired token",
			tokenString: createToken(signingKey, -1), // expired token
			expectError: true,
			errorString: "token is expired by 1h0m0s",
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.here",
			expectError: true,
			errorString: "error parsing token",
		},
		{
			// This edge case is unreachable, different signing methods require different types of keys.
			name:        "Unexpected signing method",
			tokenString: createTokenWithDifferentMethod(signingKey),
			expectError: true,
			errorString: "error parsing token: token contains an invalid number of segments",
		},
		{
			name:        "Wrong key key",
			tokenString: createToken(wrongKey, 1),
			expectError: true,
			errorString: "error parsing token: signature is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := svc.ValidateToken(tt.tokenString)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorString)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, "org123", claims.OrgId)
				assert.Equal(t, "user456", claims.UserId)
			}
		})
	}
}

func createToken(key []byte, hoursTillExpire int) string {
	claims := &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(hoursTillExpire)).Unix(),
		},
		OrgId:  "org123",
		UserId: "user456",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(key)
	return tokenString
}

func createTokenWithDifferentMethod(key []byte) string {
	claims := &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		OrgId:  "org123",
		UserId: "user456",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, _ := token.SignedString(key)
	return tokenString
}
