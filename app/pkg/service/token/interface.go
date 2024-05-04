package token

import (
	"github.com/dgrijalva/jwt-go"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_token.go "plato/app/pkg/service/token" TokenService

// TokenService defines the interface for a service that handles token validation.
type TokenService interface {
	ValidateToken(tokenString string) (*TokenClaims, error)
}

// TokenClaims extends jwt.StandardClaims to include custom claim fields used in your application.
type TokenClaims struct {
	jwt.StandardClaims
	OrgId  string `json:"orgId"`
	UserId string `json:"userId"`
}
