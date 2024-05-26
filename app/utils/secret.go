package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateSecret generates a cryptographically secure random secret string.
func GenerateSecret(length int) (string, error) {
	// Generate random bytes
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	// Encode bytes to base64 to get a string
	apiKey := base64.URLEncoding.EncodeToString(bytes)
	return apiKey, nil
}
