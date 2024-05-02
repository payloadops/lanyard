package util

import (
	"crypto/rand"
	"fmt"
)

// GenSecret generates a slice of random bytes of the specified length.
// It returns an error if the random byte generation fails.
func GenSecret(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return b, nil
}
