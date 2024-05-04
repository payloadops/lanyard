package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

var _ TokenService = (*Service)(nil)

// Service implements the TokenService interface for JWT token operations.
type Service struct {
	SecretKey []byte
}

// NewService initializes a new instance of Service.
func NewService() *Service {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		panic("JWT_SECRET_KEY must be set")
	}
	return &Service{
		SecretKey: []byte(secretKey),
	}
}

// ValidateToken validates the given JWT token and returns the extracted claims if the token is valid.
func (s *Service) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.SecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
			return nil, fmt.Errorf("token has expired")
		}

		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
