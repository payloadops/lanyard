package service

import (
	"context"
	"github.com/payloadops/plato/app/openapi"
	"go.uber.org/zap"
	"net/http"
)

const (
	HealthCheckStatus_Healthy   = "healthy"
	HealthCheckStatus_Unhealthy = "unhealthy"
)

// HealthCheckAPIService is a service that implements the logic for the HealthCheckAPIServicer
// This service should implement the business logic for every endpoint for the HealthCheckAPI API.
// Include any external packages or services that will be required by this service.
type HealthCheckAPIService struct {
	healthy bool
	logger  *zap.Logger
}

// NewHealthCheckAPIService creates a default app service
func NewHealthCheckAPIService(logger *zap.Logger) openapi.HealthCheckAPIServicer {
	return &HealthCheckAPIService{healthy: true, logger: logger}
}

// HealthCheck - Health Check Endpoint
func (s *HealthCheckAPIService) HealthCheck(ctx context.Context) (openapi.ImplResponse, error) {
	if !s.healthy {
		return openapi.Response(http.StatusInternalServerError, openapi.HealthCheckErrorResponse{
			Status: HealthCheckStatus_Unhealthy,
		}), nil
	}

	return openapi.Response(http.StatusOK, openapi.HealthCheckSuccessResponse{
		Status: HealthCheckStatus_Healthy,
	}), nil
}
