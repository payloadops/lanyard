package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/payloadops/plato/api/openapi"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	t.Run("Returns healthy status when service is healthy", func(t *testing.T) {
		service := HealthCheckAPIService{healthy: true}

		resp, err := service.HealthCheck(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)

		status, ok := resp.Body.(openapi.HealthCheck200Response)
		assert.True(t, ok)
		assert.Equal(t, HealthCheckStatus_Healthy, status.Status)
	})

	t.Run("Returns unhealthy status when service is not healthy", func(t *testing.T) {
		service := HealthCheckAPIService{healthy: false}

		resp, err := service.HealthCheck(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)

		status, ok := resp.Body.(openapi.HealthCheck500Response)
		assert.True(t, ok)
		assert.Equal(t, HealthCheckStatus_Unhealthy, status.Status)
	})
}
