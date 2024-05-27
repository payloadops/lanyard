package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKSUID(t *testing.T) {
	t.Run("Generate KSUID successfully", func(t *testing.T) {
		id, err := GenerateKSUID()
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
		assert.Equal(t, 27, len(id)) // KSUIDs are 27 characters long
	})
}
