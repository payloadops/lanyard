package service_test

import (
	"context"
	"net/http"
	"testing"

	"go.uber.org/zap"

	"github.com/payloadops/lanyard/app/openapi"
	"github.com/payloadops/lanyard/app/service"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	t.Run("Returns healthy status when service is healthy", func(t *testing.T) {
		service := service.NewHealthCheckAPIService(zap.NewNop())

		resp, err := service.HealthCheck(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)

		status, ok := resp.Body.(openapi.HealthCheckSuccessResponse)
		assert.True(t, ok)
		assert.Equal(t, "healthy", status.Status)
	})
}
