package util

import (
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenIDString(t *testing.T) {
	// Generate an ID string
	idString := GenIDString()

	// Check if the generated ID string is not empty
	assert.NotEmpty(t, idString, "Generated ID string is empty")

	// Parse the generated ID string
	parsedID, err := xid.FromString(idString)
	assert.NoError(t, err, "Failed to parse generated ID string")

	// Check if the parsed ID is nil
	assert.False(t, parsedID.IsNil(), "Parsed ID is nil")

	// Check if the parsed ID is zero
	assert.False(t, parsedID.IsZero(), "Parsed ID is zero")
}
