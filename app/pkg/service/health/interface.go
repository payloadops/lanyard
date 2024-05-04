package health

//go:generate mockgen -package=mocks -destination=mocks/mock_health.go "plato/app/pkg/service/health" HealthService

// HealthService defines an interface for checking health of the service.
type HealthService interface {
	CheckHealth() CheckHealthResponse
}

// CheckHealthResponse represents the response from a health check.
type CheckHealthResponse struct {
	// Healthy indicates whether the service is healthy (true) or not (false).
	Healthy bool `json:"healthy"`
	// Status provides a descriptive status of the health check.
	Status string `json:"status"`
}
