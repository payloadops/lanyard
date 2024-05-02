package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenSecret(t *testing.T) {
	// Test generating secret with positive length
	secret, err := GenSecret(32)
	assert.NoError(t, err, "Should not have an error for valid length")
	assert.Len(t, secret, 32, "The length of the secret should be 32")

	// Test generating secret with zero length
	secret, err = GenSecret(0)
	assert.NoError(t, err, "Should not have an error for zero length")
	assert.Len(t, secret, 0, "The length of the secret should be 0")

	// Test generating secret with negative length
	secret, err = GenSecret(-1)
	assert.Error(t, err, "Should have an error for negative length")
	assert.Nil(t, secret, "The secret should be nil for negative length")
}
