package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseRFC3339Timestamp(t *testing.T) {
	t.Run("EmptyString", func(t *testing.T) {
		result, err := ParseTimestamp("")
		assert.NoError(t, err)
		assert.Equal(t, time.Time{}, result)
	})

	t.Run("ValidTimestamp", func(t *testing.T) {
		timestamp := "2024-05-16T15:04:05Z"
		expectedTime, _ := time.Parse(time.RFC3339, timestamp)

		result, err := ParseTimestamp(timestamp)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTime, result)
	})

	t.Run("InvalidTimestamp", func(t *testing.T) {
		_, err := ParseTimestamp("invalid-timestamp")
		assert.Error(t, err)
	})
}
