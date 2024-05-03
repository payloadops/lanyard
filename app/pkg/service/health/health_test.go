package health

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestCheckHealthHealthy tests the CheckHealth method for expected healthy response.
func TestCheckHealthHealthy(t *testing.T) {
	service := NewService()
	response := service.CheckHealth()

	// Since the function is mocked to always return healthy, adjust these tests according to real conditions
	assert.True(t, response.Healthy, "Expected the service to be healthy")
	assert.Equal(t, "Service is healthy", response.Status, "Expected health status description does not match")
}
