package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID string `json:"user_id"` // This could also be `sub` if the user ID is stored in the subject claim
	jwt.StandardClaims
}

// ValidateOAuthToken validates a JWT locally by checking its signature and claims.
func ValidateOAuthToken(tokenString string) (*Claims, error) {
	// This is your public key which matches the private key used by the OAuth server.
	// It should be parsed on application start-up and passed around or made globally accessible.
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM([]byte("your-public-key-PEM-data"))

	// Parse takes the token string and a function for looking up the key to perform verification.
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Additionally check if the token is expired or not
		if claims.ExpiresAt < time.Now().Unix() {
			return nil, fmt.Errorf("token has expired")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
