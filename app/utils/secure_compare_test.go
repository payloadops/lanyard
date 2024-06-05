package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecureCompare(t *testing.T) {
	t.Run("CompareSameString", func(t *testing.T) {
		result := SecureCompare("abc123", "abc123")
		assert.Equal(t, true, result)
	})

	t.Run("CompareDiffString", func(t *testing.T) {
		result := SecureCompare("xyz987", "abc123")
		assert.Equal(t, false, result)
	})
}
