package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAPIKey(t *testing.T) {
	// Test the generation of the API key
	keyLength := 32
	apiKey, err := GenerateSecret(keyLength)
	assert.NoError(t, err)
	assert.NotEmpty(t, apiKey)

	// Since base64 encoding might pad, check that the decoded length matches the original length
	decodedLen := len(apiKey) * 3 / 4
	assert.GreaterOrEqual(t, decodedLen, keyLength)
}
